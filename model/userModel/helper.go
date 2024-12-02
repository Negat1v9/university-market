package usermodel

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
