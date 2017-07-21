package ut4updater

// ProgressTracker implements a writer for tracking the
// progress of hash generation and file downloads
type ProgressTracker struct {
	progress chan int
}

func (pt ProgressTracker) Write(data []byte) (int, error) {
	pt.progress <- len(data)
	// TODO: Fix this magic number, it's the buffer size
	if len(data) < 32768 {
		close(pt.progress)
	}
	return len(data), nil
}
