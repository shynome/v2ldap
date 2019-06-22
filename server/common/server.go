package common

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Resp json data
func Resp(w http.ResponseWriter, data interface{}, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if data == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	jsondata, err := json.Marshal(data)
	if err != nil {
		Resp(w, nil, err)
		return
	}
	headers := w.Header()
	headers["Content-Type"] = []string{"application/json"}
	w.Write(jsondata)
}

// ParseParamsFromReq parse json body
func ParseParamsFromReq(r *http.Request, v interface{}) (err error) {
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&v)
	if err != nil {
		err = fmt.Errorf("json parse fail: %v", err.Error())
	}
	return
}
