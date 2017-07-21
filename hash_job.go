package ut4updater

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// HashJob is the structure that will be passed for processing of hashes
type HashJob struct {
	filepath string
	progress chan HashProgressEvent
}

// Process is an implementation of Job.Process()
func (job HashJob) Process() {
	filename := filepath.Base(job.filepath)
	fileInfo, err := os.Stat(job.filepath)
	if err != nil {
		job.progress <- HashProgressEvent{
			Filename: filename,
			Filepath: job.filepath,
			Error:    err.Error(),
		}
		return
	}
	file, err := os.Open(job.filepath)
	if err != nil {
		if err != nil {
			job.progress <- HashProgressEvent{
				Filename: filename,
				Filepath: job.filepath,
				Error:    err.Error(),
			}
			return
		}
	}
	defer file.Close()

	// Set up an internal hash progress tracker
	hashProgressChan := make(chan int)
	hasher := sha256.New()
	fileReader := io.TeeReader(file, ProgressTracker{progress: hashProgressChan})

	// Start the hashing in the background since large files (.pak) files
	// could take quite some time to complete
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := io.Copy(hasher, fileReader)
		// Send the error back to the progress tracker
		job.progress <- HashProgressEvent{
			Filename: filename,
			Filepath: job.filepath,
			Error:    err.Error(),
		}
		return
	}()

	bytesHashed := 0
	// Report hashing progress once a second
	progressReportInterval := time.Second
	lastProgressReport := time.Now()
	// Track the bytes per second to determine estimated
	// time remaining
	bytesPerSecond := 0
	bytesLeft := fileInfo.Size()
	for hashedBytes := range hashProgressChan {
		bytesLeft -= int64(hashedBytes)
		bytesHashed += hashedBytes
		bytesPerSecond += hashedBytes
		// If enough time has passed, report the progress
		if time.Since(lastProgressReport) > progressReportInterval {
			job.progress <- HashProgressEvent{
				Filename: filename,
				Filepath: job.filepath,
				Mbps:     float64(bytesPerSecond) / 1024.00 / 1024.00,
				ETA:      float64(bytesLeft) / float64(bytesPerSecond),
			}
			// Reset counters
			lastProgressReport = time.Now()
			bytesPerSecond = 0
		}
	}
	wg.Wait()
	job.progress <- HashProgressEvent{
		Filename: filename,
		Filepath: job.filepath,
		Hash:     fmt.Sprintf("%x", hasher.Sum(nil)),
	}
}
