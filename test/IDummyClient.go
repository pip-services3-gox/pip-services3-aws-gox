package test

import (
	"context"

	cdata "github.com/pip-services3-gox/pip-services3-commons-gox/data"
)

type IDummyClient interface {
	GetDummies(ctx context.Context, correlationId string, filter *cdata.FilterParams, paging *cdata.PagingParams) (result *cdata.DataPage[Dummy], err error)
	GetDummyById(ctx context.Context, correlationId string, dummyId string) (result *Dummy, err error)
	CreateDummy(ctx context.Context, correlationId string, dummy Dummy) (result *Dummy, err error)
	UpdateDummy(ctx context.Context, correlationId string, dummy Dummy) (result *Dummy, err error)
	DeleteDummy(ctx context.Context, correlationId string, dummyId string) (result *Dummy, err error)
}
