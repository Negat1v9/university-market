package workerservice

import "fmt"

// Info: checkWorkerAlredyRespond - checks whether the worker has responded to the task return true if worker responded and return false if no
func checkWorkerAlredyRespond(wokrerID string, responds []string) bool {
	for i := 0; i < len(responds); i++ {
		if wokrerID == responds[i] {
			return true
		}
	}
	return false
}

// marks that the files were sent to the worker, returns an error if they have already been sent
func checkFilesAlredySend(workerID string, filesSendTo []string) error {
	for i := 0; i < len(filesSendTo); i++ {
		if workerID == filesSendTo[i] {
			return fmt.Errorf("%s already received files", workerID)
		}
	}
	return nil
}

// if the value of maxNum is greater than or equal to the number of responds, it returns an error
func checkCountRespondsOnTask(maxNum int, responds []string) error {
	if len(responds) >= maxNum {
		return fmt.Errorf("task already has too many responds: %d", len(responds))
	}
	return nil
}
