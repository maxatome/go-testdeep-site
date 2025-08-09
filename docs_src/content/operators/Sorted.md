---
title: "Sorted"
weight: 10
---

```go
func Sorted(how ...any) TestDeep
```

Sorted operator checks that data is an array, a slice or a pointer
on array/slice, and it is well sorted as *how* tells it should be.

*how*... can be:

- empty to check a generic ascending order;
- `nil` or a `float64`/`int` >= 0 to check a generic ascending order;
- a `float64`/`int` < 0 to check a generic descending order;
- strings specifying fields-paths (each optionally prefixed by "+"
  or "-" for respectively checking an ascending or a descending order,
  defaulting to ascending one);
- a function matching func(a, b T) `bool` signature and returning
  true if a is before b.


A fields-path, also used by [`Smuggle`]({{% ref "Smuggle" %}}) and [`Sort`]({{% ref "Sort" %}}) operators,
allows to access nested structs fields and maps & slices items. See
[`Smuggle`]({{% ref "Smuggle" %}}) for details on fields-path possibilities.

```go
type A struct{ props map[string]int }
p12 := A{props: map[string]int{"priority": 12}}
p23 := A{props: map[string]int{"priority": 23}}
p34 := A{props: map[string]int{"priority": 34}}
got := []A{p34, p23, p12}
td.Cmp(t, got, td.Sorted("-props[priority]")) // succeeds
```

*how* can be a `float64` to allow Sort to be used in expected JSON of
[`JSON`]({{% ref "JSON" %}}), [`SubJSONOf`]({{% ref "SubJSONOf" %}}) & [`SuperJSONOf`]({{% ref "SuperJSONOf" %}}) operators:

```go
got := map[string][]string{"labels": {"a", "b", "c"}}
td.Cmp(t, got, td.JSON(`{ "labels": Sorted }`)) // succeeds
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
    {"Marcel", 25},
    {"Brian", 22},
    {"Alice", 20},
    {"Bob", 19},
    {"Stephen", 19},
  },
}
// sorted by age desc, then by name asc
td.Cmp(t, got, td.JSON(`{ "people": Sorted("-age", "name") }`)) // succeeds
```

> See also [`Sort`]({{% ref "Sort" %}}).


> See also [<i class='fas fa-book'></i> Sorted godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Sorted).

### Examples

{{%expand "Basic example" %}}```go
	t := &testing.T{}

	got := []int{-3, -2, -1, 0, 1, 2, 3}

	ok := td.Cmp(t, got, td.Sorted())
	fmt.Println("is asc order (default):", ok)

	ok = td.Cmp(t, got, td.Sorted(1))
	fmt.Println("is asc order (1):", ok)

	ok = td.Cmp(t, got, td.Sorted(0))
	fmt.Println("is asc order (0):", ok)

	ok = td.Cmp(t, got, td.Sorted(nil))
	fmt.Println("is asc order (nil):", ok)

	ok = td.Cmp(t, got, td.Sorted(-1))
	fmt.Println("is desc order:", ok)

	got = []int{-3, -1, 1, 3, -2, 0, 2}
	evenHigher := func(a, b int) bool {
		if (a%2 == 0) != (b%2 == 0) {
			return a%2 != 0
		}
		return a < b
	}
	ok = td.Cmp(t, got, td.Sorted(evenHigher))
	fmt.Println("is even higher order:", ok)

	// Output:
	// is asc order (default): true
	// is asc order (1): true
	// is asc order (0): true
	// is asc order (nil): true
	// is desc order: false
	// is even higher order: true

```{{% /expand%}}
{{%expand "Fields_path example" %}}```go
	t := &testing.T{}

	type Person struct {
		Name string
		Age  int
	}

	alice := Person{Name: "Alice", Age: 20}
	bob := Person{Name: "Bob", Age: 19}
	brian := Person{Name: "Brian", Age: 22}
	marcel := Person{Name: "Marcel", Age: 25}
	stephen := Person{Name: "Stephen", Age: 19}

	got := []Person{alice, bob, brian, marcel, stephen}
	ok := td.Cmp(t, got, td.Sorted("Name"))
	fmt.Println("is sorted by name asc:", ok)

	got = []Person{stephen, marcel, brian, bob, alice}
	ok = td.Cmp(t, got, td.Sorted("-Name"))
	fmt.Println("is sorted by name desc:", ok)

	got = []Person{marcel, brian, alice, bob, stephen}
	ok = td.Cmp(t, got, td.Sorted("-Age", "Name"))
	fmt.Println("is sorted by age desc, then by name asc 1:", ok)

	got = []Person{marcel, brian, alice, stephen, bob}
	ok = td.Cmp(t, got, td.Sorted("-Age", "Name"))
	fmt.Println("is sorted by age desc, then by name asc 2:", ok)

	type A struct{ props map[string]int }
	p12 := A{props: map[string]int{"priority": 12}}
	p23 := A{props: map[string]int{"priority": 23}}
	p34 := A{props: map[string]int{"priority": 34}}
	got2 := []A{p34, p23, p12}
	ok = td.Cmp(t, got2, td.Sorted(`-props[priority]`))
	fmt.Println("is sorted by priority desc:", ok)

	ok = td.Cmp(t, got2, td.Sorted(`props.priority`))
	fmt.Println("is sorted by priority asc:", ok)

	// Output:
	// is sorted by name asc: true
	// is sorted by name desc: true
	// is sorted by age desc, then by name asc 1: true
	// is sorted by age desc, then by name asc 2: false
	// is sorted by priority desc: true
	// is sorted by priority asc: false

```{{% /expand%}}
## CmpSorted shortcut

