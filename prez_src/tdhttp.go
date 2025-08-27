package main

import (
	"testing"

	"github.com/maxatome/go-testdeep/helpers/tdhttp"
	"github.com/maxatome/go-testdeep/td"
)

func TestMyApi(t *testing.T) {
	ta := tdhttp.NewTestAPI(t, myAPI)
	// ta-end OMIT

	var id int64
	ta.Name("Retrieve a person").
		Get("/person/Bob", "Accept", "application/json").
		CmpStatus(http.StatusOK).
		CmpHeader(td.ContainsKey("X-Custom-Header")).
		CmpJSONBody(td.JSON(`{"id": $1, "name": "Bob", "age":  26}`, td.Catch(&id, td.NotZero())))

	t.Logf("Did the test succeeded? %t, ID of Bob is %d", !ta.Failed(), id)
}

func TestMyApiAnchor(t *testing.T) {
	ta := tdhttp.NewTestAPI(t, myAPI)

	var id int64
	ta.Get("/person/Bob", "Accept", "application/json").
		CmpStatus(http.StatusOK).
		CmpJSONBody(Person{
			ID:   ta.A(td.Catch(&id, td.NotZero())).(int64), // TestAPI.A method // HL
			Name: "Bob",
			Age:  td.A[int](ta.T(), td.Between(25, 30)), // td.A generic version // HL
		})
}

func TestMyApiDumpIfFailure(t *testing.T) {
	ta := tdhttp.NewTestAPI(t, myAPI)

	var id int64
	ta.Name("Person creation").
		PostJSON("/person", PersonNew{Name: "Bob"}).
		CmpStatus(http.StatusCreated).
		CmpJSONBody(td.JSON(`
			{
			  "id": $1,
			  "name": "Bob",
			  "created_at": $2
			}`,
			td.Catch(&id, td.NotZero()), // catch just created ID
			td.Gte(ta.SentAt()),         // check that created_at is ≥ request sent date
		)).
		OrDumpResponse() // if some test fails, the response is dumped // HL
}
