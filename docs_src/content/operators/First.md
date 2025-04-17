---
title: "First"
weight: 10
---

```go
func First(filter, expectedValue any) TestDeep
```

First is a [smuggler operator]({{% ref "operators#smuggler-operators" %}}). It takes an array, a slice or a
pointer on array/slice. For each item it applies *filter*, a
[TestDeep operator]({{% ref "operators" %}}) or a function returning a `bool`. It takes the
first item for which the *filter* matched and compares it to
*expectedValue*. The *filter* matches when it is a:

- [TestDeep operator]({{% ref "operators" %}}) and it matches for the item;
- function receiving the item and it returns true.


*expectedValue* can of course be a [TestDeep operator]({{% ref "operators" %}}).

```go
got := []int{-3, -2, -1, 0, 1, 2, 3}
td.Cmp(t, got, td.First(td.Gt(0), 1))                                    // succeeds
td.Cmp(t, got, td.First(func(x int) bool { return x%2 == 0 }, -2))       // succeeds
td.Cmp(t, got, td.First(func(x int) bool { return x%2 == 0 }, td.Lt(0))) // succeeds
```

If the input is empty (and/or `nil` for a slice), an "item not found"
[`error`](https://pkg.go.dev/builtin#error) is raised before comparing to *expectedValue*.

```go
var got []int
td.Cmp(t, got, td.First(td.Gt(0), td.Gt(0)))      // fails
td.Cmp(t, []int{}, td.First(td.Gt(0), td.Gt(0)))  // fails
td.Cmp(t, [0]int{}, td.First(td.Gt(0), td.Gt(0))) // fails
```

> See also [`Last`]({{% ref "Last" %}}) and [`Grep`]({{% ref "Grep" %}}).


> See also [<i class='fas fa-book'></i> First godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#First).

### Examples

{{%expand "Classic example" %}}```go
	t := &testing.T{}

	got := []int{-3, -2, -1, 0, 1, 2, 3}

	ok := td.Cmp(t, got, td.First(td.Gt(0), 1))
	fmt.Println("first positive number is 1:", ok)

	isEven := func(x int) bool { return x%2 == 0 }

	ok = td.Cmp(t, got, td.First(isEven, -2))
	fmt.Println("first even number is -2:", ok)

	ok = td.Cmp(t, got, td.First(isEven, td.Lt(0)))
	fmt.Println("first even number is < 0:", ok)

	ok = td.Cmp(t, got, td.First(isEven, td.Code(isEven)))
	fmt.Println("first even number is well even:", ok)

	// Output:
	// first positive number is 1: true
	// first even number is -2: true
	// first even number is < 0: true
	// first even number is well even: true

```{{% /expand%}}
{{%expand "Empty example" %}}```go
	t := &testing.T{}

	ok := td.Cmp(t, ([]int)(nil), td.First(td.Gt(0), td.Gt(0)))
	fmt.Println("first in nil slice:", ok)

	ok = td.Cmp(t, []int{}, td.First(td.Gt(0), td.Gt(0)))
	fmt.Println("first in empty slice:", ok)

	ok = td.Cmp(t, &[]int{}, td.First(td.Gt(0), td.Gt(0)))
	fmt.Println("first in empty pointed slice:", ok)

	ok = td.Cmp(t, [0]int{}, td.First(td.Gt(0), td.Gt(0)))
	fmt.Println("first in empty array:", ok)

	// Output:
	// first in nil slice: false
	// first in empty slice: false
	// first in empty pointed slice: false
	// first in empty array: false

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

	ok := td.Cmp(t, got, td.First(
		td.Smuggle("Age", td.Gt(30)),
		td.Smuggle("Fullname", "Bob Foobar")))
	fmt.Println("first person.Age > 30 → Bob:", ok)

	ok = td.Cmp(t, got, td.First(
		td.JSONPointer("/age", td.Gt(30)),
		td.SuperJSONOf(`{"fullname":"Bob Foobar"}`)))
	fmt.Println("first person.Age > 30 → Bob, using JSON:", ok)

	ok = td.Cmp(t, got, td.First(
		td.JSONPointer("/age", td.Gt(30)),
		td.JSONPointer("/fullname", td.HasPrefix("Bob"))))
	fmt.Println("first person.Age > 30 → Bob, using JSONPointer:", ok)

	// Output:
	// first person.Age > 30 → Bob: true
	// first person.Age > 30 → Bob, using JSON: true
	// first person.Age > 30 → Bob, using JSONPointer: true

```{{% /expand%}}
{{%expand "Json example" %}}```go
	t := &testing.T{}

	got := map[string]any{
		"values": []int{1, 2, 3, 4},
	}
	ok := td.Cmp(t, got, td.JSON(`{"values": First(Gt(2), 3)}`))
	fmt.Println("first number > 2:", ok)

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
  "persons": First(JSONPointer("/name", "Brian"), {"id": 4, "name": "Brian"})
}`))
	fmt.Println(`is "Brian" content OK:`, ok)

	ok = td.Cmp(t, got, td.JSON(`
{
  "persons": First(JSONPointer("/name", "Brian"), JSONPointer("/id", 4))
}`))
	fmt.Println(`ID of "Brian" is 4:`, ok)

	// Output:
	// first number > 2: true
	// is "Brian" content OK: true
	// ID of "Brian" is 4: true

```{{% /expand%}}
## CmpFirst shortcut

```go
func CmpFirst(t TestingT, got, filter , expectedValue any, args ...any) bool
```

CmpFirst is a shortcut for:

```go
td.Cmp(t, got, td.First(filter, expectedValue), args...)
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


