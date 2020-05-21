package main

import (
	"testing"

	"github.com/maxatome/go-testdeep/td"
)

func TestExample(t *testing.T) {
	assert, require := td.AssertRequire(t)

	person, err := GetPerson("Bob")

	require.CmpNoError(err) // exits test if it fails

	assert.RootName("PERSON").
		Cmp(person, &Person{
			ID:       assert.A(td.NotZero(), int64(0)).(int64),
			Name:     "Bob",
			Age:      assert.A(td.Between(40, 45)).(int),
			Children: assert.A(td.Len(2), ([]*Person)(nil)).([]*Person),
		})
}
