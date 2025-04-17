---
title: "Grep"
weight: 10
---

```go
func Grep(filter, expectedValue any) TestDeep
```

Grep is a [smuggler operator]({{% ref "operators#smuggler-operators" %}}). It takes an array, a slice or a
pointer on array/slice. For each item it applies *filter*, a
[TestDeep operator]({{% ref "operators" %}}) or a function returning a `bool`, and produces a
slice consisting of those items for which the *filter* matched and
compares it to *expectedValue*. The *filter* matches when it is a:

- [TestDeep operator]({{% ref "operators" %}}) and it matches for the item;
- function receiving the item and it returns true.


*expectedValue* can be a [TestDeep operator]({{% ref "operators" %}}) or a slice (but never an
array nor a pointer on a slice/array nor `any` other kind).

```go
got := []int{-3, -2, -1, 0, 1, 2, 3}
td.Cmp(t, got, td.Grep(td.Gt(0), []int{1, 2, 3})) // succeeds
td.Cmp(t, got, td.Grep(
  func(x int) bool { return x%2 == 0 },
  []int{-2, 0, 2})) // succeeds
td.Cmp(t, got, td.Grep(
  func(x int) bool { return x%2 == 0 },
  td.Set(0, 2, -2))) // succeeds
```

If Grep receives a `nil` slice or a pointer on a `nil` slice, it always
returns a `nil` slice:

```go
var got []int
td.Cmp(t, got, td.Grep(td.Gt(0), ([]int)(nil))) // succeeds
td.Cmp(t, got, td.Grep(td.Gt(0), td.Nil()))     // succeeds
td.Cmp(t, got, td.Grep(td.Gt(0), []int{}))      // fails
```

> See also [`First`]({{% ref "First" %}}), [`Last`]({{% ref "Last" %}}) and [`Flatten`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Flatten).


> See also [<i class='fas fa-book'></i> Grep godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Grep).

### Examples

{{%expand "Classic example" %}}```go
	t := &testing.T{}

	got := []int{-3, -2, -1, 0, 1, 2, 3}

	ok := td.Cmp(t, got, td.Grep(td.Gt(0), []int{1, 2, 3}))
	fmt.Println("check positive numbers:", ok)

	isEven := func(x int) bool { return x%2 == 0 }

	ok = td.Cmp(t, got, td.Grep(isEven, []int{-2, 0, 2}))
	fmt.Println("even numbers are -2, 0 and 2:", ok)

	ok = td.Cmp(t, got, td.Grep(isEven, td.Set(0, 2, -2)))
	fmt.Println("even numbers are also 0, 2 and -2:", ok)

	ok = td.Cmp(t, got, td.Grep(isEven, td.ArrayEach(td.Code(isEven))))
	fmt.Println("even numbers are each even:", ok)

	// Output:
	// check positive numbers: true
	// even numbers are -2, 0 and 2: true
	// even numbers are also 0, 2 and -2: true
	// even numbers are each even: true

```{{% /expand%}}
{{%expand "Nil example" %}}```go
	t := &testing.T{}

	var got []int
	ok := td.Cmp(t, got, td.Grep(td.Gt(0), ([]int)(nil)))
	fmt.Println("typed []int nil:", ok)

	ok = td.Cmp(t, got, td.Grep(td.Gt(0), ([]string)(nil)))
	fmt.Println("typed []string nil:", ok)

	ok = td.Cmp(t, got, td.Grep(td.Gt(0), td.Nil()))
	fmt.Println("td.Nil:", ok)

	ok = td.Cmp(t, got, td.Grep(td.Gt(0), []int{}))
	fmt.Println("empty non-nil slice:", ok)

	// Output:
	// typed []int nil: true
	// typed []string nil: false
	// td.Nil: true
	// empty non-nil slice: false

```{{% /expand%}}
{{%expand "Struct example" %}}```go
	t := &testing.T{}

	type Person struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
	}

	got := []*Person{
		{
			Fullname: "Bob Foobar",
			Age:      42,
		},
		{
			Fullname: "Alice Bingo",
			Age:      27,
		},
	}

	ok := td.Cmp(t, got, td.Grep(
		td.Smuggle("Age", td.Gt(30)),
		td.All(
			td.Len(1),
			td.ArrayEach(td.Smuggle("Fullname", "Bob Foobar")),
		)))
	fmt.Println("person.Age > 30 → only Bob:", ok)

	ok = td.Cmp(t, got, td.Grep(
		td.JSONPointer("/age", td.Gt(30)),
		td.JSON(`[ SuperMapOf({"fullname":"Bob Foobar"}) ]`)))
	fmt.Println("person.Age > 30 → only Bob, using JSON:", ok)

	// Output:
	// person.Age > 30 → only Bob: true
	// person.Age > 30 → only Bob, using JSON: true

```{{% /expand%}}
{{%expand "Json example" %}}```go
	t := &testing.T{}

	got := map[string]any{
		"values": []int{1, 2, 3, 4},
	}
	ok := td.Cmp(t, got, td.JSON(`{"values": Grep(Gt(2), [3, 4])}`))
	fmt.Println("grep a number > 2:", ok)

	got = map[string]any{
		"persons": []map[string]any{
			{"id": 1, "name": "Joe"},
			{"id": 2, "name": "Bob"},
			{"id": 3, "name": "Alice"},
			{"id": 4, "name": "Brian"},
			{"id": 5, "name": "Britt"},
		},
	}
	ok = td.Cmp(t, got, td.JSON(`
{
  "persons": Grep(JSONPointer("/name", HasPrefix("Br")), [
    {"id": 4, "name": "Brian"},
    {"id": 5, "name": "Britt"},
  ])
}`))
	fmt.Println(`grep "Br" prefix:`, ok)

	// Output:
	// grep a number > 2: true
	// grep "Br" prefix: true

```{{% /expand%}}
## CmpGrep shortcut

