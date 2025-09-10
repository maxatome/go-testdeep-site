---
title: "List"
weight: 10
---

```go
func List(expectedValues ...any) TestDeep
```

List operator compares the contents of an array or a slice (or a
pointer on array/slice) with taking care of the order of items.

[`Array`]({{% ref "Array" %}}) and [`Slice`]({{% ref "Slice" %}}) need to specify the type of array/slice being
compared then to index all expected items. List does not. It acts
as comparing a literal array/slice, but without having to specify
the type and allowing to easily use TestDeep operators:

```go
td.Cmp(t, []int{1, 9, 5}, td.List(1, 9, 5))                              // succeeds
td.Cmp(t, []int{1, 9, 5}, td.List(td.Gt(0), td.Between(8, 9), td.Lt(5))) // succeeds
td.Cmp(t, []int{1, 9, 5}, td.List(1, 9))                                 // fails, 5 is extra
td.Cmp(t, []int{1, 9, 5}, td.List(1, 9, 5, 4))                           // fails, 4 is missing

// works with slices/arrays of any type
td.Cmp(t, personSlice, td.List(
  Person{Name: "Bob", Age: 32},
  Person{Name: "Alice", Age: 26},
))
```

To flatten a non-`[]any` slice/array, use [`Flatten`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Flatten) function
and so avoid boring and inefficient copies:

```go
expected := []int{1, 2, 1}
td.Cmp(t, []int{1, 1, 2}, td.List(td.Flatten(expected))) // succeeds
// = td.Cmp(t, []int{1, 1, 2}, td.List(1, 2, 1))

// Compare only Name field of a slice of Person structs
td.Cmp(t, personSlice, td.List(td.Flatten([]string{"Bob", "Alice"}, "Smuggle:Name")))
```

[`TypeBehind`]({{% ref "operators#typebehind-method" %}}) method can return a non-`nil` [`reflect.Type`](https://pkg.go.dev/reflect#Type) if all items
known non-interface types are equal, or if only interface types
are found (mostly issued from Isa()) and they are equal.

> See also [`Bag`]({{% ref "Bag" %}}), [`Set`]({{% ref "Set" %}}) and [`Sort`]({{% ref "Sort" %}}).


> See also [<i class='fas fa-book'></i> List godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#List).

### Example

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := []int{1, 33, 8, 2}

	// Matches as all items are present
	ok := td.Cmp(t, got, td.List(1, td.Between(32, 34), td.Gt(7), 2))
	fmt.Println("checks all items match, in this order:", ok)

	// Does not match as got does not use the same order as expected
	ok = td.Cmp(t, got, td.List(1, td.Gt(7), 2, td.Between(32, 34)))
	fmt.Println("checks all items match, in wrong order:", ok)

	// When expected is already a non-[]any slice, it cannot be
	// flattened directly using expected... without copying it to a new
	// []any slice, then use td.Flatten!
	expected := []int{1, 33, 8}
	ok = td.Cmp(t, got, td.List(td.Flatten(expected), td.Lte(2)))
	fmt.Println("checks all expected items are present + last one ≤ 2:", ok)

	// Output:
	// checks all items match, in this order: true
	// checks all items match, in wrong order: false
	// checks all expected items are present + last one ≤ 2: true

```{{% /expand%}}
## CmpList shortcut

```go
func CmpList(t TestingT, got any, expectedValues []any, args ...any) bool
```

CmpList is a shortcut for:

```go
td.Cmp(t, got, td.List(expectedValues...), args...)
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


> See also [<i class='fas fa-book'></i> CmpList godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpList).

### Example

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := []int{1, 33, 8, 2}

	// Matches as all items are present
	ok := td.CmpList(t, got, []any{1, td.Between(32, 34), td.Gt(7), 2})
	fmt.Println("checks all items match, in this order:", ok)

	// Does not match as got does not use the same order as expected
	ok = td.CmpList(t, got, []any{1, td.Gt(7), 2, td.Between(32, 34)})
	fmt.Println("checks all items match, in wrong order:", ok)

	// When expected is already a non-[]any slice, it cannot be
	// flattened directly using expected... without copying it to a new
	// []any slice, then use td.Flatten!
	expected := []int{1, 33, 8}
	ok = td.CmpList(t, got, []any{td.Flatten(expected), td.Lte(2)})
	fmt.Println("checks all expected items are present + last one ≤ 2:", ok)

	// Output:
	// checks all items match, in this order: true
	// checks all items match, in wrong order: false
	// checks all expected items are present + last one ≤ 2: true

```{{% /expand%}}
## T.List shortcut

```go
func (t *T) List(got any, expectedValues []any, args ...any) bool
```

List is a shortcut for:

```go
t.Cmp(got, td.List(expectedValues...), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.List godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.List).

### Example

{{%expand "Base example" %}}```go
	t := td.NewT(&testing.T{})

	got := []int{1, 33, 8, 2}

	// Matches as all items are present
	ok := t.List(got, []any{1, td.Between(32, 34), td.Gt(7), 2})
	fmt.Println("checks all items match, in this order:", ok)

	// Does not match as got does not use the same order as expected
	ok = t.List(got, []any{1, td.Gt(7), 2, td.Between(32, 34)})
	fmt.Println("checks all items match, in wrong order:", ok)

	// When expected is already a non-[]any slice, it cannot be
	// flattened directly using expected... without copying it to a new
	// []any slice, then use td.Flatten!
	expected := []int{1, 33, 8}
	ok = t.List(got, []any{td.Flatten(expected), td.Lte(2)})
	fmt.Println("checks all expected items are present + last one ≤ 2:", ok)

	// Output:
	// checks all items match, in this order: true
	// checks all items match, in wrong order: false
	// checks all expected items are present + last one ≤ 2: true

```{{% /expand%}}
