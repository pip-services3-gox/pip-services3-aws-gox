package test

import (
	"context"

	cdata "github.com/pip-services3-gox/pip-services3-commons-gox/data"
)

type IDummyController interface {
	GetPageByFilter(ctx context.Context, correlationId string, filter *cdata.FilterParams, paging *cdata.PagingParams) (result *cdata.DataPage[Dummy], err error)
	GetOneById(ctx context.Context, correlationId string, id string) (result *Dummy, err error)
	Create(ctx context.Context, correlationId string, entity Dummy) (result *Dummy, err error)
	Update(ctx context.Context, correlationId string, entity Dummy) (result *Dummy, err error)
	DeleteById(ctx context.Context, correlationId string, id string) (result *Dummy, err error)
}
