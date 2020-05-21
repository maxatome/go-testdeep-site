package main

func checkPerson(t *testing.T, p *Person, name string, age int) bool {
	t.Helper()
	td.Cmp(t, p.Name, name)
	td.Cmp(t, p.Age, age)
}

func TestPerson(t *testing.T) {
	checkPerson(t, GetPerson("Bob"), "Bob", 37)
}

func checkPerson(name string, age int) TestDeep {
	return td.Struct(Person{Name: name, Age: age})
}

func TestPerson(t *testing.T) {
	td.Cmp(t, GetPerson("Bob"), checkPerson("Bob", 37))
}

func CheckDateGte(t time.Time, catch *time.Time) td.TestDeep {
	op := td.Gte(t.Truncate(time.Millisecond))
	if catch != nil {
		op = td.Catch(catch, op)
	}
	return td.All(
		td.HasSuffix("Z"),
		td.Smuggle(func(s string) (time.Time, error) {
			t, err := time.Parse(time.RFC3339Nano, s)
			if err == nil && t.IsZero() {
				err = fmt.Errorf("zero time")
			}
			return t, err
		}, op))
}

func TestCreateArticle(t *testing.T) {
	type Article struct {
		ID        int64     `json:"id"`
		Code      string    `json:"code"`
		CreatedAt time.Time `json:"created_at"`
	}

	beforeCreation := time.Now()
	var createdAt time.Time
	td.Cmp(t, CreateArticle("Car"),
		td.JSON(`{"id": $^NotZero, "code": "Car", "created_at": $1}`,
			CheckDateGte(beforeCreation, &createdAt)))

	// If the test succeeds, then "created_at" value is well a RFC3339
	// datetime in UTC timezone and its value is directly exploitable as
	// time.Time thanks to createdAt variable
}

func TestCustom(t *testing.T) {
	// Code-begin OMIT
	td.Cmp(t, gotTime,
		td.Code(func(date time.Time) bool {
			return date.Year() == 2018
		}))
	// Code-end OMIT

	// Smuggle-begin OMIT
	td.Cmp(t, gotTime,
		td.Smuggle(func(date time.Time) int { return date.Year() },
			td.Between(2010, 2020)),
	)
	// Smuggle-end OMIT
}
