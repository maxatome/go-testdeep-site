package main

type T struct {
	testing.TB               // implemented by *testing.T
	Config     ContextConfig // defaults to td.DefaultContextConfig
}

func TestVals(t *testing.T) {
	got := GetPerson("Bob")
	td.Cmp(t, got.Age, td.Between(40, 45))
	td.Cmp(t, got.Children, td.Len(2))
}

// tdt-begin OMIT
func TestVals(t *testing.T) {
	assert := td.Assert(t) // HL

	got := GetPerson("Bob")
	assert.Cmp(got.Age, td.Between(40, 45))
	assert.Cmp(got.Children, td.Len(2))
}