> See also [<i class='fas fa-book'></i> CmpFirst godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpFirst).

### Examples

{{%expand "Classic example" %}}```go
	t := &testing.T{}

	got := []int{-3, -2, -1, 0, 1, 2, 3}

	ok := td.CmpFirst(t, got, td.Gt(0), 1)
	fmt.Println("first positive number is 1:", ok)

	isEven := func(x int) bool { return x%2 == 0 }

	ok = td.CmpFirst(t, got, isEven, -2)
	fmt.Println("first even number is -2:", ok)

	ok = td.CmpFirst(t, got, isEven, td.Lt(0))
	fmt.Println("first even number is < 0:", ok)

	ok = td.CmpFirst(t, got, isEven, td.Code(isEven))
	fmt.Println("first even number is well even:", ok)

	// Output:
	// first positive number is 1: true
	// first even number is -2: true
	// first even number is < 0: true
	// first even number is well even: true

```{{% /expand%}}
{{%expand "Empty example" %}}```go
	t := &testing.T{}

	ok := td.CmpFirst(t, ([]int)(nil), td.Gt(0), td.Gt(0))
	fmt.Println("first in nil slice:", ok)

	ok = td.CmpFirst(t, []int{}, td.Gt(0), td.Gt(0))
	fmt.Println("first in empty slice:", ok)

	ok = td.CmpFirst(t, &[]int{}, td.Gt(0), td.Gt(0))
	fmt.Println("first in empty pointed slice:", ok)

	ok = td.CmpFirst(t, [0]int{}, td.Gt(0), td.Gt(0))
	fmt.Println("first in empty array:", ok)

	// Output:
	// first in nil slice: false
	// first in empty slice: false
	// first in empty pointed slice: false
	// first in empty array: false

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

	ok := td.CmpFirst(t, got, td.Smuggle("Age", td.Gt(30)), td.Smuggle("Fullname", "Bob Foobar"))
	fmt.Println("first person.Age > 30 → Bob:", ok)

	ok = td.CmpFirst(t, got, td.JSONPointer("/age", td.Gt(30)), td.SuperJSONOf(`{"fullname":"Bob Foobar"}`))
	fmt.Println("first person.Age > 30 → Bob, using JSON:", ok)

	ok = td.CmpFirst(t, got, td.JSONPointer("/age", td.Gt(30)), td.JSONPointer("/fullname", td.HasPrefix("Bob")))
	fmt.Println("first person.Age > 30 → Bob, using JSONPointer:", ok)

	// Output:
	// first person.Age > 30 → Bob: true
	// first person.Age > 30 → Bob, using JSON: true
	// first person.Age > 30 → Bob, using JSONPointer: true

```{{% /expand%}}
## T.First shortcut

```go
func (t *T) First(got, filter , expectedValue any, args ...any) bool
```

First is a shortcut for:

```go
t.Cmp(got, td.First(filter, expectedValue), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.First godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.First).

### Examples

{{%expand "Classic example" %}}```go
	t := td.NewT(&testing.T{})

	got := []int{-3, -2, -1, 0, 1, 2, 3}

	ok := t.First(got, td.Gt(0), 1)
	fmt.Println("first positive number is 1:", ok)

	isEven := func(x int) bool { return x%2 == 0 }

	ok = t.First(got, isEven, -2)
	fmt.Println("first even number is -2:", ok)

	ok = t.First(got, isEven, td.Lt(0))
	fmt.Println("first even number is < 0:", ok)

	ok = t.First(got, isEven, td.Code(isEven))
	fmt.Println("first even number is well even:", ok)

	// Output:
	// first positive number is 1: true
	// first even number is -2: true
	// first even number is < 0: true
	// first even number is well even: true

```{{% /expand%}}
{{%expand "Empty example" %}}```go
	t := td.NewT(&testing.T{})

	ok := t.First(([]int)(nil), td.Gt(0), td.Gt(0))
	fmt.Println("first in nil slice:", ok)

	ok = t.First([]int{}, td.Gt(0), td.Gt(0))
	fmt.Println("first in empty slice:", ok)

	ok = t.First(&[]int{}, td.Gt(0), td.Gt(0))
	fmt.Println("first in empty pointed slice:", ok)

	ok = t.First([0]int{}, td.Gt(0), td.Gt(0))
	fmt.Println("first in empty array:", ok)

	// Output:
	// first in nil slice: false
	// first in empty slice: false
	// first in empty pointed slice: false
	// first in empty array: false

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

	ok := t.First(got, td.Smuggle("Age", td.Gt(30)), td.Smuggle("Fullname", "Bob Foobar"))
	fmt.Println("first person.Age > 30 → Bob:", ok)

	ok = t.First(got, td.JSONPointer("/age", td.Gt(30)), td.SuperJSONOf(`{"fullname":"Bob Foobar"}`))
	fmt.Println("first person.Age > 30 → Bob, using JSON:", ok)

	ok = t.First(got, td.JSONPointer("/age", td.Gt(30)), td.JSONPointer("/fullname", td.HasPrefix("Bob")))
	fmt.Println("first person.Age > 30 → Bob, using JSONPointer:", ok)

	// Output:
	// first person.Age > 30 → Bob: true
	// first person.Age > 30 → Bob, using JSON: true
	// first person.Age > 30 → Bob, using JSONPointer: true

```{{% /expand%}}
