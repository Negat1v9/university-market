package usermodel

import "fmt"

// checking which fields need to be updated and return WorkerInfo with updated fields

func (i *WorkerInfo) FindUpdates(new *WorkerInfo) *WorkerInfo {
	if new.FullName != "" {
		i.FullName = new.FullName
	}
	if new.Education != "" {
		i.Education = new.Education
	}
	if new.Experience != "" {
		i.Experience = new.Experience
	}
	if new.Description != "" {
		i.Description = new.Description
	}

	return i
}

// checking the length of fields before creating or updating
func (i *WorkerInfo) ValidateFields() error {
	switch {
	case len(i.FullName) > 100:
		return fmt.Errorf("max field \"full_name\" length is %d", 100)
	case len(i.Education) > 100:
		return fmt.Errorf("max field \"education\" length is %d", 100)
	case len(i.Experience) > 200:
		return fmt.Errorf("max field \"experience\" length is %d", 200)
	case len(i.Description) > 400:
		return fmt.Errorf("max field \"description\" length is %d", 400)
	}
	return nil
}
