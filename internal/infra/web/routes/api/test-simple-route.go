package routes

import (
	"encoding/json"
	"net/http"
)

type TestSimpleRoute struct {
}

func NewTestSimpleRoute() *TestSimpleRoute {
	return &TestSimpleRoute{}
}

type output struct {
	Message string
}

func (cr *TestSimpleRoute) Handler(w http.ResponseWriter, r *http.Request) {
	output := &output{
		Message: "Success",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}
