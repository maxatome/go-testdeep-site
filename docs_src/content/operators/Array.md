---
title: "Array"
weight: 10
---

```go
func Array(model any, expectedEntries ArrayEntries) TestDeep
```

Array operator compares the contents of an array or a pointer on an
array against the values of *model* and the values of
*expectedEntries*. Entries with zero values of *model* are ignored
if the same entry is present in *expectedEntries*, otherwise they
are taken into account. An entry cannot be present in both *model*
and *expectedEntries*, except if it is a zero-value in *model*. At
the end, all entries are checked. To check only some entries of an
array, see [`SuperSliceOf`]({{% ref "SuperSliceOf" %}}) operator.

*model* must be the same type as compared data.

*expectedEntries* can be `nil`, if no zero entries are expected and
no [TestDeep operators]({{% ref "operators" %}}) are involved.

```go
got := [3]int{12, 14, 17}
td.Cmp(t, got, td.Array([3]int{0, 14}, td.ArrayEntries{0: 12, 2: 17})) // succeeds
td.Cmp(t, &got,
  td.Array(&[3]int{0, 14}, td.ArrayEntries{0: td.Gt(10), 2: td.Gt(15)})) // succeeds
```

[`TypeBehind`]({{% ref "operators#typebehind-method" %}}) method returns the [`reflect.Type`](https://pkg.go.dev/reflect#Type) of *model*.

> See also [`Slice`]({{% ref "Slice" %}}) and [`SuperSliceOf`]({{% ref "SuperSliceOf" %}}).


> See also [<i class='fas fa-book'></i> Array godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Array).

### Examples

{{%expand "Array example" %}}```go
	t := &testing.T{}

	got := [3]int{42, 58, 26}

	ok := td.Cmp(t, got,
		td.Array([3]int{42}, td.ArrayEntries{1: 58, 2: td.Ignore()}),
		"checks array %v", got)
	fmt.Println("Simple array:", ok)

	ok = td.Cmp(t, &got,
		td.Array(&[3]int{42}, td.ArrayEntries{1: 58, 2: td.Ignore()}),
		"checks array %v", got)
	fmt.Println("Array pointer:", ok)

	ok = td.Cmp(t, &got,
		td.Array((*[3]int)(nil), td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()}),
		"checks array %v", got)
	fmt.Println("Array pointer, nil model:", ok)

	// Output:
	// Simple array: true
	// Array pointer: true
	// Array pointer, nil model: true

```{{% /expand%}}
{{%expand "TypedArray example" %}}```go
	t := &testing.T{}

	type MyArray [3]int

	got := MyArray{42, 58, 26}

	ok := td.Cmp(t, got,
		td.Array(MyArray{42}, td.ArrayEntries{1: 58, 2: td.Ignore()}),
		"checks typed array %v", got)
	fmt.Println("Typed array:", ok)

	ok = td.Cmp(t, &got,
		td.Array(&MyArray{42}, td.ArrayEntries{1: 58, 2: td.Ignore()}),
		"checks pointer on typed array %v", got)
	fmt.Println("Pointer on a typed array:", ok)

	ok = td.Cmp(t, &got,
		td.Array(&MyArray{}, td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()}),
		"checks pointer on typed array %v", got)
	fmt.Println("Pointer on a typed array, empty model:", ok)

	ok = td.Cmp(t, &got,
		td.Array((*MyArray)(nil), td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()}),
		"checks pointer on typed array %v", got)
	fmt.Println("Pointer on a typed array, nil model:", ok)

	// Output:
	// Typed array: true
	// Pointer on a typed array: true
	// Pointer on a typed array, empty model: true
	// Pointer on a typed array, nil model: true

```{{% /expand%}}
## CmpArray shortcut

```go
func CmpArray(t TestingT, got, model any, expectedEntries ArrayEntries, args ...any) bool
```

CmpArray is a shortcut for:

```go
td.Cmp(t, got, td.Array(model, expectedEntries), args...)
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


> See also [<i class='fas fa-book'></i> CmpArray godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpArray).

### Examples

{{%expand "Array example" %}}```go
	t := &testing.T{}

	got := [3]int{42, 58, 26}

	ok := td.CmpArray(t, got, [3]int{42}, td.ArrayEntries{1: 58, 2: td.Ignore()},
		"checks array %v", got)
	fmt.Println("Simple array:", ok)

	ok = td.CmpArray(t, &got, &[3]int{42}, td.ArrayEntries{1: 58, 2: td.Ignore()},
		"checks array %v", got)
	fmt.Println("Array pointer:", ok)

	ok = td.CmpArray(t, &got, (*[3]int)(nil), td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()},
		"checks array %v", got)
	fmt.Println("Array pointer, nil model:", ok)

	// Output:
	// Simple array: true
	// Array pointer: true
	// Array pointer, nil model: true

