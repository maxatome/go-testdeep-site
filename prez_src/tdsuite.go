package main

import (
	"testing"

	"github.com/maxatome/go-testdeep/helpers/tdsuite"
	"github.com/maxatome/go-testdeep/td"
)

func TestPerson(t *testing.T) {
	tdsuite.Run(t, &PersonSuite{ // entrypoint of the suite
		db: InitDB(), // a DB handler probably used in each tests
	})
}

type PersonSuite struct{ db MyDBHandler }

func (ps *PersonSuite) TestGet(assert *td.T) {
	// …
}

func (ps *PersonSuite) TestPost(assert, require *td.T) {
	// …
}

// very-end OMIT

// compose-begin OMIT
func TestAnother(t *testing.T) { //                    🔗 https://goplay.tools/snippet/5cbM9eHbx33
	tdsuite.Run(t, &AnotherSuite{}) // entrypoint of the suite
}

// BaseSuite is the base test suite used by all tests suite using the DB.
type BaseSuite struct{ db MyDBHandler }

func (bs *BaseSuite) Setup(t *td.T) (err error) {
	bs.db, err = InitDB()
	return
}

func (bs *BaseSuite) Destroy(t *td.T) error {
	return bs.db.Exec(`TRUNCATE x, y, z CASCADE`)
}

// AnotherSuite is the final test suite blah blah blah…
type AnotherSuite struct{ BaseSuite }

func (as *AnotherSuite) TestGet(assert, require *td.T) {
	res, err := as.db.Query(`SELECT 42`)
	require.CmpNoError(err)
	assert.Cmp(res, 42)
}

// compose-end OMIT
