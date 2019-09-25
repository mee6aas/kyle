package pool

// Destory stops management for the pool.
// It does not deletes spawned runtimes.
func Destory() {
	mngCancel()
}
