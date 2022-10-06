package test_container

import (
	"context"

	awscont "github.com/pip-services3-gox/pip-services3-aws-gox/container"
	awstest "github.com/pip-services3-gox/pip-services3-aws-gox/test"
	cref "github.com/pip-services3-gox/pip-services3-commons-gox/refer"
)

type DummyCommandableLambdaFunction struct {
	*awscont.CommandableLambdaFunction
}

func NewDummyCommandableLambdaFunction() *DummyCommandableLambdaFunction {
	c := &DummyCommandableLambdaFunction{}
	c.CommandableLambdaFunction = awscont.NewCommandableLambdaFunction("dummy", "Dummy lambda function")

	c.DependencyResolver.Put(context.Background(), "controller", cref.NewDescriptor("pip-services-dummies", "controller", "default", "*", "*"))
	c.AddFactory(awstest.NewDummyFactory())
	return c
}
