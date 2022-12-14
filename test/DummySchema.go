package test

import (
	cconv "github.com/pip-services3-gox/pip-services3-commons-gox/convert"
	cvalid "github.com/pip-services3-gox/pip-services3-commons-gox/validate"
)

type DummySchema struct {
	*cvalid.ObjectSchema
}

func NewDummySchema() *DummySchema {
	ds := DummySchema{}
	ds.ObjectSchema = cvalid.NewObjectSchema()
	ds.WithOptionalProperty("id", cconv.String)
	ds.WithRequiredProperty("key", cconv.String)
	ds.WithOptionalProperty("content", cconv.String)
	return &ds
}
