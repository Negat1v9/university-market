package taskmodel

import usermodel "github.com/Negat1v9/work-marketplace/model/userModel"

// Info: CalculateRespondStarPrice - calculation of the price of the respond to the task, how much it will cost the worker (amount in telegram stars valut)
// if meta or wi == nil return 0
func CalculateRespondStarPrice(meta *TaskMeta, wi *usermodel.WorkerInfo) int {
	if meta == nil || wi == nil {
		return 0
	}
	switch {
	// case wi.WorkCompleted == 0:
	// 	return 1
	case meta.MinPrice <= 500:
		return 100
	case meta.MinPrice <= 1000:
		return 150
	default:
		return 200
	}

}
