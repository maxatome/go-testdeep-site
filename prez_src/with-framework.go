package main

import (
	"testing"

	"github.com/maxatome/go-testdeep/td"
)

func TestGetPerson(t *testing.T) {
	person, err := GetPerson("Bob")
	if td.CmpNoError(t, err, "GetPerson does not return an error") {
		td.Cmp(t, person, Person{Name: "Bob", Age: 42}, "GetPerson returns Bob")
	}
}

// without-test-name OMIT
func TestGetPerson(t *testing.T) {
	person, err := GetPerson("Bob")
	if td.CmpNoError(t, err) {
		td.Cmp(t, person, Person{Name: "Bob", Age: 42})
	}
}
