package test

import (
	"context"

	ccomand "github.com/pip-services3-gox/pip-services3-commons-gox/commands"
	cdata "github.com/pip-services3-gox/pip-services3-commons-gox/data"
)

type DummyController struct {
	commandSet *DummyCommandSet
	entities   []Dummy
}

func NewDummyController() *DummyController {
	dc := DummyController{}
	dc.entities = make([]Dummy, 0)
	return &dc
}

func (c *DummyController) GetCommandSet() *ccomand.CommandSet {
	if c.commandSet == nil {
		c.commandSet = NewDummyCommandSet(c)
	}
	return c.commandSet.CommandSet
}

func (c *DummyController) GetPageByFilter(ctx context.Context, correlationId string, filter *cdata.FilterParams, paging *cdata.PagingParams) (items *cdata.DataPage[Dummy], err error) {

	if filter == nil {
		filter = cdata.NewEmptyFilterParams()
	}
	var key string = filter.GetAsString("key")

	if paging == nil {
		paging = cdata.NewEmptyPagingParams()
	}
	var skip int64 = paging.GetSkip(0)
	var take int64 = paging.GetTake(100)

	var result []Dummy
	for i := 0; i < len(c.entities); i++ {
		var entity Dummy = c.entities[i]
		if key != "" && key != entity.Key {
			continue
		}

		skip--
		if skip >= 0 {
			continue
		}

		take--
		if take < 0 {
			break
		}

		result = append(result, entity)
	}
	var total int64 = (int64)(len(result))
	return cdata.NewDataPage(result, int(total)), nil
}

func (c *DummyController) GetOneById(ctx context.Context, correlationId string, id string) (result *Dummy, err error) {
	for i := 0; i < len(c.entities); i++ {
		var entity Dummy = c.entities[i]
		if id == entity.Id {
			return &entity, nil
		}
	}
	return nil, nil
}

func (c *DummyController) Create(ctx context.Context, correlationId string, entity Dummy) (result *Dummy, err error) {
	if entity.Id == "" {
		entity.Id = cdata.IdGenerator.NextLong()
	}
	c.entities = append(c.entities, entity)
	return &entity, nil
}

func (c *DummyController) Update(ctx context.Context, correlationId string, newEntity Dummy) (result *Dummy, err error) {
	for index := 0; index < len(c.entities); index++ {
		var entity Dummy = c.entities[index]
		if entity.Id == newEntity.Id {
			c.entities[index] = newEntity
			return &newEntity, nil

		}
	}
	return nil, nil
}

func (c *DummyController) DeleteById(ctx context.Context, correlationId string, id string) (result *Dummy, err error) {
	var entity Dummy

	for i := 0; i < len(c.entities); {
		entity = c.entities[i]
		if entity.Id == id {
			if i == len(c.entities)-1 {
				c.entities = c.entities[:i]
			} else {
				c.entities = append(c.entities[:i], c.entities[i+1:]...)
			}
			return &entity, nil
		}
		i++
	}
	return nil, nil
}