```go
func CmpGrep(t TestingT, got, filter , expectedValue any, args ...any) bool
```

CmpGrep is a shortcut for:

```go
td.Cmp(t, got, td.Grep(filter, expectedValue), args...)
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


> See also [<i class='fas fa-book'></i> CmpGrep godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpGrep).

### Examples

{{%expand "Classic example" %}}```go
	t := &testing.T{}

	got := []int{-3, -2, -1, 0, 1, 2, 3}

	ok := td.CmpGrep(t, got, td.Gt(0), []int{1, 2, 3})
	fmt.Println("check positive numbers:", ok)

	isEven := func(x int) bool { return x%2 == 0 }

	ok = td.CmpGrep(t, got, isEven, []int{-2, 0, 2})
	fmt.Println("even numbers are -2, 0 and 2:", ok)

	ok = td.CmpGrep(t, got, isEven, td.Set(0, 2, -2))
	fmt.Println("even numbers are also 0, 2 and -2:", ok)

	ok = td.CmpGrep(t, got, isEven, td.ArrayEach(td.Code(isEven)))
	fmt.Println("even numbers are each even:", ok)

	// Output:
	// check positive numbers: true
	// even numbers are -2, 0 and 2: true
	// even numbers are also 0, 2 and -2: true
	// even numbers are each even: true

```{{% /expand%}}
{{%expand "Nil example" %}}```go
	t := &testing.T{}

	var got []int
	ok := td.CmpGrep(t, got, td.Gt(0), ([]int)(nil))
	fmt.Println("typed []int nil:", ok)

	ok = td.CmpGrep(t, got, td.Gt(0), ([]string)(nil))
	fmt.Println("typed []string nil:", ok)

	ok = td.CmpGrep(t, got, td.Gt(0), td.Nil())
	fmt.Println("td.Nil:", ok)

	ok = td.CmpGrep(t, got, td.Gt(0), []int{})
	fmt.Println("empty non-nil slice:", ok)

	// Output:
	// typed []int nil: true
	// typed []string nil: false
	// td.Nil: true
	// empty non-nil slice: false

