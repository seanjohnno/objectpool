## Object Pool for golang

### Description

Used to add and retrieve objects/structs from a common pool rather than having them GC'd and new ones created. Interface is contained in [object_pool.go](https://github.com/seanjohnno/objectpool/blob/master/object_pool.go):

```
// ObjectPool allows us to store & reuse objects (so they're not GC'd)
type ObjectPool interface {

	// Add adds an object to our object pool
	Add(obj interface{})
	
	// Gets an object from our object pool
	Retrieve() (interface{}, bool)
}
```

Created the interface so we can have different underlying implementations of the object pool. Currently the only implementation is [ExpiringObjectPool](https://github.com/seanjohnno/objectpool/blob/master/expiring_obj_pool.go)

### ExpiringObjectPool

Implementation ensures that objects expire and are removed from the pool if they're not accessed within a given time period

### When would I want to use this?

If you find yourself creating lots of a specific type of obejct rapidly and then discard them quickly i.e. they have a short lifespan . Instead of discarding them and having them GC'd they can be added to our pool and re-used if they're requested within the expiry period.

The incentive for creating this was creating buffers to serve HTTP content. If there's a spike in traffic we add used buffers to a shared pool and have them re-used for subsequent requests rather than creating another buffer. When traffic dies down and they're not required they'll be auto removed from our pool. 

### Quick Example

```
  func main() {
  	completeChan := make(chan bool, 100)
  
  	pool := objpool.NewTimedExiryPool(3000)
  	for i := 0; i < 100; i++ {
  		go ReuseFunc(pool, completeChan)
  	}
  
  	completeCount := 0
  	for {
  		<- completeChan
  		if completeCount++; completeCount == 100 {
  			return
  		}
  	}
  }
  
  func ReuseFunc(pool objpool.ObjectPool, completeChan chan<- bool) {
  	if item, present := pool.Retrieve(); present {
  		fmt.Println("Woohoo! found an existing buffer...")
  		// ...here we'd do something with our buffer...
  		pool.Add(item)
  	} else {
  		fmt.Println("Have to create new buffer...")
  		item := bytes.NewBuffer(make([]byte, 50))			// Nothing found in pool so create new
  		// ...here we'd do something with our buffer...
  		pool.Add(item)										// we're done with it, lets add back into pool so something else can use
  	}
  	completeChan <- true
  }
```

### Full Example

You can see the example above with imports and stuff [here](https://github.com/seanjohnno/goexamples/blob/master/object_pool_example.go). Continue reading if you require instructions on how to grab the sourcecode and/or example from within the command-line...

### Setup

Create your Go folder structure on the filesystem (if you have't already):

```
GoProjects
  |- src
  |- pkg
  |- bin
```
In your command-line set your **GOPATH** environment variable:

* Linux: `export GOPATH=<Replace_me_with_path_to>\GoProjects`
* Windows: `set GOPATH="<Replace_me_with_path_to>\GoProjects"`

Browse to your *GoProjects* folder in the command-line and enter:

  `go get github.com/seanjohnno/objectpool`

You should see the folders */github.com/seanjohnno/objectpool* under *src* and the code inside *objectpool*

If you want to run the example then make sure you're in your *GoProjects* folder and run:

  `go get github.com/seanjohnno/goexamples`

Navigate to the *goexamples directory* and run the following:

```
  go build object_pool_example.go
```

...and then depending on your OS:

* Linux: `./object_pool_example`
* Windows: `object_pool_example.exe`
  



