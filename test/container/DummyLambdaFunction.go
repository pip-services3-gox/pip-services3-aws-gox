package test_container

import (
	"context"
	"encoding/json"

	awscont "github.com/pip-services3-gox/pip-services3-aws-gox/container"
	awstest "github.com/pip-services3-gox/pip-services3-aws-gox/test"
	cconv "github.com/pip-services3-gox/pip-services3-commons-gox/convert"
	cdata "github.com/pip-services3-gox/pip-services3-commons-gox/data"
	cref "github.com/pip-services3-gox/pip-services3-commons-gox/refer"
	cvalid "github.com/pip-services3-gox/pip-services3-commons-gox/validate"
)

type DummyLambdaFunction struct {
	*awscont.LambdaFunction
	controller awstest.IDummyController
}

func NewDummyLambdaFunction() *DummyLambdaFunction {
	c := &DummyLambdaFunction{}
	c.LambdaFunction = awscont.InheriteLambdaFunction(c, "dummy", "Dummy lambda function")

	c.DependencyResolver.Put(context.Background(), "controller", cref.NewDescriptor("pip-services-dummies", "controller", "default", "*", "*"))
	c.AddFactory(awstest.NewDummyFactory())
	return c
}

func (c *DummyLambdaFunction) SetReferences(ctx context.Context, references cref.IReferences) {
	c.LambdaFunction.SetReferences(ctx, references)
	depRes, depErr := c.DependencyResolver.GetOneRequired("controller")
	if depErr == nil && depRes != nil {
		c.controller = depRes.(awstest.IDummyController)
	}
}

func (c *DummyLambdaFunction) getPageByFilter(ctx context.Context, params map[string]any) (interface{}, error) {
	correlationId, _ := params["correlation_id"].(string)
	return c.controller.GetPageByFilter(
		ctx,
		correlationId,
		cdata.NewFilterParamsFromValue(params["filter"]),
		cdata.NewPagingParamsFromValue(params["paging"]),
	)
}

func (c *DummyLambdaFunction) getOneById(ctx context.Context, params map[string]any) (interface{}, error) {
	correlationId, _ := params["correlation_id"].(string)
	return c.controller.GetOneById(
		ctx,
		correlationId,
		params["dummy_id"].(string),
	)
}

func (c *DummyLambdaFunction) create(ctx context.Context, params map[string]any) (interface{}, error) {
	correlationId, _ := params["correlation_id"].(string)
	val, _ := json.Marshal(params["dummy"])
	var entity = awstest.Dummy{}
	json.Unmarshal(val, &entity)

	c.Logger().Debug(ctx, correlationId, "Create method called Dummy %v", entity)

	res, err := c.controller.Create(
		ctx,
		correlationId,
		entity,
	)

	c.Logger().Debug(ctx, correlationId, "Create method called Result: %v Err: %v", res, err)

	return res, err
}

func (c *DummyLambdaFunction) update(ctx context.Context, params map[string]any) (interface{}, error) {
	correlationId, _ := params["correlation_id"].(string)
	val, _ := json.Marshal(params["dummy"])
	var entity = awstest.Dummy{}
	json.Unmarshal(val, &entity)
	return c.controller.Update(
		ctx,
		correlationId,
		entity,
	)
}

func (c *DummyLambdaFunction) deleteById(ctx context.Context, params map[string]any) (interface{}, error) {
	correlationId, _ := params["correlation_id"].(string)

	c.Logger().Debug(ctx, correlationId, "DeleteById method called Id %v", params["dummy_id"].(string))

	res, err := c.controller.DeleteById(
		ctx,
		correlationId,
		params["dummy_id"].(string),
	)
	c.Logger().Debug(ctx, correlationId, "DeleteById method called Result: %v Err: %v", res, err)

	return res, err
}

func (c *DummyLambdaFunction) Register() {

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
