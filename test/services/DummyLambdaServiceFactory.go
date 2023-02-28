package test_services

import (
	cref "github.com/pip-services3-gox/pip-services3-commons-gox/refer"
	cbuild "github.com/pip-services3-gox/pip-services3-components-gox/build"
)

type DummyLambdaServiceFactory struct {
	*cbuild.Factory
	Descriptor                 *cref.Descriptor
	ControllerDescriptor       *cref.Descriptor
	LambdaServiceDescriptor    *cref.Descriptor
	CmdLambdaServiceDescriptor *cref.Descriptor
}

func NewDummyLambdaServiceFactory() *DummyLambdaServiceFactory {

	c := DummyLambdaServiceFactory{
		Factory:                    cbuild.NewFactory(),
		Descriptor:                 cref.NewDescriptor("pip-services-dummies", "factory", "default", "default", "1.0"),
		LambdaServiceDescriptor:    cref.NewDescriptor("pip-services-dummies", "service", "awslambda", "*", "1.0"),
		CmdLambdaServiceDescriptor: cref.NewDescriptor("pip-services-dummies", "service", "commandable-awslambda", "*", "1.0"),
	}

	c.RegisterType(c.LambdaServiceDescriptor, NewDummyLambdaService)
	c.RegisterType(c.CmdLambdaServiceDescriptor, NewDummyCommandableLambdaService)
	return &c
}
