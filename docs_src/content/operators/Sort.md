---
title: "Sort"
weight: 10
---

```go
func Sort(how any, expectedValue any) TestDeep
```

Sort is a [smuggler operator]({{% ref "operators#smuggler-operators" %}}). It takes an array, a slice or a
pointer on array/slice, it sorts it using *how* and compares the
sorted result to *expectedValue*. It can be seen as an alternative to
[`Bag`]({{% ref "Bag" %}}).

*how* can be:

- `nil` or a `float64`/`int` >= 0 for a generic ascending order;
- a `float64`/`int` < 0 for a generic descending order;
- a `string` specifying a fields-path (optionally prefixed by "+"
  or "-" for respectively an ascending or a descending order,
  defaulting to ascending one);
- a `[]string` containing a list of fields-paths (as above), second
  and next fields-paths are checked when the previous ones are equal;
- a function matching func(a, b T) `bool` signature and returning
  true if a is before b.


A fields-path, also used by [`Smuggle`]({{% ref "Smuggle" %}}) and [`Sorted`]({{% ref "Sorted" %}}) operators,
allows to access nested structs fields and maps & slices items. See
[`Smuggle`]({{% ref "Smuggle" %}}) for details on fields-path possibilities.

```go
type A struct{ props map[string]int }
p12 := A{props: map[string]int{"priority": 12}}
p23 := A{props: map[string]int{"priority": 23}}
p34 := A{props: map[string]int{"priority": 34}}
got := []A{p23, p12, p34}
td.Cmp(t, got, td.Sort("-props[priority]", []A{p34, p23, p12})) // succeeds
```

*how* can be a `float64` to allow Sort to be used in expected JSON of
[`JSON`]({{% ref "JSON" %}}), [`SubJSONOf`]({{% ref "SubJSONOf" %}}) & [`SuperJSONOf`]({{% ref "SuperJSONOf" %}}) operators:

```go
got := map[string][]string{"labels": {"c", "a", "b"}}
td.Cmp(t, got, td.JSON(`{ "labels": Sort(1, ["a", "b", "c"]) }`)) // succeeds
```

or using fields-path feature:

```go
type Person struct {
  Name string `json:"name"`
  Age  int    `json:"age"`
}
got := struct {
  People []Person `json:"people"`
}{
  People: []Person{
    {"Brian", 22},
    {"Bob", 19},
    {"Stephen", 19},
    {"Alice", 20},
    {"Marcel", 25},
  },
}
td.Cmp(t, got, td.JSON(`{
  "people": Sort("name", [ // sort by name ascending
    {"name": "Alice",   "age": 20},
    {"name": "Bob",     "age": 19},
    {"name": "Brian",   "age": 22},
    {"name": "Marcel",  "age": 25},
    {"name": "Stephen", "age": 19},
  ])
}`)) // succeeds
td.Cmp(t, got, td.JSON(`{
  "people": Sort([ "-age", "name" ], [ // sort by age desc, then by name asc
    {"name": "Marcel",  "age": 25},
    {"name": "Brian",   "age": 22},
    {"name": "Alice",   "age": 20},
    {"name": "Bob",     "age": 19},
    {"name": "Stephen", "age": 19},
  ])
}`)) // succeeds
```

> See also [`Sorted`]({{% ref "Sorted" %}}), [`Smuggle`]({{% ref "Smuggle" %}}) and [`Bag`]({{% ref "Bag" %}}).


> See also [<i class='fas fa-book'></i> Sort godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Sort).

### Examples

