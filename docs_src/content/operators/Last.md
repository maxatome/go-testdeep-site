---
title: "Last"
weight: 10
---

```go
func Last(filter, expectedValue any) TestDeep
```

Last is a [smuggler operator]({{% ref "operators#smuggler-operators" %}}). It takes an array, a slice or a
pointer on array/slice. For each item it applies *filter*, a
[TestDeep operator]({{% ref "operators" %}}) or a function returning a `bool`. It takes the
last item for which the *filter* matched and compares it to
*expectedValue*. The *filter* matches when it is a:

- [TestDeep operator]({{% ref "operators" %}}) and it matches for the item;
- function receiving the item and it returns true.


*expectedValue* can of course be a [TestDeep operator]({{% ref "operators" %}}).

```go
got := []int{-3, -2, -1, 0, 1, 2, 3}
td.Cmp(t, got, td.Last(td.Lt(0), -1))                                   // succeeds
td.Cmp(t, got, td.Last(func(x int) bool { return x%2 == 0 }, 2))        // succeeds
td.Cmp(t, got, td.Last(func(x int) bool { return x%2 == 0 }, td.Gt(0))) // succeeds
```

If the input is empty (and/or `nil` for a slice), an "item not found"
[`error`](https://pkg.go.dev/builtin#error) is raised before comparing to *expectedValue*.

```go
var got []int
td.Cmp(t, got, td.Last(td.Gt(0), td.Gt(0)))      // fails
td.Cmp(t, []int{}, td.Last(td.Gt(0), td.Gt(0)))  // fails
td.Cmp(t, [0]int{}, td.Last(td.Gt(0), td.Gt(0))) // fails
```

> See also [`First`]({{% ref "First" %}}) and [`Grep`]({{% ref "Grep" %}}).


> See also [<i class='fas fa-book'></i> Last godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Last).

### Examples

{{%expand "Classic example" %}}```go
	t := &testing.T{}

	got := []int{-3, -2, -1, 0, 1, 2, 3}

	ok := td.Cmp(t, got, td.Last(td.Lt(0), -1))
	fmt.Println("last negative number is -1:", ok)

	isEven := func(x int) bool { return x%2 == 0 }

	ok = td.Cmp(t, got, td.Last(isEven, 2))
	fmt.Println("last even number is 2:", ok)

	ok = td.Cmp(t, got, td.Last(isEven, td.Gt(0)))
	fmt.Println("last even number is > 0:", ok)

	ok = td.Cmp(t, got, td.Last(isEven, td.Code(isEven)))
	fmt.Println("last even number is well even:", ok)

	// Output:
	// last negative number is -1: true
	// last even number is 2: true
	// last even number is > 0: true
	// last even number is well even: true

```{{% /expand%}}
{{%expand "Empty example" %}}```go
	t := &testing.T{}

	ok := td.Cmp(t, ([]int)(nil), td.Last(td.Gt(0), td.Gt(0)))
	fmt.Println("last in nil slice:", ok)

	ok = td.Cmp(t, []int{}, td.Last(td.Gt(0), td.Gt(0)))
	fmt.Println("last in empty slice:", ok)

	ok = td.Cmp(t, &[]int{}, td.Last(td.Gt(0), td.Gt(0)))
	fmt.Println("last in empty pointed slice:", ok)

	ok = td.Cmp(t, [0]int{}, td.Last(td.Gt(0), td.Gt(0)))
	fmt.Println("last in empty array:", ok)

	// Output:
	// last in nil slice: false
	// last in empty slice: false
	// last in empty pointed slice: false
	// last in empty array: false

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
			Age:      37,
		},
	}

	ok := td.Cmp(t, got, td.Last(
		td.Smuggle("Age", td.Gt(30)),
		td.Smuggle("Fullname", "Alice Bingo")))
	fmt.Println("last person.Age > 30 → Alice:", ok)

	ok = td.Cmp(t, got, td.Last(
		td.JSONPointer("/age", td.Gt(30)),
		td.SuperJSONOf(`{"fullname":"Alice Bingo"}`)))
	fmt.Println("last person.Age > 30 → Alice, using JSON:", ok)

	ok = td.Cmp(t, got, td.Last(
		td.JSONPointer("/age", td.Gt(30)),
		td.JSONPointer("/fullname", td.HasPrefix("Alice"))))
	fmt.Println("first person.Age > 30 → Alice, using JSONPointer:", ok)

	// Output:
	// last person.Age > 30 → Alice: true
	// last person.Age > 30 → Alice, using JSON: true
	// first person.Age > 30 → Alice, using JSONPointer: true

```{{% /expand%}}
{{%expand "Json example" %}}```go
	t := &testing.T{}

	got := map[string]any{
		"values": []int{1, 2, 3, 4},
	}
	ok := td.Cmp(t, got, td.JSON(`{"values": Last(Lt(3), 2)}`))
	fmt.Println("last number < 3:", ok)

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
  "persons": Last(JSONPointer("/name", "Brian"), {"id": 4, "name": "Brian"})
}`))
	fmt.Println(`is "Brian" content OK:`, ok)

	ok = td.Cmp(t, got, td.JSON(`
{
  "persons": Last(JSONPointer("/name", "Brian"), JSONPointer("/id", 4))
}`))
	fmt.Println(`ID of "Brian" is 4:`, ok)

	// Output:
	// last number < 3: true
	// is "Brian" content OK: true
	// ID of "Brian" is 4: true

