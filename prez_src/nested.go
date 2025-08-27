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
	td.Cmp(t, GetPerson("Bob"),
		td.SStruct(Person{Name: "Bob"}, // HL
			td.StructFields{
				"ID":  td.NotZero(),
				"Age": td.Between(40, 45),
				"Children": td.Bag(
					td.SStruct(&Person{Name: "Alice", Age: 20}, // HL
						td.StructFields{"ID": td.NotZero()}),
					td.SStruct(&Person{Name: "Brian", Age: 18}, // HL
						td.StructFields{"ID": td.NotZero()}),
				),
			},
		))
	// nested2-end OMIT

	// nested3-begin OMIT
	td.Cmp(t, GetPerson("Bob"), td.JSON(`
		{
			"id":   $1,     // ← placeholder (could be "$1" or $BobAge, see JSON operator doc)
			"name": "Bob",
			"age":  Between(40, 45),     // yes, most operators are embedable
			"children": [
				{
					"id":       NotZero,
					"name":     "Alice",
					"age":      20,
					"children": Empty, /* null is "empty" */
				},
				{
					"id":       NotZero,
					"name":     "Brian",
					"age":      18,
					"children": Nil,
				}
			]
		}`,
		td.Catch(&bobID, td.NotZero()), // $1 catches the ID of Bob on the fly and sets bobID var
	))
	// nested3-end OMIT

	// nested3-bag-begin OMIT
	td.Cmp(t, GetPerson("Bob"), td.JSON(`
		{
			"id":   $1,     // ← placeholder (could be "$1" or $BobAge, see JSON operator doc)
			"name": "Bob",
			"age":  Between(40, 45),     // yes, most operators are embedable
			"children": Bag(             // ← Bag HERE
				{
					"id":       NotZero,
					"name":     "Brian",
					"age":      18,
					"children": null,
				},
				{
					"id":       NotZero(),
					"name":     "Alice",
					"age":      20,
					"children": Empty, /* null is "empty" */
				},
			)
		}`,
		td.Catch(&bobID, td.NotZero()), // $1 catches the ID of Bob on the fly and sets bobID var
	))
	// nested3-bag-end OMIT

	// nested4-begin OMIT
	assert := td.Assert(t)
	assert.Cmp(GetPerson("Bob"),
		Person{
			ID:   td.A[int64](assert, td.Catch(&bobID, td.NotZero())), // HL
			Name: "Bob",
			Age:  td.A[int](assert, td.Between(40, 45)), // HL
			Children: []*Person{
				{
					ID:   td.A[int64](assert, td.NotZero()), // HL
					Name: "Alice",
					Age:  20,
				},
				{
					ID:   td.A[int64](assert, td.NotZero()), // HL
					Name: "Brian",
					Age:  18,
				},
			},
		})
	// nested4-end OMIT
}
