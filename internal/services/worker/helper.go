package workerservice

// Info: checkWorkerAlredyRespond - checks whether the worker has responded to the task return true if worker responded and return false if no
func checkWorkerAlredyRespond(wokrerID string, responds []string) bool {
	for i := 0; i < len(responds); i++ {
		if wokrerID == responds[i] {
			return true
		}
	}
	return false
}

// // checking the task before choosing a worker
// func checkTaskFree(status taskmodel.TaskStatus) error {
// 	if status != taskmodel.WaitingExecution {
// 		return
// 	}
// 	return nil
// }