```go
func CmpSorted(t TestingT, got, how any, args ...any) bool
```

CmpSorted is a shortcut for:

```go
td.Cmp(t, got, td.Sorted(how), args...)
```

See above for details.

[`Sorted`]({{% ref "Sorted" %}}) optional parameter *how* is here mandatory.
`nil` value should be passed to mimic its absence in
original [`Sorted`]({{% ref "Sorted" %}}) call.

Returns true if the test is OK, false if it fails.

If *t* is a [`*T`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T) then its Config field is inherited.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpSorted godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpSorted).

### Examples

{{%expand "Basic example" %}}```go
	t := &testing.T{}

	got := []int{-3, -2, -1, 0, 1, 2, 3}

	ok := td.CmpSorted(t, got, nil)
	fmt.Println("is asc order (default):", ok)

	ok = td.CmpSorted(t, got, 1)
	fmt.Println("is asc order (1):", ok)

	ok = td.CmpSorted(t, got, 0)
	fmt.Println("is asc order (0):", ok)

	ok = td.CmpSorted(t, got, nil)
	fmt.Println("is asc order (nil):", ok)

	ok = td.CmpSorted(t, got, -1)
	fmt.Println("is desc order:", ok)

	got = []int{-3, -1, 1, 3, -2, 0, 2}
	evenHigher := func(a, b int) bool {
		if (a%2 == 0) != (b%2 == 0) {
			return a%2 != 0
		}
		return a < b
	}
	ok = td.CmpSorted(t, got, evenHigher)
	fmt.Println("is even higher order:", ok)

	// Output:
	// is asc order (default): true
	// is asc order (1): true
	// is asc order (0): true
	// is asc order (nil): true
	// is desc order: false
	// is even higher order: true

```{{% /expand%}}
{{%expand "Fields_path example" %}}```go
	t := &testing.T{}

	type Person struct {
		Name string
		Age  int
	}

	alice := Person{Name: "Alice", Age: 20}
	bob := Person{Name: "Bob", Age: 19}
	brian := Person{Name: "Brian", Age: 22}
	marcel := Person{Name: "Marcel", Age: 25}
	stephen := Person{Name: "Stephen", Age: 19}

	got := []Person{alice, bob, brian, marcel, stephen}
	ok := td.CmpSorted(t, got, "Name")
	fmt.Println("is sorted by name asc:", ok)

	got = []Person{stephen, marcel, brian, bob, alice}
	ok = td.CmpSorted(t, got, "-Name")
	fmt.Println("is sorted by name desc:", ok)

	got = []Person{marcel, brian, alice, bob, stephen}
	ok = td.CmpSorted(t, got, []string{"-Age", "Name"})
	fmt.Println("is sorted by age desc, then by name asc 1:", ok)

	got = []Person{marcel, brian, alice, stephen, bob}
	ok = td.CmpSorted(t, got, []string{"-Age", "Name"})
	fmt.Println("is sorted by age desc, then by name asc 2:", ok)

	type A struct{ props map[string]int }
	p12 := A{props: map[string]int{"priority": 12}}
	p23 := A{props: map[string]int{"priority": 23}}
	p34 := A{props: map[string]int{"priority": 34}}
	got2 := []A{p34, p23, p12}
	ok = td.CmpSorted(t, got2, `-props[priority]`)
	fmt.Println("is sorted by priority desc:", ok)

	ok = td.CmpSorted(t, got2, `props.priority`)
	fmt.Println("is sorted by priority asc:", ok)

	// Output:
	// is sorted by name asc: true
	// is sorted by name desc: true
	// is sorted by age desc, then by name asc 1: true
	// is sorted by age desc, then by name asc 2: false
	// is sorted by priority desc: true
	// is sorted by priority asc: false

```{{% /expand%}}
## T.Sorted shortcut

