package objpool

// ObjectPool allows us to store & reuse objects (so they're not GC'd)
type ObjectPool interface {

	// Add adds an object to our object pool
	Add(obj interface{})
	
	// Gets an object from our object pool
	Retrieve() (interface{}, bool)
}