---
title: "Struct"
weight: 10
---

```go
func Struct(model any, expectedFields ...StructFields) TestDeep
```

Struct operator compares the contents of a struct or a pointer on a
struct against the non-zero values of *model* (if `any`) and the
values of *expectedFields*. See [`SStruct`]({{% ref "SStruct" %}}) to compares against zero
fields without specifying them in *expectedFields*.

*model* must be the same type as compared data. If the expected type
is anonymous or private, *model* can be `nil`. In this case it is
considered lazy and determined each time the operator is involved
in a match, see below.

*expectedFields* can be omitted, if no zero entries are expected
and no [TestDeep operators]({{% ref "operators" %}}) are involved. If *expectedFields*
contains more than one item, all items are merged before their use,
from left to right.

```go
td.Cmp(t, got, td.Struct(
  Person{
    Name: "John Doe",
  },
  td.StructFields{
    "Children": 4,
  },
  td.StructFields{
    "Age":      td.Between(40, 45),
    "Children": 0, // overwrite 4
  }),
)
```

It is an [`error`](https://pkg.go.dev/builtin#error) to set a non-zero field in *model* AND to set the
same field in *expectedFields*, as in such cases the Struct
operator does not know if the user wants to override the non-zero
*model* field value or if it is an [`error`](https://pkg.go.dev/builtin#error). To explicitly override a
non-zero *model* in *expectedFields*, just prefix its name with a
">" (followed by some optional spaces), as in:

```go
td.Cmp(t, got, td.Struct(
  Person{
    Name:     "John Doe",
    Age:      23,
    Children: 4,
  },
  td.StructFields{
    "> Age":     td.Between(40, 45),
    ">Children": 0, // spaces after ">" are optional
  }),
)
```

*expectedFields* can also contain regexps or shell patterns to
match multiple fields not explicitly listed in *model* and in
*expectedFields*. Regexps are prefixed by "=\~" or "!\~" to
respectively match or don't-match. Shell patterns are prefixed by "="
or "!" to respectively match or don't-match.

```go
td.Cmp(t, got, td.Struct(
  Person{
    Name: "John Doe",
  },
  td.StructFields{
    "=*At":     td.Lte(time.Now()), // matches CreatedAt & UpdatedAt fields using shell pattern
    "=~^[a-z]": td.Ignore(),        // explicitly ignore private fields using a regexp
  }),
)
```

When several patterns can match a same field, it is advised to tell
go-testdeep in which order patterns should be tested, as once a
pattern matches a field, the other patterns are ignored for this
field. To do so, each pattern can be prefixed by a number, as in:

```go
td.Cmp(t, got, td.Struct(
  Person{
    Name: "John Doe",
  },
  td.StructFields{
    "1=*At":     td.Lte(time.Now()),
    "2=~^[a-z]": td.NotNil(),
  }),
)
```

This way, "*At" shell pattern is always used before "^[a-z]"
regexp, so if a field "createdAt" exists it is tested against
time.Now() and never against [`NotNil`]({{% ref "NotNil" %}}). A pattern without a
prefix number is the same as specifying "0" as prefix.

To make it clearer, some spaces can be added, as well as bigger
numbers used:

```go
td.Cmp(t, got, td.Struct(
  Person{
    Name: "John Doe",
  },
  td.StructFields{
    " 900 =  *At":    td.Lte(time.Now()),
    "2000 =~ ^[a-z]": td.NotNil(),
  }),
)
```

The following example combines all possibilities:

```go
td.Cmp(t, got, td.Struct(
  Person{
    NickName: "Joe",
  },
  td.StructFields{
    "Firstname":               td.Any("John", "Johnny"),
    "1 =  *[nN]ame":           td.NotEmpty(), // matches LastName, lastname, …
    "2 !  [A-Z]*":             td.NotZero(),  // matches all private fields
    "3 =~ ^(Crea|Upda)tedAt$": td.Gte(time.Now()),
    "4 !~ ^(Dogs|Children)$":  td.Zero(),   // matches all remaining fields except Dogs and Children
    "5 =~ .":                  td.NotNil(), // matches all remaining fields (same as "5 = *")
  }),
)
```

If the expected type is private to the current package, it cannot
be passed as *model*. To overcome this limitation, *model* can be `nil`,
it is then considered as lazy. This way, the *model* is automatically
set during each match to the same type (still requiring struct or
struct pointer) of the compared data. Similarly, testing an
anonymous struct can be boring as all fields have to be re-declared
to define *model*. A `nil` *model* avoids that:

```go
got := struct {
  name string
  age  int
}{"Bob", 42}
td.Cmp(t, got, td.Struct(nil, td.StructFields{"age": td.Between(40, 42)}))
```

During a match, all expected fields must be found to
succeed. Non-expected fields (and so zero *model* fields) are
ignored.

[`TypeBehind`]({{% ref "operators#typebehind-method" %}}) method returns the [`reflect.Type`](https://pkg.go.dev/reflect#Type) of *model*.

> See also [`SStruct`]({{% ref "SStruct" %}}).


> See also [<i class='fas fa-book'></i> Struct godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Struct).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	type Person struct {
		Name        string
		Age         int
		NumChildren int
	}

	got := Person{
		Name:        "Foobar",
		Age:         42,
		NumChildren: 3,
	}

	// As NumChildren is zero in Struct() call, it is not checked
	ok := td.Cmp(t, got,
		td.Struct(Person{Name: "Foobar"}, td.StructFields{
			"Age": td.Between(40, 50),
		}),
		"checks %v is the right Person")
	fmt.Println("Foobar is between 40 & 50:", ok)

	// Model can be empty
	ok = td.Cmp(t, got,
		td.Struct(Person{}, td.StructFields{
			"Name":        "Foobar",
			"Age":         td.Between(40, 50),
			"NumChildren": td.Not(0),
		}),
		"checks %v is the right Person")
	fmt.Println("Foobar has some children:", ok)

	// Works with pointers too
	ok = td.Cmp(t, &got,
		td.Struct(&Person{}, td.StructFields{
			"Name":        "Foobar",
			"Age":         td.Between(40, 50),
			"NumChildren": td.Not(0),
		}),
		"checks %v is the right Person")
	fmt.Println("Foobar has some children (using pointer):", ok)

	// Model does not need to be instanciated
	ok = td.Cmp(t, &got,
		td.Struct((*Person)(nil), td.StructFields{
			"Name":        "Foobar",
			"Age":         td.Between(40, 50),
			"NumChildren": td.Not(0),
		}),
		"checks %v is the right Person")
	fmt.Println("Foobar has some children (using nil model):", ok)

	// Output:
	// Foobar is between 40 & 50: true
	// Foobar has some children: true
	// Foobar has some children (using pointer): true
	// Foobar has some children (using nil model): true

```{{% /expand%}}
{{%expand "Overwrite_model example" %}}```go
	t := &testing.T{}

	type Person struct {
		Name        string
		Age         int
		NumChildren int
	}

	got := Person{
		Name:        "Foobar",
		Age:         42,
		NumChildren: 3,
	}

	ok := td.Cmp(t, got,
		td.Struct(
			Person{
				Name: "Foobar",
				Age:  53,
			},
			td.StructFields{
				">Age":        td.Between(40, 50), // ">" to overwrite Age:53 in model
				"NumChildren": td.Gt(2),
			}),
		"checks %v is the right Person")
	fmt.Println("Foobar is between 40 & 50:", ok)

	ok = td.Cmp(t, got,
		td.Struct(
			Person{
				Name: "Foobar",
				Age:  53,
			},
			td.StructFields{
				"> Age":       td.Between(40, 50), // same, ">" can be followed by spaces
				"NumChildren": td.Gt(2),
			}),
		"checks %v is the right Person")
	fmt.Println("Foobar is between 40 & 50:", ok)

	// Output:
	// Foobar is between 40 & 50: true
	// Foobar is between 40 & 50: true

```{{% /expand%}}
{{%expand "Patterns example" %}}```go
	t := &testing.T{}

	type Person struct {
		Firstname string
		Lastname  string
		Surname   string
		Nickname  string
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt *time.Time
	}

	now := time.Now()
	got := Person{
		Firstname: "Maxime",
		Lastname:  "Foo",
		Surname:   "Max",
		Nickname:  "max",
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil, // not deleted yet
	}

	ok := td.Cmp(t, got,
		td.Struct(Person{Lastname: "Foo"}, td.StructFields{
			`DeletedAt`: nil,
			`=  *name`:  td.Re(`^(?i)max`),  // shell pattern, matches all names except Lastname as in model
			`=~ At\z`:   td.Lte(time.Now()), // regexp, matches CreatedAt & UpdatedAt
		}),
		"mix shell & regexp patterns")
	fmt.Println("Patterns match only remaining fields:", ok)

	ok = td.Cmp(t, got,
		td.Struct(Person{Lastname: "Foo"}, td.StructFields{
			`DeletedAt`:  nil,
			`1 =  *name`: td.Re(`^(?i)max`),  // shell pattern, matches all names except Lastname as in model
			`2 =~ At\z`:  td.Lte(time.Now()), // regexp, matches CreatedAt & UpdatedAt
		}),
		"ordered patterns")
	fmt.Println("Ordered patterns match only remaining fields:", ok)

	// Output:
	// Patterns match only remaining fields: true
	// Ordered patterns match only remaining fields: true

```{{% /expand%}}
{{%expand "Lazy_model example" %}}```go
	t := &testing.T{}

	got := struct {
		name string
		age  int
	}{
		name: "Foobar",
		age:  42,
	}

	ok := td.Cmp(t, got, td.Struct(nil, td.StructFields{
		"name": "Foobar",
		"age":  td.Between(40, 45),
	}))
	fmt.Println("Lazy model:", ok)

	ok = td.Cmp(t, got, td.Struct(nil, td.StructFields{
		"name": "Foobar",
		"zip":  666,
	}))
	fmt.Println("Lazy model with unknown field:", ok)

	// Output:
	// Lazy model: true
	// Lazy model with unknown field: false

```{{% /expand%}}
## CmpStruct shortcut

```go
func CmpStruct(t TestingT, got, model any, expectedFields StructFields, args ...any) bool
```

CmpStruct is a shortcut for:

```go
td.Cmp(t, got, td.Struct(model, expectedFields), args...)
```

See above for details.

[`Struct`]({{% ref "Struct" %}}) optional parameter *expectedFields* is here mandatory.
`nil` value should be passed to mimic its absence in
original [`Struct`]({{% ref "Struct" %}}) call.

Returns true if the test is OK, false if it fails.

If *t* is a [`*T`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T) then its Config field is inherited.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpStruct godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpStruct).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	type Person struct {
		Name        string
		Age         int
		NumChildren int
	}

	got := Person{
		Name:        "Foobar",
		Age:         42,
		NumChildren: 3,
	}

	// As NumChildren is zero in Struct() call, it is not checked
	ok := td.CmpStruct(t, got, Person{Name: "Foobar"}, td.StructFields{
		"Age": td.Between(40, 50),
	},
		"checks %v is the right Person")
	fmt.Println("Foobar is between 40 & 50:", ok)

	// Model can be empty
	ok = td.CmpStruct(t, got, Person{}, td.StructFields{
		"Name":        "Foobar",
		"Age":         td.Between(40, 50),
		"NumChildren": td.Not(0),
	},
		"checks %v is the right Person")
	fmt.Println("Foobar has some children:", ok)

	// Works with pointers too
	ok = td.CmpStruct(t, &got, &Person{}, td.StructFields{
		"Name":        "Foobar",
		"Age":         td.Between(40, 50),
		"NumChildren": td.Not(0),
	},
		"checks %v is the right Person")
	fmt.Println("Foobar has some children (using pointer):", ok)

	// Model does not need to be instanciated
	ok = td.CmpStruct(t, &got, (*Person)(nil), td.StructFields{
		"Name":        "Foobar",
		"Age":         td.Between(40, 50),
		"NumChildren": td.Not(0),
	},
		"checks %v is the right Person")
	fmt.Println("Foobar has some children (using nil model):", ok)

	// Output:
	// Foobar is between 40 & 50: true
	// Foobar has some children: true
	// Foobar has some children (using pointer): true
	// Foobar has some children (using nil model): true

```{{% /expand%}}
{{%expand "Overwrite_model example" %}}```go
	t := &testing.T{}

	type Person struct {
		Name        string
		Age         int
		NumChildren int
	}

	got := Person{
		Name:        "Foobar",
		Age:         42,
		NumChildren: 3,
	}

	ok := td.CmpStruct(t, got, Person{
		Name: "Foobar",
		Age:  53,
	}, td.StructFields{
		">Age":        td.Between(40, 50), // ">" to overwrite Age:53 in model
		"NumChildren": td.Gt(2),
	},
		"checks %v is the right Person")
	fmt.Println("Foobar is between 40 & 50:", ok)

	ok = td.CmpStruct(t, got, Person{
		Name: "Foobar",
		Age:  53,
	}, td.StructFields{
		"> Age":       td.Between(40, 50), // same, ">" can be followed by spaces
		"NumChildren": td.Gt(2),
	},
		"checks %v is the right Person")
	fmt.Println("Foobar is between 40 & 50:", ok)

	// Output:
	// Foobar is between 40 & 50: true
	// Foobar is between 40 & 50: true

```{{% /expand%}}
{{%expand "Patterns example" %}}```go
	t := &testing.T{}

	type Person struct {
		Firstname string
		Lastname  string
		Surname   string
		Nickname  string
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt *time.Time
	}

	now := time.Now()
	got := Person{
		Firstname: "Maxime",
		Lastname:  "Foo",
		Surname:   "Max",
		Nickname:  "max",
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil, // not deleted yet
	}

	ok := td.CmpStruct(t, got, Person{Lastname: "Foo"}, td.StructFields{
		`DeletedAt`: nil,
		`=  *name`:  td.Re(`^(?i)max`),  // shell pattern, matches all names except Lastname as in model
		`=~ At\z`:   td.Lte(time.Now()), // regexp, matches CreatedAt & UpdatedAt
	},
		"mix shell & regexp patterns")
	fmt.Println("Patterns match only remaining fields:", ok)

	ok = td.CmpStruct(t, got, Person{Lastname: "Foo"}, td.StructFields{
		`DeletedAt`:  nil,
		`1 =  *name`: td.Re(`^(?i)max`),  // shell pattern, matches all names except Lastname as in model
		`2 =~ At\z`:  td.Lte(time.Now()), // regexp, matches CreatedAt & UpdatedAt
	},
		"ordered patterns")
	fmt.Println("Ordered patterns match only remaining fields:", ok)

	// Output:
	// Patterns match only remaining fields: true
	// Ordered patterns match only remaining fields: true

```{{% /expand%}}
{{%expand "Lazy_model example" %}}```go
	t := &testing.T{}

	got := struct {
		name string
		age  int
	}{
		name: "Foobar",
		age:  42,
	}

	ok := td.CmpStruct(t, got, nil, td.StructFields{
		"name": "Foobar",
		"age":  td.Between(40, 45),
	})
	fmt.Println("Lazy model:", ok)

	ok = td.CmpStruct(t, got, nil, td.StructFields{
		"name": "Foobar",
		"zip":  666,
	})
	fmt.Println("Lazy model with unknown field:", ok)

	// Output:
	// Lazy model: true
	// Lazy model with unknown field: false

```{{% /expand%}}
## T.Struct shortcut

```go
func (t *T) Struct(got, model any, expectedFields StructFields, args ...any) bool
```

Struct is a shortcut for:

```go
t.Cmp(got, td.Struct(model, expectedFields), args...)
```

See above for details.

[`Struct`]({{% ref "Struct" %}}) optional parameter *expectedFields* is here mandatory.
`nil` value should be passed to mimic its absence in
original [`Struct`]({{% ref "Struct" %}}) call.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.Struct godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.Struct).

### Examples

{{%expand "Base example" %}}```go
	t := td.NewT(&testing.T{})

	type Person struct {
		Name        string
		Age         int
		NumChildren int
	}

	got := Person{
		Name:        "Foobar",
		Age:         42,
		NumChildren: 3,
	}

	// As NumChildren is zero in Struct() call, it is not checked
	ok := t.Struct(got, Person{Name: "Foobar"}, td.StructFields{
		"Age": td.Between(40, 50),
	},
		"checks %v is the right Person")
	fmt.Println("Foobar is between 40 & 50:", ok)

	// Model can be empty
	ok = t.Struct(got, Person{}, td.StructFields{
		"Name":        "Foobar",
		"Age":         td.Between(40, 50),
		"NumChildren": td.Not(0),
	},
		"checks %v is the right Person")
	fmt.Println("Foobar has some children:", ok)

	// Works with pointers too
	ok = t.Struct(&got, &Person{}, td.StructFields{
		"Name":        "Foobar",
		"Age":         td.Between(40, 50),
		"NumChildren": td.Not(0),
	},
		"checks %v is the right Person")
	fmt.Println("Foobar has some children (using pointer):", ok)

	// Model does not need to be instanciated
	ok = t.Struct(&got, (*Person)(nil), td.StructFields{
		"Name":        "Foobar",
		"Age":         td.Between(40, 50),
		"NumChildren": td.Not(0),
	},
		"checks %v is the right Person")
	fmt.Println("Foobar has some children (using nil model):", ok)

	// Output:
	// Foobar is between 40 & 50: true
	// Foobar has some children: true
	// Foobar has some children (using pointer): true
	// Foobar has some children (using nil model): true

```{{% /expand%}}
{{%expand "Overwrite_model example" %}}```go
	t := td.NewT(&testing.T{})

	type Person struct {
		Name        string
		Age         int
		NumChildren int
	}

	got := Person{
		Name:        "Foobar",
		Age:         42,
		NumChildren: 3,
	}

	ok := t.Struct(got, Person{
		Name: "Foobar",
		Age:  53,
	}, td.StructFields{
		">Age":        td.Between(40, 50), // ">" to overwrite Age:53 in model
		"NumChildren": td.Gt(2),
	},
		"checks %v is the right Person")
	fmt.Println("Foobar is between 40 & 50:", ok)

	ok = t.Struct(got, Person{
		Name: "Foobar",
		Age:  53,
	}, td.StructFields{
		"> Age":       td.Between(40, 50), // same, ">" can be followed by spaces
		"NumChildren": td.Gt(2),
	},
		"checks %v is the right Person")
	fmt.Println("Foobar is between 40 & 50:", ok)

	// Output:
	// Foobar is between 40 & 50: true
	// Foobar is between 40 & 50: true

```{{% /expand%}}
{{%expand "Patterns example" %}}```go
	t := td.NewT(&testing.T{})

	type Person struct {
		Firstname string
		Lastname  string
		Surname   string
		Nickname  string
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt *time.Time
	}

	now := time.Now()
	got := Person{
		Firstname: "Maxime",
		Lastname:  "Foo",
		Surname:   "Max",
		Nickname:  "max",
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil, // not deleted yet
	}

	ok := t.Struct(got, Person{Lastname: "Foo"}, td.StructFields{
		`DeletedAt`: nil,
		`=  *name`:  td.Re(`^(?i)max`),  // shell pattern, matches all names except Lastname as in model
		`=~ At\z`:   td.Lte(time.Now()), // regexp, matches CreatedAt & UpdatedAt
	},
		"mix shell & regexp patterns")
	fmt.Println("Patterns match only remaining fields:", ok)

	ok = t.Struct(got, Person{Lastname: "Foo"}, td.StructFields{
		`DeletedAt`:  nil,
		`1 =  *name`: td.Re(`^(?i)max`),  // shell pattern, matches all names except Lastname as in model
		`2 =~ At\z`:  td.Lte(time.Now()), // regexp, matches CreatedAt & UpdatedAt
	},
		"ordered patterns")
	fmt.Println("Ordered patterns match only remaining fields:", ok)

	// Output:
	// Patterns match only remaining fields: true
	// Ordered patterns match only remaining fields: true

```{{% /expand%}}
{{%expand "Lazy_model example" %}}```go
	t := td.NewT(&testing.T{})

	got := struct {
		name string
		age  int
	}{
		name: "Foobar",
		age:  42,
	}

	ok := t.Struct(got, nil, td.StructFields{
		"name": "Foobar",
		"age":  td.Between(40, 45),
	})
	fmt.Println("Lazy model:", ok)

	ok = t.Struct(got, nil, td.StructFields{
		"name": "Foobar",
		"zip":  666,
	})
	fmt.Println("Lazy model with unknown field:", ok)

	// Output:
	// Lazy model: true
	// Lazy model with unknown field: false

```{{% /expand%}}