```go
func (t *T) Sorted(got, how any, args ...any) bool
```

Sorted is a shortcut for:

```go
t.Cmp(got, td.Sorted(how), args...)
```

See above for details.

[`Sorted`]({{% ref "Sorted" %}}) optional parameter *how* is here mandatory.
`nil` value should be passed to mimic its absence in
original [`Sorted`]({{% ref "Sorted" %}}) call.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.Sorted godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.Sorted).

### Examples

{{%expand "Basic example" %}}```go
	t := td.NewT(&testing.T{})

	got := []int{-3, -2, -1, 0, 1, 2, 3}

	ok := t.Sorted(got, nil)
	fmt.Println("is asc order (default):", ok)

	ok = t.Sorted(got, 1)
	fmt.Println("is asc order (1):", ok)

	ok = t.Sorted(got, 0)
	fmt.Println("is asc order (0):", ok)

	ok = t.Sorted(got, nil)
	fmt.Println("is asc order (nil):", ok)

	ok = t.Sorted(got, -1)
	fmt.Println("is desc order:", ok)

	got = []int{-3, -1, 1, 3, -2, 0, 2}
	evenHigher := func(a, b int) bool {
		if (a%2 == 0) != (b%2 == 0) {
			return a%2 != 0
		}
		return a < b
	}
	ok = t.Sorted(got, evenHigher)
	fmt.Println("is even higher order:", ok)

	// Output:
	// is asc order (default): true
	// is asc order (1): true
	// is asc order (0): true
	// is asc order (nil): true
	// is desc order: false
	// is even higher order: true

```{{% /expand%}}
{{%expand "Fields_path example" %}}```go
	t := td.NewT(&testing.T{})

	type Person struct {
		Name string
		Age  int
	}

	alice := Person{Name: "Alice", Age: 20}
	bob := Person{Name: "Bob", Age: 19}
	brian := Person{Name: "Brian", Age: 22}
	marcel := Person{Name: "Marcel", Age: 25}
	stephen := Person{Name: "Stephen", Age: 19}

	got := []Person{alice, bob, brian, marcel, stephen}
	ok := t.Sorted(got, "Name")
	fmt.Println("is sorted by name asc:", ok)

	got = []Person{stephen, marcel, brian, bob, alice}
	ok = t.Sorted(got, "-Name")
	fmt.Println("is sorted by name desc:", ok)

	got = []Person{marcel, brian, alice, bob, stephen}
	ok = t.Sorted(got, []string{"-Age", "Name"})
	fmt.Println("is sorted by age desc, then by name asc 1:", ok)

	got = []Person{marcel, brian, alice, stephen, bob}
	ok = t.Sorted(got, []string{"-Age", "Name"})
	fmt.Println("is sorted by age desc, then by name asc 2:", ok)

	type A struct{ props map[string]int }
	p12 := A{props: map[string]int{"priority": 12}}
	p23 := A{props: map[string]int{"priority": 23}}
	p34 := A{props: map[string]int{"priority": 34}}
	got2 := []A{p34, p23, p12}
	ok = t.Sorted(got2, `-props[priority]`)
	fmt.Println("is sorted by priority desc:", ok)

	ok = t.Sorted(got2, `props.priority`)
	fmt.Println("is sorted by priority asc:", ok)

	// Output:
	// is sorted by name asc: true
	// is sorted by name desc: true
	// is sorted by age desc, then by name asc 1: true
	// is sorted by age desc, then by name asc 2: false
	// is sorted by priority desc: true
	// is sorted by priority asc: false

```{{% /expand%}}
