package taskmodel

import (
	"fmt"

	usermodel "github.com/Negat1v9/work-marketplace/model/userModel"
)

// Info: CalculateRespondStarPrice - calculation of the price of the respond to the task, how much it will cost the worker (amount in telegram stars valut)
// if meta or wi == nil return 0
func CalculateRespondStarPrice(meta *TaskMeta, balance *usermodel.Balance) int {
	if meta == nil || balance == nil {
		return 0
	}
	sumMinMax := meta.MinPrice + meta.MaxPrice
	switch {
	// discount for workers with only two refills
	case balance.NumberPayments <= 2:
		return 10
	case sumMinMax <= 1000:
		return 80
	case sumMinMax <= 5000:
		return 120
	default:
		return 180
	}

}

// checking the length of fields before creating or updating
func (m *TaskMeta) ValidateMetaFields() error {
	switch {
	case len(m.FormEducation) > 100:
		return fmt.Errorf("max field \"form_education\" length is %d", 100)
	case len(m.University) > 150:
		return fmt.Errorf("max field \"university\" length is %d", 150)
	case len(m.Subject) > 100:
		return fmt.Errorf("max field \"subject\" length is %d", 100)
	case len(m.TaskType) > 50:
		return fmt.Errorf("max field \"task_type\" length is %d", 50)
	case len(m.Description) > 400:
		return fmt.Errorf("max field \"description\" length is %d", 400)
	}
	return nil
}
