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
			ID:       td.A[int64](assert, td.NotZero()), // HL
			Name:     "Bob",
			Age:      td.A[int](assert, td.Between(40, 45)), // HL
			Children: td.A[[]*Person](assert, td.Len(2)),    // HL
		})
}
