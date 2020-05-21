package main

type Person struct {
	ID       int64     `json:"id"`
	Name     string    `json:"name"`
	Age      int       `json:"age"`
	Children []*Person `json:"children"`
}

func Test(t *testing.T) {
	// nested1-begin OMIT
	got := GetPerson("Bob")

	td.Cmp(t, got.ID, td.NotZero())
	td.Cmp(t, got.Name, "Bob")
	td.Cmp(t, got.Age, td.Between(40, 45))
	if td.Cmp(t, got.Children, td.Len(2)) {
		// Alice
		td.Cmp(t, got.Children[0].ID, td.NotZero())
		td.Cmp(t, got.Children[0].Name, "Alice")
		td.Cmp(t, got.Children[0].Age, 20)
		td.Cmp(t, got.Children[0].Children, td.Len(0))
		// Brian
		td.Cmp(t, got.Children[1].ID, td.NotZero())
		td.Cmp(t, got.Children[1].Name, "Brian")
		td.Cmp(t, got.Children[1].Age, 18)
		td.Cmp(t, got.Children[1].Children, td.Len(0))
	}
	// nested1-end OMIT

	// nested2-begin OMIT
	got := GetPerson("Bob")
	td.Cmp(t, got,
		td.SStruct(
			Person{Name: "Bob"},
			td.StructFields{
				"ID":  td.NotZero(),
				"Age": td.Between(40, 45),
				"Children": td.Bag(
					td.SStruct(&Person{Name: "Alice", Age: 20},
						td.StructFields{"ID": td.NotZero()}),
					td.SStruct(&Person{Name: "Brian", Age: 18},
						td.StructFields{"ID": td.NotZero()}),
				),
			},
		))
	// nested2-end OMIT

	// nested3-begin OMIT
	td.Cmp(t, got, td.JSON(`
    {
        "id":       $^NotZero, // ← simple operator (8 others are eligible)
        "name":     "Bob"
        "age":      $1,        // ← placeholder (could be "$1" or $BobAge, see JSON operator doc)
        "children": [
            {
                "id":       $^NotZero,
                "name":     "Alice",
                "age":      20,
                "children": null,
            },
            {
                "id":       $^NotZero,
                "name":     "Brian",
                "age":      18,
                "children": null,
            }
        ]
    }`,
		td.Between(40, 45),
	))
	// nested3-end OMIT

	// nested4-begin OMIT
	tdt := td.NewT(t)
	tdt.Cmp(got,
		Person{
			ID:   tdt.A(td.NotZero(), int64(0)).(int64), // HL
			Name: "Bob",
			Age:  tdt.A(td.Between(40, 45)).(int), // HL
			Children: []*Person{
				{
					ID:   tdt.A(td.NotZero(), int64(0)).(int64), // HL
					Name: "Alice",
					Age:  20,
				},
				{
					ID:   tdt.A(td.NotZero(), int64(0)).(int64), // HL
					Name: "Brian",
					Age:  18,
				},
			},
		})
	// nested4-end OMIT
}
