package test

import (
	"context"
	"testing"

	cdata "github.com/pip-services3-gox/pip-services3-commons-gox/data"
	"github.com/stretchr/testify/assert"
)

type DummyClientFixture struct {
	client IDummyClient
}

func NewDummyClientFixture(client IDummyClient) *DummyClientFixture {
	dcf := DummyClientFixture{client: client}
	return &dcf
}

func (c *DummyClientFixture) TestCrudOperations(t *testing.T) {
	ctx := context.Background()
	dummy1 := Dummy{Id: "", Key: "Key 1", Content: "Content 1"}
	dummy2 := Dummy{Id: "", Key: "Key 2", Content: "Content 2"}

	dummy1.Id = cdata.IdGenerator.NextLong()
	dummy2.Id = cdata.IdGenerator.NextLong()

	// Create one dummy
	dummy, err := c.client.CreateDummy(ctx, "", dummy1)
	assert.Nil(t, err)
	assert.NotNil(t, dummy)
	assert.Equal(t, dummy.Content, dummy1.Content)
	assert.Equal(t, dummy.Key, dummy1.Key)
	dummy1 = *dummy

	// Create another dummy
	dummy, err = c.client.CreateDummy(ctx, "", dummy2)
	assert.Nil(t, err)
	assert.NotNil(t, dummy)
	assert.Equal(t, dummy.Content, dummy2.Content)
	assert.Equal(t, dummy.Key, dummy2.Key)
	dummy2 = *dummy

	// Get all dummies
	dummies, err := c.client.GetDummies(ctx, "", cdata.NewEmptyFilterParams(), cdata.NewPagingParams(0, 15, false))
	assert.Nil(t, err)
	assert.NotNil(t, dummies)
	assert.True(t, len(dummies.Data) >= 2)

	// Update the dummy
	dummy1.Content = "Updated Content 1"
	dummy, err = c.client.UpdateDummy(ctx, "", dummy1)
	assert.Nil(t, err)
	assert.NotNil(t, dummy)
	assert.Equal(t, dummy.Content, "Updated Content 1")
	assert.Equal(t, dummy.Key, dummy1.Key)
	dummy1 = *dummy

	// Delete dummy
	dummy, err = c.client.DeleteDummy(ctx, "", dummy1.Id)
	assert.Nil(t, err)

	// Try to get delete dummy
	dummy, err = c.client.GetDummyById(ctx, "", dummy1.Id)
	assert.Nil(t, err)
	assert.Nil(t, dummy)

}
