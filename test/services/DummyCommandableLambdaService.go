package test_services

import (
	"context"

	awsserv "github.com/pip-services3-gox/pip-services3-aws-gox/services"
	cref "github.com/pip-services3-gox/pip-services3-commons-gox/refer"
)

type DummyCommandableLambdaService struct {
	*awsserv.CommandableLambdaService
}

func NewDummyCommandableLambdaService() *DummyCommandableLambdaService {
	c := &DummyCommandableLambdaService{}
	c.CommandableLambdaService = awsserv.InheritCommandableLambdaService(c, "dummy")
	c.DependencyResolver.Put(context.Background(), "controller", cref.NewDescriptor("pip-services-dummies", "controller", "default", "*", "*"))
	return c
}