{{%expand "Basic example" %}}```go
	t := &testing.T{}

	got := []int{-1, 1, 2, -3, 3, -2, 0}

	// Generic ascending order (≥0 or nil)
	ok := td.Cmp(t, got, td.Sort(1, []int{-3, -2, -1, 0, 1, 2, 3}))
	fmt.Println("asc order:", ok)

	ok = td.Cmp(t, got, td.Sort(0, []int{-3, -2, -1, 0, 1, 2, 3}))
	fmt.Println("asc order:", ok)

	ok = td.Cmp(t, got, td.Sort(nil, []int{-3, -2, -1, 0, 1, 2, 3}))
	fmt.Println("asc order:", ok)

	// Generic descending order (< 0)
	ok = td.Cmp(t, got, td.Sort(-1, []int{3, 2, 1, 0, -1, -2, -3}))
	fmt.Println("desc order:", ok)

	evenHigher := func(a, b int) bool {
		if (a%2 == 0) != (b%2 == 0) {
			return a%2 != 0
		}
		return a < b
	}
	ok = td.Cmp(t, got, td.Sort(evenHigher, []int{-3, -1, 1, 3, -2, 0, 2}))
	fmt.Println("even higher order:", ok)

	// Output:
	// asc order: true
	// asc order: true
	// asc order: true
	// desc order: true
	// even higher order: true

```{{% /expand%}}
{{%expand "Fields_path example" %}}```go
	t := &testing.T{}

	type Person struct {
		Name string
		Age  int
	}

	brian := Person{Name: "Brian", Age: 22}
	bob := Person{Name: "Bob", Age: 19}
	stephen := Person{Name: "Stephen", Age: 19}
	alice := Person{Name: "Alice", Age: 20}
	marcel := Person{Name: "Marcel", Age: 25}
	got := []Person{brian, bob, stephen, alice, marcel}

	ok := td.Cmp(t, got,
		td.Sort("Name", []Person{alice, bob, brian, marcel, stephen}))
	fmt.Println("by name asc:", ok)

	ok = td.Cmp(t, got,
		td.Sort("-Name", []Person{stephen, marcel, brian, bob, alice}))
	fmt.Println("by name desc:", ok)

	ok = td.Cmp(t, got,
		td.Sort([]string{"-Age", "Name"}, []Person{marcel, brian, alice, bob, stephen}))
	fmt.Println("by age desc, then by name asc:", ok)

	type A struct{ props map[string]int }
	p12 := A{props: map[string]int{"priority": 12}}
	p23 := A{props: map[string]int{"priority": 23}}
	p34 := A{props: map[string]int{"priority": 34}}
	got2 := []A{p23, p12, p34}
	ok = td.Cmp(t, got2, td.Sort(`-props[priority]`, []A{p34, p23, p12}))
	fmt.Println("by priority desc:", ok)

	ok = td.Cmp(t, got2, td.Sort(`props.priority`, []A{p12, p23, p34}))
	fmt.Println("by priority asc:", ok)

	// Output:
	// by name asc: true
	// by name desc: true
	// by age desc, then by name asc: true
	// by priority desc: true
	// by priority asc: true

```{{% /expand%}}
## CmpSort shortcut

```go
func CmpSort(t TestingT, got, how , expectedValue any, args ...any) bool
```

CmpSort is a shortcut for:

```go
td.Cmp(t, got, td.Sort(how, expectedValue), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

If *t* is a [`*T`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T) then its Config field is inherited.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpSort godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpSort).

### Examples

{{%expand "Basic example" %}}```go
	t := &testing.T{}

	got := []int{-1, 1, 2, -3, 3, -2, 0}

	// Generic ascending order (≥0 or nil)
	ok := td.CmpSort(t, got, 1, []int{-3, -2, -1, 0, 1, 2, 3})
	fmt.Println("asc order:", ok)

	ok = td.CmpSort(t, got, 0, []int{-3, -2, -1, 0, 1, 2, 3})
	fmt.Println("asc order:", ok)

	ok = td.CmpSort(t, got, nil, []int{-3, -2, -1, 0, 1, 2, 3})
	fmt.Println("asc order:", ok)

	// Generic descending order (< 0)
	ok = td.CmpSort(t, got, -1, []int{3, 2, 1, 0, -1, -2, -3})
	fmt.Println("desc order:", ok)

	evenHigher := func(a, b int) bool {
		if (a%2 == 0) != (b%2 == 0) {
			return a%2 != 0
		}
		return a < b
	}
	ok = td.CmpSort(t, got, evenHigher, []int{-3, -1, 1, 3, -2, 0, 2})
	fmt.Println("even higher order:", ok)

	// Output:
	// asc order: true
	// asc order: true
	// asc order: true
	// desc order: true
	// even higher order: true

