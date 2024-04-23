package fxplus

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/modern-go/reflect2"
)

type Named interface {
	Name() string
}

type Healther interface {
	Health(context.Context) error
}

type HealtherFunc func(context.Context) error

func (h HealtherFunc) Health(ctx context.Context) error {
	return h(ctx)
}

type HealthReport struct {
	Name    string
	Error   string `json:",omitempty"`
	Success bool
}

func RespondHealth(w http.ResponseWriter, reports ...*HealthReport) error {
	code := 200
	for _, x := range reports {
		if !x.Success {
			code = 500
			break
		}
	}
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(reports)
}

func HealthCheck(ctx context.Context, x Healther) *HealthReport {
	var name, errString string
	if val, ok := x.(Named); ok {
		name = val.Name()
	} else {
		typ := reflect2.TypeOf(x)
		name = typ.String()
	}
	success := true
	err := x.Health(ctx)
	if err != nil {
		success = false
		errString = err.Error()
	}
	return &HealthReport{
		Name:    name,
		Error:   errString,
		Success: success,
	}
}
