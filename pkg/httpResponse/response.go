package httpresponse

import (
	"encoding/json"
	"net/http"
)

func Response(w http.ResponseWriter, code int, data any) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

type responseErrType struct {
	Code  int    `json:"status_code"`
	Error string `json:"error"`
}

// Info: if type error is httperrors.HttpError use data stored in the error
// if the error type is unknown, the response will be with the code that was sent
func ResponseError(w http.ResponseWriter, code int, err error) {
	var res responseErrType
	if myErr, ok := err.(*HttpError); ok {
		res = responseErrType{
			Code:  myErr.Code,
			Error: myErr.Msg,
		}
	} else {
		res = responseErrType{
			Code:  code,
			Error: "unknown",
		}
	}

	w.WriteHeader(res.Code)
	json.NewEncoder(w).Encode(res)
}