```{{% /expand%}}
{{%expand "Fields_path example" %}}```go
	t := &testing.T{}

	type Person struct {
		Name string
		Age  int
	}

	brian := Person{Name: "Brian", Age: 22}
	bob := Person{Name: "Bob", Age: 19}
	stephen := Person{Name: "Stephen", Age: 19}
	alice := Person{Name: "Alice", Age: 20}
	marcel := Person{Name: "Marcel", Age: 25}
	got := []Person{brian, bob, stephen, alice, marcel}

	ok := td.CmpSort(t, got, "Name", []Person{alice, bob, brian, marcel, stephen})
	fmt.Println("by name asc:", ok)

	ok = td.CmpSort(t, got, "-Name", []Person{stephen, marcel, brian, bob, alice})
	fmt.Println("by name desc:", ok)

	ok = td.CmpSort(t, got, []string{"-Age", "Name"}, []Person{marcel, brian, alice, bob, stephen})
	fmt.Println("by age desc, then by name asc:", ok)

	type A struct{ props map[string]int }
	p12 := A{props: map[string]int{"priority": 12}}
	p23 := A{props: map[string]int{"priority": 23}}
	p34 := A{props: map[string]int{"priority": 34}}
	got2 := []A{p23, p12, p34}
	ok = td.CmpSort(t, got2, `-props[priority]`, []A{p34, p23, p12})
	fmt.Println("by priority desc:", ok)

	ok = td.CmpSort(t, got2, `props.priority`, []A{p12, p23, p34})
	fmt.Println("by priority asc:", ok)

	// Output:
	// by name asc: true
	// by name desc: true
	// by age desc, then by name asc: true
	// by priority desc: true
	// by priority asc: true

```{{% /expand%}}
## T.Sort shortcut

```go
func (t *T) Sort(got, how , expectedValue any, args ...any) bool
```

Sort is a shortcut for:

```go
t.Cmp(got, td.Sort(how, expectedValue), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.Sort godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.Sort).

### Examples

{{%expand "Basic example" %}}```go
	t := td.NewT(&testing.T{})

	got := []int{-1, 1, 2, -3, 3, -2, 0}

	// Generic ascending order (≥0 or nil)
	ok := t.Sort(got, 1, []int{-3, -2, -1, 0, 1, 2, 3})
	fmt.Println("asc order:", ok)

	ok = t.Sort(got, 0, []int{-3, -2, -1, 0, 1, 2, 3})
	fmt.Println("asc order:", ok)

	ok = t.Sort(got, nil, []int{-3, -2, -1, 0, 1, 2, 3})
	fmt.Println("asc order:", ok)

	// Generic descending order (< 0)
	ok = t.Sort(got, -1, []int{3, 2, 1, 0, -1, -2, -3})
	fmt.Println("desc order:", ok)

	evenHigher := func(a, b int) bool {
		if (a%2 == 0) != (b%2 == 0) {
			return a%2 != 0
		}
		return a < b
	}
	ok = t.Sort(got, evenHigher, []int{-3, -1, 1, 3, -2, 0, 2})
	fmt.Println("even higher order:", ok)

	// Output:
	// asc order: true
	// asc order: true
	// asc order: true
	// desc order: true
	// even higher order: true

```{{% /expand%}}
{{%expand "Fields_path example" %}}```go
	t := td.NewT(&testing.T{})

	type Person struct {
		Name string
		Age  int
	}

	brian := Person{Name: "Brian", Age: 22}
	bob := Person{Name: "Bob", Age: 19}
	stephen := Person{Name: "Stephen", Age: 19}
	alice := Person{Name: "Alice", Age: 20}
	marcel := Person{Name: "Marcel", Age: 25}
	got := []Person{brian, bob, stephen, alice, marcel}

	ok := t.Sort(got, "Name", []Person{alice, bob, brian, marcel, stephen})
	fmt.Println("by name asc:", ok)

	ok = t.Sort(got, "-Name", []Person{stephen, marcel, brian, bob, alice})
	fmt.Println("by name desc:", ok)

	ok = t.Sort(got, []string{"-Age", "Name"}, []Person{marcel, brian, alice, bob, stephen})
	fmt.Println("by age desc, then by name asc:", ok)

	type A struct{ props map[string]int }
	p12 := A{props: map[string]int{"priority": 12}}
	p23 := A{props: map[string]int{"priority": 23}}
	p34 := A{props: map[string]int{"priority": 34}}
	got2 := []A{p23, p12, p34}
	ok = t.Sort(got2, `-props[priority]`, []A{p34, p23, p12})
	fmt.Println("by priority desc:", ok)

	ok = t.Sort(got2, `props.priority`, []A{p12, p23, p34})
	fmt.Println("by priority asc:", ok)

	// Output:
	// by name asc: true
	// by name desc: true
	// by age desc, then by name asc: true
	// by priority desc: true
	// by priority asc: true

```{{% /expand%}}