```{{% /expand%}}
{{%expand "Struct example" %}}```go
	t := &testing.T{}

	type Person struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
	}

	got := []*Person{
		{
			Fullname: "Bob Foobar",
			Age:      42,
		},
		{
			Fullname: "Alice Bingo",
			Age:      27,
		},
	}

	ok := td.CmpGrep(t, got, td.Smuggle("Age", td.Gt(30)), td.All(
		td.Len(1),
		td.ArrayEach(td.Smuggle("Fullname", "Bob Foobar")),
	))
	fmt.Println("person.Age > 30 → only Bob:", ok)

	ok = td.CmpGrep(t, got, td.JSONPointer("/age", td.Gt(30)), td.JSON(`[ SuperMapOf({"fullname":"Bob Foobar"}) ]`))
	fmt.Println("person.Age > 30 → only Bob, using JSON:", ok)

	// Output:
	// person.Age > 30 → only Bob: true
	// person.Age > 30 → only Bob, using JSON: true

```{{% /expand%}}
## T.Grep shortcut

```go
func (t *T) Grep(got, filter , expectedValue any, args ...any) bool
```

Grep is a shortcut for:

```go
t.Cmp(got, td.Grep(filter, expectedValue), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.Grep godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.Grep).

### Examples

{{%expand "Classic example" %}}```go
	t := td.NewT(&testing.T{})

	got := []int{-3, -2, -1, 0, 1, 2, 3}

	ok := t.Grep(got, td.Gt(0), []int{1, 2, 3})
	fmt.Println("check positive numbers:", ok)

	isEven := func(x int) bool { return x%2 == 0 }

	ok = t.Grep(got, isEven, []int{-2, 0, 2})
	fmt.Println("even numbers are -2, 0 and 2:", ok)

	ok = t.Grep(got, isEven, td.Set(0, 2, -2))
	fmt.Println("even numbers are also 0, 2 and -2:", ok)

	ok = t.Grep(got, isEven, td.ArrayEach(td.Code(isEven)))
	fmt.Println("even numbers are each even:", ok)

	// Output:
	// check positive numbers: true
	// even numbers are -2, 0 and 2: true
	// even numbers are also 0, 2 and -2: true
	// even numbers are each even: true

```{{% /expand%}}
{{%expand "Nil example" %}}```go
	t := td.NewT(&testing.T{})

	var got []int
	ok := t.Grep(got, td.Gt(0), ([]int)(nil))
	fmt.Println("typed []int nil:", ok)

	ok = t.Grep(got, td.Gt(0), ([]string)(nil))
	fmt.Println("typed []string nil:", ok)

	ok = t.Grep(got, td.Gt(0), td.Nil())
	fmt.Println("td.Nil:", ok)

	ok = t.Grep(got, td.Gt(0), []int{})
	fmt.Println("empty non-nil slice:", ok)

	// Output:
	// typed []int nil: true
	// typed []string nil: false
	// td.Nil: true
	// empty non-nil slice: false

```{{% /expand%}}
{{%expand "Struct example" %}}```go
	t := td.NewT(&testing.T{})

	type Person struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
	}

	got := []*Person{
		{
			Fullname: "Bob Foobar",
			Age:      42,
		},
		{
			Fullname: "Alice Bingo",
			Age:      27,
		},
	}

	ok := t.Grep(got, td.Smuggle("Age", td.Gt(30)), td.All(
		td.Len(1),
		td.ArrayEach(td.Smuggle("Fullname", "Bob Foobar")),
	))
	fmt.Println("person.Age > 30 → only Bob:", ok)

	ok = t.Grep(got, td.JSONPointer("/age", td.Gt(30)), td.JSON(`[ SuperMapOf({"fullname":"Bob Foobar"}) ]`))
	fmt.Println("person.Age > 30 → only Bob, using JSON:", ok)

	// Output:
	// person.Age > 30 → only Bob: true
	// person.Age > 30 → only Bob, using JSON: true

```{{% /expand%}}
