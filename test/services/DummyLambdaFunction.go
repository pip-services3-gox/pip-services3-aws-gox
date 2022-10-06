package test_services

import (
	awscont "github.com/pip-services3-gox/pip-services3-aws-gox/container"
	awstest "github.com/pip-services3-gox/pip-services3-aws-gox/test"
)

type DummyLambdaFunction struct {
	*awscont.LambdaFunction
}

func NewDummyLambdaFunction() *DummyLambdaFunction {
	c := &DummyLambdaFunction{}
	c.LambdaFunction = awscont.InheriteLambdaFunction(c, "dummy", "Dummy lambda function")
	c.AddFactory(awstest.NewDummyFactory())
	c.AddFactory(NewDummyLambdaServiceFactory())
	return c
}

func (c *DummyLambdaFunction) Register() {}
