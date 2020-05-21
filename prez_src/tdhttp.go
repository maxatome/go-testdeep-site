package main

import (
	"testing"

	"github.com/maxatome/go-testdeep/helpers/tdhttp"
	"github.com/maxatome/go-testdeep/td"
)

func TestMyApi(t *testing.T) {
	ta := tdhttp.NewTestAPI(t, myAPI)
	// ta-end OMIT

	ta.Get("/person/Bob", "Accept", "application/json").
		CmpStatus(htp.StatusOK).
		CmpHeader(td.ContainsKey("X-Custom-Header")).
		CmpJSONBody(td.JSON(`{"id": $1, "name": "Bob", "age":  26}`, td.NotZero()))

	if !ta.Failed() {
		t.Log("Good job pal!")
	}
}

func TestMyApiAnchor(t *testing.T) {
	ta := tdhttp.NewTestAPI(t, myAPI)

	ta.Get("/person/Bob", "Accept", "application/json").
		CmpStatus(htp.StatusOK).
		CmpHeader(td.ContainsKey("X-Custom-Header")).
		CmpJSONBody(Person{
			ID:   ta.A(td.NotZero(), int64(0)).(int64), // HL
			Name: "Bob",
			Age:  ta.A(td.Between(40, 45)).(int), // HL
		})
}