```{{% /expand%}}
## CmpLast shortcut

```go
func CmpLast(t TestingT, got, filter , expectedValue any, args ...any) bool
```

CmpLast is a shortcut for:

```go
td.Cmp(t, got, td.Last(filter, expectedValue), args...)
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


> See also [<i class='fas fa-book'></i> CmpLast godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpLast).

### Examples

{{%expand "Classic example" %}}```go
	t := &testing.T{}

	got := []int{-3, -2, -1, 0, 1, 2, 3}

	ok := td.CmpLast(t, got, td.Lt(0), -1)
	fmt.Println("last negative number is -1:", ok)

	isEven := func(x int) bool { return x%2 == 0 }

	ok = td.CmpLast(t, got, isEven, 2)
	fmt.Println("last even number is 2:", ok)

	ok = td.CmpLast(t, got, isEven, td.Gt(0))
	fmt.Println("last even number is > 0:", ok)

	ok = td.CmpLast(t, got, isEven, td.Code(isEven))
	fmt.Println("last even number is well even:", ok)

	// Output:
	// last negative number is -1: true
	// last even number is 2: true
	// last even number is > 0: true
	// last even number is well even: true

```{{% /expand%}}
{{%expand "Empty example" %}}```go
	t := &testing.T{}

	ok := td.CmpLast(t, ([]int)(nil), td.Gt(0), td.Gt(0))
	fmt.Println("last in nil slice:", ok)

	ok = td.CmpLast(t, []int{}, td.Gt(0), td.Gt(0))
	fmt.Println("last in empty slice:", ok)

	ok = td.CmpLast(t, &[]int{}, td.Gt(0), td.Gt(0))
	fmt.Println("last in empty pointed slice:", ok)

	ok = td.CmpLast(t, [0]int{}, td.Gt(0), td.Gt(0))
	fmt.Println("last in empty array:", ok)

	// Output:
	// last in nil slice: false
	// last in empty slice: false
	// last in empty pointed slice: false
	// last in empty array: false

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
			Age:      37,
		},
	}

	ok := td.CmpLast(t, got, td.Smuggle("Age", td.Gt(30)), td.Smuggle("Fullname", "Alice Bingo"))
	fmt.Println("last person.Age > 30 → Alice:", ok)

	ok = td.CmpLast(t, got, td.JSONPointer("/age", td.Gt(30)), td.SuperJSONOf(`{"fullname":"Alice Bingo"}`))
	fmt.Println("last person.Age > 30 → Alice, using JSON:", ok)

	ok = td.CmpLast(t, got, td.JSONPointer("/age", td.Gt(30)), td.JSONPointer("/fullname", td.HasPrefix("Alice")))
	fmt.Println("first person.Age > 30 → Alice, using JSONPointer:", ok)

	// Output:
	// last person.Age > 30 → Alice: true
	// last person.Age > 30 → Alice, using JSON: true
	// first person.Age > 30 → Alice, using JSONPointer: true

```{{% /expand%}}
## T.Last shortcut

```go
func (t *T) Last(got, filter , expectedValue any, args ...any) bool
```

Last is a shortcut for:

```go
t.Cmp(got, td.Last(filter, expectedValue), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.Last godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.Last).

### Examples

{{%expand "Classic example" %}}```go
	t := td.NewT(&testing.T{})

	got := []int{-3, -2, -1, 0, 1, 2, 3}

	ok := t.Last(got, td.Lt(0), -1)
	fmt.Println("last negative number is -1:", ok)

	isEven := func(x int) bool { return x%2 == 0 }

	ok = t.Last(got, isEven, 2)
	fmt.Println("last even number is 2:", ok)

	ok = t.Last(got, isEven, td.Gt(0))
	fmt.Println("last even number is > 0:", ok)

	ok = t.Last(got, isEven, td.Code(isEven))
	fmt.Println("last even number is well even:", ok)

	// Output:
	// last negative number is -1: true
	// last even number is 2: true
	// last even number is > 0: true
	// last even number is well even: true

```{{% /expand%}}
{{%expand "Empty example" %}}```go
	t := td.NewT(&testing.T{})

	ok := t.Last(([]int)(nil), td.Gt(0), td.Gt(0))
	fmt.Println("last in nil slice:", ok)

	ok = t.Last([]int{}, td.Gt(0), td.Gt(0))
	fmt.Println("last in empty slice:", ok)

	ok = t.Last(&[]int{}, td.Gt(0), td.Gt(0))
	fmt.Println("last in empty pointed slice:", ok)

	ok = t.Last([0]int{}, td.Gt(0), td.Gt(0))
	fmt.Println("last in empty array:", ok)

	// Output:
	// last in nil slice: false
	// last in empty slice: false
	// last in empty pointed slice: false
	// last in empty array: false

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
			Age:      37,
		},
	}

	ok := t.Last(got, td.Smuggle("Age", td.Gt(30)), td.Smuggle("Fullname", "Alice Bingo"))
	fmt.Println("last person.Age > 30 → Alice:", ok)

	ok = t.Last(got, td.JSONPointer("/age", td.Gt(30)), td.SuperJSONOf(`{"fullname":"Alice Bingo"}`))
	fmt.Println("last person.Age > 30 → Alice, using JSON:", ok)

	ok = t.Last(got, td.JSONPointer("/age", td.Gt(30)), td.JSONPointer("/fullname", td.HasPrefix("Alice")))
	fmt.Println("first person.Age > 30 → Alice, using JSONPointer:", ok)

	// Output:
	// last person.Age > 30 → Alice: true
	// last person.Age > 30 → Alice, using JSON: true
	// first person.Age > 30 → Alice, using JSONPointer: true

```{{% /expand%}}
