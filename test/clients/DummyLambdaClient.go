package test

import (
	"context"

	"github.com/pip-services3-gox/pip-services3-aws-gox/clients"
	awsclient "github.com/pip-services3-gox/pip-services3-aws-gox/clients"
	awstest "github.com/pip-services3-gox/pip-services3-aws-gox/test"
	cdata "github.com/pip-services3-gox/pip-services3-commons-gox/data"
)

type DummyLambdaClient struct {
	*awsclient.LambdaClient
}

func NewDummyLambdaClient() *DummyLambdaClient {
	c := &DummyLambdaClient{
		LambdaClient: awsclient.NewLambdaClient(),
	}
	return c
}
func (c *DummyLambdaClient) GetDummies(ctx context.Context, correlationId string, filter *cdata.FilterParams,
	paging *cdata.PagingParams) (result *cdata.DataPage[awstest.Dummy], err error) {
	timing := c.Instrument(ctx, correlationId, "dummy.get_dummies")

	params := cdata.NewEmptyAnyValueMap()
	params.SetAsObject("filter", filter)
	params.SetAsObject("paging", paging)

	calValue, calErr := c.Call(ctx, "get_dummies", correlationId, params.Value())
	if calErr != nil {
		return nil, calErr
	}

	defer timing.EndTiming(ctx, err)

	return awsclient.HandleLambdaResponse[*cdata.DataPage[awstest.Dummy]](calValue)
}

func (c *DummyLambdaClient) GetDummyById(ctx context.Context, correlationId string, dummyId string) (result *awstest.Dummy, err error) {
	timing := c.Instrument(ctx, correlationId, "dummy.get_one_by_id")
	params := cdata.NewEmptyAnyValueMap()
	params.SetAsObject("dummy_id", dummyId)

	calValue, calErr := c.Call(ctx, "get_dummy_by_id", correlationId, params.Value())

	if calErr != nil {
		return nil, calErr
	}

	defer timing.EndTiming(ctx, err)

	return clients.HandleLambdaResponse[*awstest.Dummy](calValue)
}

func (c *DummyLambdaClient) CreateDummy(ctx context.Context, correlationId string, dummy awstest.Dummy) (result *awstest.Dummy, err error) {
	timing := c.Instrument(ctx, correlationId, "dummy.create_dummy")
	params := cdata.NewEmptyAnyValueMap()
	params.SetAsObject("dummy", dummy)

	calValue, calErr := c.Call(ctx, "create_dummy", correlationId, params.Value())
	if calErr != nil {
		return nil, calErr
	}

	defer timing.EndTiming(ctx, err)

	return clients.HandleLambdaResponse[*awstest.Dummy](calValue)
}

func (c *DummyLambdaClient) UpdateDummy(ctx context.Context, correlationId string, dummy awstest.Dummy) (result *awstest.Dummy, err error) {
	timing := c.Instrument(ctx, correlationId, "dummy.update_dummy")
	params := cdata.NewEmptyAnyValueMap()
	params.SetAsObject("dummy", dummy)

	calValue, calErr := c.Call(ctx, "update_dummy", correlationId, params.Value())
	if calErr != nil {
		return nil, calErr
	}

	defer timing.EndTiming(ctx, err)
	return clients.HandleLambdaResponse[*awstest.Dummy](calValue)
}

func (c *DummyLambdaClient) DeleteDummy(ctx context.Context, correlationId string, dummyId string) (result *awstest.Dummy, err error) {
	timing := c.Instrument(ctx, correlationId, "dummy.delete_dummy")
	params := cdata.NewEmptyAnyValueMap()
	params.SetAsObject("dummy_id", dummyId)
	calValue, calErr := c.Call(ctx, "delete_dummy", correlationId, params.Value())
	if calErr != nil {
		return nil, calErr
	}

	defer timing.EndTiming(ctx, err)
	return clients.HandleLambdaResponse[*awstest.Dummy](calValue)
}
