package fxplus

import (
	"context"
	"encoding/json"
	"net/http"

	"go.uber.org/fx"
)

type Named interface {
	Name() string
}

type Healther interface {
	Health(context.Context) error
}

type HealthReport struct {
	Name    string
	Error   string
	Success bool
}

type HealtherGroup struct {
	fx.In
	Healthers []Healther `group:"fxplus"`
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

func HealthCheck(ctx context.Context, xs Healther) *HealthReport {
	var name, errString string
	if val, ok := xs.(Named); ok {
		name = val.Name()
	}
	success := true
	err := xs.Health(ctx)
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
