package filehasher

import (
	"filesigner/stringhelper"
	"path/filepath"
	"runtime"
	"sync"
)

// ******** Public types ********

type HashResult struct {
	FilePath  string
	HashValue []byte
	Err       error
}

// ******** Public functions ********

// FileHashes computes the hashes of the supplied files in an asynchronous manner.
func FileHashes(filePaths []string, contextId string) map[string]*HashResult {
	// hasherWaitGroup is used to wait for all hashers to finish
	var hasherWaitGroup sync.WaitGroup

	// hasherResultChannel is where the hashers place their results to be picked off by this function
	hasherResultChannel := make(chan *HashResult, runtime.NumCPU())

	contextBytes := stringhelper.UnsafeStringBytes(contextId)
	// Start an asynchronous hasher for each file to hash
	numHashes := startFileHashers(filePaths, contextBytes, &hasherWaitGroup, &hasherResultChannel)

	// Start an asynchronous function that waits for all hashers to finish and then close the hasherResultChannel
	go waitForAllHashers(&hasherWaitGroup, &hasherResultChannel)

	// Collect all results and return when hasherResultChannel is closed
	return makeResultList(numHashes, &hasherResultChannel)
}

// ******** Private functions ********

// startFileHashers starts the file hasher processes asynchronously
func startFileHashers(filePaths []string,
	contextBytes []byte,
	hasherWaitGroup *sync.WaitGroup,
	hasherResultChannel *chan *HashResult) int {
	numHashes := 0

	var normalizedPath string

	for _, aFilePath := range filePaths {
		normalizedPath = filepath.ToSlash(aFilePath)

		numHashes++
		hasherWaitGroup.Add(1) // This must be done before the start of the goroutine, so that the waiter will have to wait for the first goroutine to start
		go fileHashWorker(normalizedPath, contextBytes, hasherWaitGroup, hasherResultChannel)
	}

	return numHashes
}

// makeResultList reads hash results from the result channel and returns when the hash result channel is closed.
func makeResultList(numHashes int, hasherResultChannel *chan *HashResult) map[string]*HashResult {
	resultList := make(map[string]*HashResult, numHashes)

	var aResult *HashResult
	var resultPresent bool
	for {
		aResult, resultPresent = <-*hasherResultChannel

		if resultPresent {
			resultList[aResult.FilePath] = aResult
		} else {
			// We come here when the hash result channel has been closed.
			// So we can stop waiting for results and return.
			break
		}
	}

	return resultList
}

// waitForAllHashers waits for all hashers to complete and closes the hash result channel then
func waitForAllHashers(hasherWaitGroup *sync.WaitGroup, hasherResultChannel *chan *HashResult) {
	// This function assumes that the hasherWaitGroup already has the number of running goroutines set
	hasherWaitGroup.Wait()

	close(*hasherResultChannel)
}

// fileHashWorker calculates the hash value of one file
func fileHashWorker(filePath string,
	contextBytes []byte,
	hasherWaitGroup *sync.WaitGroup,
	hasherResultChannel *chan *HashResult) {
	defer hasherWaitGroup.Done()

	result := &HashResult{}
	result.FilePath = filePath
	fileHasher, err := NewFileHasher(contextBytes)
	if err == nil {
		result.HashValue, result.Err = fileHasher.HashFile(filePath)
	} else {
		result.HashValue = nil
		result.Err = err
	}

	*hasherResultChannel <- result
}
