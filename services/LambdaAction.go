package services

import (
	"context"

	cvalid "github.com/pip-services3-gox/pip-services3-commons-gox/validate"
)

type LambdaAction struct {

	// Command to call the action
	Cmd string

	// Schema to validate action parameters
	Schema *cvalid.Schema

	// Action to be executed
	Action func(ctx context.Context, params map[string]any) (any, error)
}
