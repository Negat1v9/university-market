package adminservice

import (
	eventmodel "github.com/Negat1v9/work-marketplace/model/event"
	httpresponse "github.com/Negat1v9/work-marketplace/pkg/httpResponse"
)

// if the user ID is in the slice it returns true
func checkIsAdmin(userTgID int64, adminsIds []int64) bool {
	for i := 0; i < len(adminsIds); i++ {
		if userTgID == adminsIds[i] {
			return true
		}
	}
	return false
}

func validateEvent(event *eventmodel.Event) error {
	switch {

	case event.UserType != "worker" && event.UserType != "all":
		return httpresponse.NewError(422, "user_type field must contain worker or all")
	case event.Caption == "":
		return httpresponse.NewError(422, "caption field is required")
	}

	return nil
}
