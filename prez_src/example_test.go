// -*- mode: fundamental; tab-width: 4; -*-
package main

import (
	"errors"
	"github.com/maxatome/go-testdeep/td"
	"testing"
)

type Person struct {
	ID       int64     `json:"id"`
	Name     string    `json:"name"`
	Age      int       `json:"age"`
	Children []*Person `json:"children"`
}

var alice = Person{23, "Alice", 20, nil}
var brian = Person{21, "Brian", 18, nil}
var bob = Person{10, "Bob", 41, []*Person{&alice, &brian}}

func GetPerson(name string) (*Person, error) {
	switch name {
	case "Bob":
		return &bob, nil
	case "Alice":
		return &alice, nil
	case "Brian":
		return &brian, nil
	}
	return nil, errors.New("User not found")
}

var personTests = []struct {
	name           string
	expectedErr    td.TestDeep
	expectedPerson td.TestDeep
}{
	{"Bob", nil, td.JSON(`{"name":"Bob","age":41,"id":NotZero(),"children":Len(2)}`)},
	{"Marcel", td.String("User not found"), td.Nil()},
	{"Alice", nil, td.SStruct(&Person{Name: "Alice", Age: 20}, td.StructFields{"ID": td.NotZero()})},
	{"Brian", nil, td.SStruct(&Person{Name: "Brian", Age: 18}, td.StructFields{"ID": td.NotZero()})},
}
                                                        === RUN   TestGetPerson
func TestGetPerson(t *testing.T) {                      === RUN   TestGetPerson/Bob
	assert := td.Assert(t)                              === RUN   TestGetPerson/Marcel
	for _, pt := range personTests {                    === RUN   TestGetPerson/Alice
		assert.Run(pt.name, func(assert *td.T) {        === RUN   TestGetPerson/Brian
			person, err := GetPerson(pt.name)           --- PASS: TestGetPerson (0.00s)
			assert.Cmp(err, pt.expectedErr)                 --- PASS: TestGetPerson/Bob (0.00s)
			assert.Cmp(person, pt.expectedPerson)           --- PASS: TestGetPerson/Marcel (0.00s)
		})                                                  --- PASS: TestGetPerson/Alice (0.00s)
	}                                                       --- PASS: TestGetPerson/Brian (0.00s)
}                                                       PASS
// end OMIT
