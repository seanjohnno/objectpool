package objpool

import (
	"testing"
	"time"

	//"fmt"
)

var (
	Expiry = uint64(100)
	HalfExpiry = time.Duration(50 * uint64(time.Millisecond))
	BeforeExpiry = time.Duration(70 * uint64(time.Millisecond))
	SleepPastExpiry = time.Duration(120 * uint64(time.Millisecond))
)

func TestExpiringObjectPool(t * testing.T) {

	// Create ExpiringObjectPool and check its not nil
	eop := NewTimedExiryPool(Expiry)
	if eop == nil {
		t.Error("Shouldn't be nil")
	}

	// Check we can add and retrieve the object (assumption is that this should easily complete within 100ms)
	eop.Add(&Dummy{})
	if _, present := eop.Retrieve(); !present {
		t.Error("Should have been able to retrieve dummy item")
	}

	// Check that we cant grab anther item (it was just removed)
	if _, present := eop.Retrieve(); present {
		t.Error("Shouldn't have been an item present")
	}

	// Add an item and sleep past expiry time
	eop.Add(&Dummy{})
	time.Sleep(SleepPastExpiry)

	// Check that we cant grab an item as they should have expired
	if _, present := eop.Retrieve(); present {
		t.Error("Item should have expired")
	}

	// Add a bunch of items
	for i := 0; i < 50; i++ {
		eop.Add(&Dummy{})
	}
	time.Sleep(SleepPastExpiry)
	// Check that we cant grab an item as they should have expired
	if _, present := eop.Retrieve(); present {
		t.Error("All the items should have expired")
	}

	// Check that we're still present after some delay
	eop.Add(&Dummy{})
	time.Sleep(BeforeExpiry)
	if _, present := eop.Retrieve(); !present {
		t.Error("Should have been present after the shorter delay")
	}

	// Test that the fist expires but the second and 3rd are still present
	eop.Add(&Dummy{})
	time.Sleep(HalfExpiry)
	eop.Add(&Dummy{})
	eop.Add(&Dummy{})
	time.Sleep(HalfExpiry)

	i := 0
	for _, present := eop.Retrieve(); present == true;  _, present = eop.Retrieve() {
		i++
	}
	if i != 2 {
		t.Error("Should have been 2 items left")
	}
}

type Dummy struct {
}