package schema_test

import (
	"math/big"
	"testing"

	"github.com/hsfzxjy/imbed/schema"
	"github.com/stretchr/testify/assert"
)

func TestReflect(t *testing.T) {
	type Foo struct {
		int64    `imbed:""`
		*big.Rat `imbed:""`
		string   `imbed:""`
		bool     `imbed:""`

		unannotatedField int
		UnannotatedField int
	}

	var foo = Foo{
		int64:            42,
		Rat:              big.NewRat(314, 100),
		string:           "foo_bar_baz",
		bool:             true,
		unannotatedField: 114,
		UnannotatedField: 514,
	}
	schema1 := schema.Struct(&foo,
		schema.F("int64", &foo.int64, schema.Int()),
		schema.F("Rat", &foo.Rat, schema.Rat()),
		schema.F("string", &foo.string, schema.String()),
		schema.F("bool", &foo.bool, schema.Bool()),
	).Build()

	store := schema.NewStore()
	schema2, err := schema.Register[Foo](store)
	assert.ErrorIs(t, nil, err)

	{
		schema3, err := store.Get(&foo)
		assert.ErrorIs(t, nil, err)
		assert.True(t, schema3 == schema2)

		_, err = store.Get(foo)
		assert.ErrorContains(t, err, "Store.Get() expect *T, got ")

		_, err = store.Get(&struct{}{})
		assert.ErrorContains(t, err, "no schema registered for struct {}")
	}

	encoded1, err := schema.EncodeBytes(schema1, &foo)
	assert.ErrorIs(t, nil, err)

	encoded2, err := schema.EncodeBytesAny(schema2, &foo)
	assert.ErrorIs(t, nil, err)

	assert.Equal(t, encoded1, encoded2)
}

func TestReflectBad(t *testing.T) {
	store := schema.NewStore()
	_, err := schema.Register[int](store)
	assert.ErrorContains(t, err, "expect a struct, got int")
	_, err = schema.Register[*struct{ int }](store)
	assert.ErrorContains(t, err, "expect a struct, got *struct { int }")
	_, err = schema.Register[struct {
		foo struct{ int } `imbed:""`
	}](store)
	assert.ErrorContains(t, err, "field \"foo\": cannot create schema for struct { int }")
}