```{{% /expand%}}
{{%expand "TypedArray example" %}}```go
	t := &testing.T{}

	type MyArray [3]int

	got := MyArray{42, 58, 26}

	ok := td.CmpArray(t, got, MyArray{42}, td.ArrayEntries{1: 58, 2: td.Ignore()},
		"checks typed array %v", got)
	fmt.Println("Typed array:", ok)

	ok = td.CmpArray(t, &got, &MyArray{42}, td.ArrayEntries{1: 58, 2: td.Ignore()},
		"checks pointer on typed array %v", got)
	fmt.Println("Pointer on a typed array:", ok)

	ok = td.CmpArray(t, &got, &MyArray{}, td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()},
		"checks pointer on typed array %v", got)
	fmt.Println("Pointer on a typed array, empty model:", ok)

	ok = td.CmpArray(t, &got, (*MyArray)(nil), td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()},
		"checks pointer on typed array %v", got)
	fmt.Println("Pointer on a typed array, nil model:", ok)

	// Output:
	// Typed array: true
	// Pointer on a typed array: true
	// Pointer on a typed array, empty model: true
	// Pointer on a typed array, nil model: true

```{{% /expand%}}
## T.Array shortcut

```go
func (t *T) Array(got, model any, expectedEntries ArrayEntries, args ...any) bool
```

Array is a shortcut for:

```go
t.Cmp(got, td.Array(model, expectedEntries), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.Array godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.Array).

### Examples

{{%expand "Array example" %}}```go
	t := td.NewT(&testing.T{})

	got := [3]int{42, 58, 26}

	ok := t.Array(got, [3]int{42}, td.ArrayEntries{1: 58, 2: td.Ignore()},
		"checks array %v", got)
	fmt.Println("Simple array:", ok)

	ok = t.Array(&got, &[3]int{42}, td.ArrayEntries{1: 58, 2: td.Ignore()},
		"checks array %v", got)
	fmt.Println("Array pointer:", ok)

	ok = t.Array(&got, (*[3]int)(nil), td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()},
		"checks array %v", got)
	fmt.Println("Array pointer, nil model:", ok)

	// Output:
	// Simple array: true
	// Array pointer: true
	// Array pointer, nil model: true

```{{% /expand%}}
{{%expand "TypedArray example" %}}```go
	t := td.NewT(&testing.T{})

	type MyArray [3]int

	got := MyArray{42, 58, 26}

	ok := t.Array(got, MyArray{42}, td.ArrayEntries{1: 58, 2: td.Ignore()},
		"checks typed array %v", got)
	fmt.Println("Typed array:", ok)

	ok = t.Array(&got, &MyArray{42}, td.ArrayEntries{1: 58, 2: td.Ignore()},
		"checks pointer on typed array %v", got)
	fmt.Println("Pointer on a typed array:", ok)

	ok = t.Array(&got, &MyArray{}, td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()},
		"checks pointer on typed array %v", got)
	fmt.Println("Pointer on a typed array, empty model:", ok)

	ok = t.Array(&got, (*MyArray)(nil), td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()},
		"checks pointer on typed array %v", got)
	fmt.Println("Pointer on a typed array, nil model:", ok)

	// Output:
	// Typed array: true
	// Pointer on a typed array: true
	// Pointer on a typed array, empty model: true
	// Pointer on a typed array, nil model: true

```{{% /expand%}}
