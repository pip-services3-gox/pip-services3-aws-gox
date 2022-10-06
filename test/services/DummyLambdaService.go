package test_services

import (
	"context"
	"encoding/json"

	awsserv "github.com/pip-services3-gox/pip-services3-aws-gox/services"
	awstest "github.com/pip-services3-gox/pip-services3-aws-gox/test"
	cconv "github.com/pip-services3-gox/pip-services3-commons-gox/convert"
	cdata "github.com/pip-services3-gox/pip-services3-commons-gox/data"
	cref "github.com/pip-services3-gox/pip-services3-commons-gox/refer"
	cvalid "github.com/pip-services3-gox/pip-services3-commons-gox/validate"
)

type DummyLambdaService struct {
	*awsserv.LambdaService
	controller awstest.IDummyController
}

func NewDummyLambdaService() *DummyLambdaService {
	c := &DummyLambdaService{}
	c.LambdaService = awsserv.InheritLambdaService(c, "dummy")

	c.DependencyResolver.Put(context.Background(), "controller", cref.NewDescriptor("pip-services-dummies", "controller", "default", "*", "*"))
	return c
}

func (c *DummyLambdaService) SetReferences(ctx context.Context, references cref.IReferences) {
	c.LambdaService.SetReferences(ctx, references)
	depRes, depErr := c.DependencyResolver.GetOneRequired("controller")
	if depErr == nil && depRes != nil {
		c.controller = depRes.(awstest.IDummyController)
	}
}

func (c *DummyLambdaService) getPageByFilter(ctx context.Context, params map[string]any) (interface{}, error) {
	correlationId, _ := params["correlation_id"].(string)
	return c.controller.GetPageByFilter(ctx,
		correlationId,
		cdata.NewFilterParamsFromValue(params["filter"]),
		cdata.NewPagingParamsFromValue(params["paging"]),
	)
}

func (c *DummyLambdaService) getOneById(ctx context.Context, params map[string]any) (interface{}, error) {
	correlationId, _ := params["correlation_id"].(string)
	return c.controller.GetOneById(ctx,
		correlationId,
		params["dummy_id"].(string),
	)
}

func (c *DummyLambdaService) create(ctx context.Context, params map[string]any) (interface{}, error) {
	correlationId, _ := params["correlation_id"].(string)
	val, _ := json.Marshal(params["dummy"])
	var entity awstest.Dummy
	json.Unmarshal(val, &entity)
	return c.controller.Create(ctx,
		correlationId,
		entity,
	)
}

func (c *DummyLambdaService) update(ctx context.Context, params map[string]any) (interface{}, error) {
	correlationId, _ := params["correlation_id"].(string)
	val, _ := json.Marshal(params["dummy"])
	var entity awstest.Dummy
	json.Unmarshal(val, &entity)
	return c.controller.Update(ctx,
		correlationId,
		entity,
	)
}

func (c *DummyLambdaService) deleteById(ctx context.Context, params map[string]any) (interface{}, error) {
	correlationId, _ := params["correlation_id"].(string)
	return c.controller.DeleteById(ctx,
		correlationId,
		params["dummy_id"].(string),
	)
}

func (c *DummyLambdaService) Register() {

	c.RegisterAction(
		"get_dummies",
		cvalid.NewObjectSchema().
			WithOptionalProperty("filter", cvalid.NewFilterParamsSchema()).
			WithOptionalProperty("paging", cvalid.NewPagingParamsSchema()).Schema,
		c.getPageByFilter)

	c.RegisterAction(
		"get_dummy_by_id",
		cvalid.NewObjectSchema().
			WithOptionalProperty("dummy_id", cconv.String).Schema,
		c.getOneById)

	c.RegisterAction(
		"create_dummy",
		cvalid.NewObjectSchema().
			WithRequiredProperty("dummy", awstest.NewDummySchema()).Schema,
		c.create)

	c.RegisterAction(
		"update_dummy",
		cvalid.NewObjectSchema().
			WithRequiredProperty("dummy", awstest.NewDummySchema()).Schema,
		c.update)

	c.RegisterAction(
		"delete_dummy",
		cvalid.NewObjectSchema().
			WithOptionalProperty("dummy_id", cconv.String).Schema,
		c.deleteById)
}
