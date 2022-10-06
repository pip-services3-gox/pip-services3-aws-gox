package container

import cref "github.com/pip-services3-gox/pip-services3-commons-gox/refer"

type ILambdaFunctionOverrides interface {
	cref.IReferenceable
	// Perform required registration steps.
	Register()
}
