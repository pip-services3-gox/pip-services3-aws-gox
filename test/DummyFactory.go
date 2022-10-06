package test

import (
	cref "github.com/pip-services3-gox/pip-services3-commons-gox/refer"
	cbuild "github.com/pip-services3-gox/pip-services3-components-gox/build"
)

type DummyFactory struct {
	*cbuild.Factory
	Descriptor                 *cref.Descriptor
	ControllerDescriptor       *cref.Descriptor
	LambdaServiceDescriptor    *cref.Descriptor
	CmdLambdaServiceDescriptor *cref.Descriptor
}

func NewDummyFactory() *DummyFactory {

	c := DummyFactory{
		Factory:              cbuild.NewFactory(),
		Descriptor:           cref.NewDescriptor("pip-services-dummies", "factory", "default", "default", "1.0"),
		ControllerDescriptor: cref.NewDescriptor("pip-services-dummies", "controller", "default", "*", "1.0"),
	}

	c.RegisterType(c.ControllerDescriptor, NewDummyController)
	return &c
}
