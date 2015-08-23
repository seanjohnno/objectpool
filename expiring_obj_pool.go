package objpool

import(
	"github.com/seanjohnno/semaphore"
	"sync"
	"container/list"
	"time"
)

// ------------------------------------------------------------------------------------------------------------------------
// Struct: ExpiryElement
// ------------------------------------------------------------------------------------------------------------------------

// ExpiryElement is used to store an expiry along with the stored object
type ExpiryElement struct {

	// Item is the underlying object we're storing
	Item interface{}

	// ExpiryTime is the time at which this object will expire and be removed from the list
	ExpiryTime time.Time
}

// ------------------------------------------------------------------------------------------------------------------------
// Struct: ExpiringObjectPool
// ------------------------------------------------------------------------------------------------------------------------

// ExpiringObjectPool implements ObjectPool. Objects that aren't accessed are auto removed within a specified time limit
type ExpiringObjectPool struct {

	// LinkedList is used to maintain list of objects
	LinkedList *list.List

	// ExpiryTime stores the amount of milliseconds that an object is valid for
	ExpiryPeriod uint64

	//
	Sem *semaphore.CountingSemaphore

	Mutex sync.Mutex
}

// ------------------------------------------------------------------------------------------------------------------------
// Implementing ObjectPool
// ------------------------------------------------------------------------------------------------------------------------

// Add adds an object to our object pool
func (this *ExpiringObjectPool) Add(obj interface{}) {

	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Add item to the front of the list
	expiryTime := time.Now().Add(time.Duration(this.ExpiryPeriod * uint64(time.Millisecond)))
	this.LinkedList.PushFront( &ExpiryElement{ Item:obj, ExpiryTime: expiryTime } )

	// Signals semaphore, an exiting Wait() will unblock or future Wait() won't block
	this.Sem.Signal()
}
	
// Gets an object from our object pool (if one is available)
//
// The object is removed from the pool
func (this *ExpiringObjectPool) Retrieve() (interface{}, bool) {

	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	elem := this.LinkedList.Front()
	if elem != nil {
		this.LinkedList.Remove(elem)
		return elem.Value.(*ExpiryElement).Item, true
	} else {
		return nil, false
	}
}

// removeExpiredItems removes items in the list that have expired
func (this *ExpiringObjectPool) removeExpiredItems() {
	for {
		
		// It'll sit here until Unlock is called from Add method
		this.Sem.Wait()

		if lastElem := this.LinkedList.Back(); lastElem != nil {

			// Sleep until its time for item to expire
			expiryTime := lastElem.Value.(*ExpiryElement).ExpiryTime
			time.Sleep(expiryTime.Sub(time.Now()))

			// Altering LL so obtain lock
			this.Mutex.Lock()

			// Check last element hasn't been retrieved. Remove if it hasn't as its expired
			
			if lastElem = this.LinkedList.Back(); lastElem != nil {

				expiryTime := lastElem.Value.(*ExpiryElement).ExpiryTime
				timeNow := time.Now()

				if (timeNow.After(expiryTime) || timeNow.Equal(expiryTime)) {
					this.LinkedList.Remove(lastElem)
				}
			}

			// Done altering so release lock
			this.Mutex.Unlock()
		}
	}
}

// ------------------------------------------------------------------------------------------------------------------------
// Construction
// ------------------------------------------------------------------------------------------------------------------------

// NewTimedExiryPool creates a ...NewTimedExiryPool
func NewTimedExiryPool(expiryTimeInMilli uint64) (*ExpiringObjectPool) {
	tep := &ExpiringObjectPool{ LinkedList: list.New(), ExpiryPeriod: expiryTimeInMilli, Sem: semaphore.New() }
	go tep.removeExpiredItems()
	return tep
}