package cellar

type FileLock interface {
	TryLock() (bool, error)
	Lock() error
	Unlock() error
}
