---
title: "SuperSliceOf"
weight: 10
---

```go
func SuperSliceOf(model any, expectedEntries ArrayEntries) TestDeep
```

SuperSliceOf operator compares the contents of an array, a pointer
on an array, a slice or a pointer on a slice against the non-zero
values of *model* (if `any`) and the values of *expectedEntries*. So
entries with zero value of *model* are always ignored. If a zero
value check is needed, this zero value has to be set in
*expectedEntries*. An entry cannot be present in both *model* and
*expectedEntries*, except if it is a zero-value in *model*. At the
end, only entries present in *expectedEntries* and non-zero ones
present in *model* are checked. To check all entries of an array
see [`Array`]({{< ref "Array" >}}) operator. To check all entries of a slice see [`Slice`]({{< ref "Slice" >}})
operator.

*model* must be the same type as compared data.

*expectedEntries* can be `nil`, if no zero entries are expected and
no [TestDeep operators]({{< ref "operators" >}}) are involved.

Works with slices:

```go
got := []int{12, 14, 17}
td.Cmp(t, got, td.SuperSliceOf([]int{12}, nil))                                // succeeds
td.Cmp(t, got, td.SuperSliceOf([]int{12}, td.ArrayEntries{2: 17}))             // succeeds
td.Cmp(t, &got, td.SuperSliceOf(&[]int{0, 14}, td.ArrayEntries{2: td.Gt(16)})) // succeeds
```

and arrays:

```go
got := [5]int{12, 14, 17, 26, 56}
td.Cmp(t, got, td.SuperSliceOf([5]int{12}, nil))                                // succeeds
td.Cmp(t, got, td.SuperSliceOf([5]int{12}, td.ArrayEntries{2: 17}))             // succeeds
td.Cmp(t, &got, td.SuperSliceOf(&[5]int{0, 14}, td.ArrayEntries{2: td.Gt(16)})) // succeeds
```

> See also [`Array`]({{< ref "Array" >}}) and [`Slice`]({{< ref "Slice" >}}).


> See also [<i class='fas fa-book'></i> SuperSliceOf godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#SuperSliceOf).

### Examples

{{%expand "Array example" %}}```go
	t := &testing.T{}

	got := [4]int{42, 58, 26, 666}

	ok := td.Cmp(t, got,
		td.SuperSliceOf([4]int{1: 58}, td.ArrayEntries{3: td.Gt(660)}),
		"checks array %v", got)
	fmt.Println("Only check items #1 & #3:", ok)

	ok = td.Cmp(t, got,
		td.SuperSliceOf([4]int{}, td.ArrayEntries{0: 42, 3: td.Between(660, 670)}),
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3:", ok)

	ok = td.Cmp(t, &got,
		td.SuperSliceOf(&[4]int{}, td.ArrayEntries{0: 42, 3: td.Between(660, 670)}),
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3 of an array pointer:", ok)

	ok = td.Cmp(t, &got,
		td.SuperSliceOf((*[4]int)(nil), td.ArrayEntries{0: 42, 3: td.Between(660, 670)}),
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3 of an array pointer, using nil model:", ok)

	// Output:
	// Only check items #1 & #3: true
	// Only check items #0 & #3: true
	// Only check items #0 & #3 of an array pointer: true
	// Only check items #0 & #3 of an array pointer, using nil model: true

```{{% /expand%}}
{{%expand "TypedArray example" %}}```go
	t := &testing.T{}

	type MyArray [4]int

	got := MyArray{42, 58, 26, 666}

	ok := td.Cmp(t, got,
		td.SuperSliceOf(MyArray{1: 58}, td.ArrayEntries{3: td.Gt(660)}),
		"checks typed array %v", got)
	fmt.Println("Only check items #1 & #3:", ok)

	ok = td.Cmp(t, got,
		td.SuperSliceOf(MyArray{}, td.ArrayEntries{0: 42, 3: td.Between(660, 670)}),
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3:", ok)

	ok = td.Cmp(t, &got,
		td.SuperSliceOf(&MyArray{}, td.ArrayEntries{0: 42, 3: td.Between(660, 670)}),
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3 of an array pointer:", ok)

	ok = td.Cmp(t, &got,
		td.SuperSliceOf((*MyArray)(nil), td.ArrayEntries{0: 42, 3: td.Between(660, 670)}),
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3 of an array pointer, using nil model:", ok)

	// Output:
	// Only check items #1 & #3: true
	// Only check items #0 & #3: true
	// Only check items #0 & #3 of an array pointer: true
	// Only check items #0 & #3 of an array pointer, using nil model: true

```{{% /expand%}}
{{%expand "Slice example" %}}```go
	t := &testing.T{}

	got := []int{42, 58, 26, 666}

	ok := td.Cmp(t, got,
		td.SuperSliceOf([]int{1: 58}, td.ArrayEntries{3: td.Gt(660)}),
		"checks array %v", got)
	fmt.Println("Only check items #1 & #3:", ok)

	ok = td.Cmp(t, got,
		td.SuperSliceOf([]int{}, td.ArrayEntries{0: 42, 3: td.Between(660, 670)}),
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3:", ok)

	ok = td.Cmp(t, &got,
		td.SuperSliceOf(&[]int{}, td.ArrayEntries{0: 42, 3: td.Between(660, 670)}),
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3 of a slice pointer:", ok)

	ok = td.Cmp(t, &got,
		td.SuperSliceOf((*[]int)(nil), td.ArrayEntries{0: 42, 3: td.Between(660, 670)}),
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3 of a slice pointer, using nil model:", ok)

	// Output:
	// Only check items #1 & #3: true
	// Only check items #0 & #3: true
	// Only check items #0 & #3 of a slice pointer: true
	// Only check items #0 & #3 of a slice pointer, using nil model: true

```{{% /expand%}}
{{%expand "TypedSlice example" %}}```go
	t := &testing.T{}

	type MySlice []int

	got := MySlice{42, 58, 26, 666}

	ok := td.Cmp(t, got,
		td.SuperSliceOf(MySlice{1: 58}, td.ArrayEntries{3: td.Gt(660)}),
		"checks typed array %v", got)
	fmt.Println("Only check items #1 & #3:", ok)

	ok = td.Cmp(t, got,
		td.SuperSliceOf(MySlice{}, td.ArrayEntries{0: 42, 3: td.Between(660, 670)}),
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3:", ok)

	ok = td.Cmp(t, &got,
		td.SuperSliceOf(&MySlice{}, td.ArrayEntries{0: 42, 3: td.Between(660, 670)}),
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3 of a slice pointer:", ok)

	ok = td.Cmp(t, &got,
		td.SuperSliceOf((*MySlice)(nil), td.ArrayEntries{0: 42, 3: td.Between(660, 670)}),
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3 of a slice pointer, using nil model:", ok)

	// Output:
	// Only check items #1 & #3: true
	// Only check items #0 & #3: true
	// Only check items #0 & #3 of a slice pointer: true
	// Only check items #0 & #3 of a slice pointer, using nil model: true

```{{% /expand%}}
## CmpSuperSliceOf shortcut

```go
func CmpSuperSliceOf(t TestingT, got, model any, expectedEntries ArrayEntries, args ...any) bool
```

CmpSuperSliceOf is a shortcut for:

```go
td.Cmp(t, got, td.SuperSliceOf(model, expectedEntries), args...)
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


> See also [<i class='fas fa-book'></i> CmpSuperSliceOf godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpSuperSliceOf).

### Examples

{{%expand "Array example" %}}```go
	t := &testing.T{}

	got := [4]int{42, 58, 26, 666}

	ok := td.CmpSuperSliceOf(t, got, [4]int{1: 58}, td.ArrayEntries{3: td.Gt(660)},
		"checks array %v", got)
	fmt.Println("Only check items #1 & #3:", ok)

	ok = td.CmpSuperSliceOf(t, got, [4]int{}, td.ArrayEntries{0: 42, 3: td.Between(660, 670)},
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3:", ok)

	ok = td.CmpSuperSliceOf(t, &got, &[4]int{}, td.ArrayEntries{0: 42, 3: td.Between(660, 670)},
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3 of an array pointer:", ok)

	ok = td.CmpSuperSliceOf(t, &got, (*[4]int)(nil), td.ArrayEntries{0: 42, 3: td.Between(660, 670)},
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3 of an array pointer, using nil model:", ok)

	// Output:
	// Only check items #1 & #3: true
	// Only check items #0 & #3: true
	// Only check items #0 & #3 of an array pointer: true
	// Only check items #0 & #3 of an array pointer, using nil model: true

```{{% /expand%}}
{{%expand "TypedArray example" %}}```go
	t := &testing.T{}

	type MyArray [4]int

	got := MyArray{42, 58, 26, 666}

	ok := td.CmpSuperSliceOf(t, got, MyArray{1: 58}, td.ArrayEntries{3: td.Gt(660)},
		"checks typed array %v", got)
	fmt.Println("Only check items #1 & #3:", ok)

	ok = td.CmpSuperSliceOf(t, got, MyArray{}, td.ArrayEntries{0: 42, 3: td.Between(660, 670)},
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3:", ok)

	ok = td.CmpSuperSliceOf(t, &got, &MyArray{}, td.ArrayEntries{0: 42, 3: td.Between(660, 670)},
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3 of an array pointer:", ok)

	ok = td.CmpSuperSliceOf(t, &got, (*MyArray)(nil), td.ArrayEntries{0: 42, 3: td.Between(660, 670)},
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3 of an array pointer, using nil model:", ok)

	// Output:
	// Only check items #1 & #3: true
	// Only check items #0 & #3: true
	// Only check items #0 & #3 of an array pointer: true
	// Only check items #0 & #3 of an array pointer, using nil model: true

```{{% /expand%}}
{{%expand "Slice example" %}}```go
	t := &testing.T{}

	got := []int{42, 58, 26, 666}

	ok := td.CmpSuperSliceOf(t, got, []int{1: 58}, td.ArrayEntries{3: td.Gt(660)},
		"checks array %v", got)
	fmt.Println("Only check items #1 & #3:", ok)

	ok = td.CmpSuperSliceOf(t, got, []int{}, td.ArrayEntries{0: 42, 3: td.Between(660, 670)},
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3:", ok)

	ok = td.CmpSuperSliceOf(t, &got, &[]int{}, td.ArrayEntries{0: 42, 3: td.Between(660, 670)},
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3 of a slice pointer:", ok)

	ok = td.CmpSuperSliceOf(t, &got, (*[]int)(nil), td.ArrayEntries{0: 42, 3: td.Between(660, 670)},
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3 of a slice pointer, using nil model:", ok)

	// Output:
	// Only check items #1 & #3: true
	// Only check items #0 & #3: true
	// Only check items #0 & #3 of a slice pointer: true
	// Only check items #0 & #3 of a slice pointer, using nil model: true

```{{% /expand%}}
{{%expand "TypedSlice example" %}}```go
	t := &testing.T{}

	type MySlice []int

	got := MySlice{42, 58, 26, 666}

	ok := td.CmpSuperSliceOf(t, got, MySlice{1: 58}, td.ArrayEntries{3: td.Gt(660)},
		"checks typed array %v", got)
	fmt.Println("Only check items #1 & #3:", ok)

	ok = td.CmpSuperSliceOf(t, got, MySlice{}, td.ArrayEntries{0: 42, 3: td.Between(660, 670)},
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3:", ok)

	ok = td.CmpSuperSliceOf(t, &got, &MySlice{}, td.ArrayEntries{0: 42, 3: td.Between(660, 670)},
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3 of a slice pointer:", ok)

	ok = td.CmpSuperSliceOf(t, &got, (*MySlice)(nil), td.ArrayEntries{0: 42, 3: td.Between(660, 670)},
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3 of a slice pointer, using nil model:", ok)

	// Output:
	// Only check items #1 & #3: true
	// Only check items #0 & #3: true
	// Only check items #0 & #3 of a slice pointer: true
	// Only check items #0 & #3 of a slice pointer, using nil model: true

```{{% /expand%}}
## T.SuperSliceOf shortcut

```go
func (t *T) SuperSliceOf(got, model any, expectedEntries ArrayEntries, args ...any) bool
```

SuperSliceOf is a shortcut for:

```go
t.Cmp(got, td.SuperSliceOf(model, expectedEntries), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.SuperSliceOf godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.SuperSliceOf).

### Examples

{{%expand "Array example" %}}```go
	t := td.NewT(&testing.T{})

	got := [4]int{42, 58, 26, 666}

	ok := t.SuperSliceOf(got, [4]int{1: 58}, td.ArrayEntries{3: td.Gt(660)},
		"checks array %v", got)
	fmt.Println("Only check items #1 & #3:", ok)

	ok = t.SuperSliceOf(got, [4]int{}, td.ArrayEntries{0: 42, 3: td.Between(660, 670)},
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3:", ok)

	ok = t.SuperSliceOf(&got, &[4]int{}, td.ArrayEntries{0: 42, 3: td.Between(660, 670)},
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3 of an array pointer:", ok)

	ok = t.SuperSliceOf(&got, (*[4]int)(nil), td.ArrayEntries{0: 42, 3: td.Between(660, 670)},
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3 of an array pointer, using nil model:", ok)

	// Output:
	// Only check items #1 & #3: true
	// Only check items #0 & #3: true
	// Only check items #0 & #3 of an array pointer: true
	// Only check items #0 & #3 of an array pointer, using nil model: true

```{{% /expand%}}
{{%expand "TypedArray example" %}}```go
	t := td.NewT(&testing.T{})

	type MyArray [4]int

	got := MyArray{42, 58, 26, 666}

	ok := t.SuperSliceOf(got, MyArray{1: 58}, td.ArrayEntries{3: td.Gt(660)},
		"checks typed array %v", got)
	fmt.Println("Only check items #1 & #3:", ok)

	ok = t.SuperSliceOf(got, MyArray{}, td.ArrayEntries{0: 42, 3: td.Between(660, 670)},
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3:", ok)

	ok = t.SuperSliceOf(&got, &MyArray{}, td.ArrayEntries{0: 42, 3: td.Between(660, 670)},
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3 of an array pointer:", ok)

	ok = t.SuperSliceOf(&got, (*MyArray)(nil), td.ArrayEntries{0: 42, 3: td.Between(660, 670)},
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3 of an array pointer, using nil model:", ok)

	// Output:
	// Only check items #1 & #3: true
	// Only check items #0 & #3: true
	// Only check items #0 & #3 of an array pointer: true
	// Only check items #0 & #3 of an array pointer, using nil model: true

```{{% /expand%}}
{{%expand "Slice example" %}}```go
	t := td.NewT(&testing.T{})

	got := []int{42, 58, 26, 666}

	ok := t.SuperSliceOf(got, []int{1: 58}, td.ArrayEntries{3: td.Gt(660)},
		"checks array %v", got)
	fmt.Println("Only check items #1 & #3:", ok)

	ok = t.SuperSliceOf(got, []int{}, td.ArrayEntries{0: 42, 3: td.Between(660, 670)},
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3:", ok)

	ok = t.SuperSliceOf(&got, &[]int{}, td.ArrayEntries{0: 42, 3: td.Between(660, 670)},
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3 of a slice pointer:", ok)

	ok = t.SuperSliceOf(&got, (*[]int)(nil), td.ArrayEntries{0: 42, 3: td.Between(660, 670)},
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3 of a slice pointer, using nil model:", ok)

	// Output:
	// Only check items #1 & #3: true
	// Only check items #0 & #3: true
	// Only check items #0 & #3 of a slice pointer: true
	// Only check items #0 & #3 of a slice pointer, using nil model: true

```{{% /expand%}}
{{%expand "TypedSlice example" %}}```go
	t := td.NewT(&testing.T{})

	type MySlice []int

	got := MySlice{42, 58, 26, 666}

	ok := t.SuperSliceOf(got, MySlice{1: 58}, td.ArrayEntries{3: td.Gt(660)},
		"checks typed array %v", got)
	fmt.Println("Only check items #1 & #3:", ok)

	ok = t.SuperSliceOf(got, MySlice{}, td.ArrayEntries{0: 42, 3: td.Between(660, 670)},
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3:", ok)

	ok = t.SuperSliceOf(&got, &MySlice{}, td.ArrayEntries{0: 42, 3: td.Between(660, 670)},
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3 of a slice pointer:", ok)

	ok = t.SuperSliceOf(&got, (*MySlice)(nil), td.ArrayEntries{0: 42, 3: td.Between(660, 670)},
		"checks array %v", got)
	fmt.Println("Only check items #0 & #3 of a slice pointer, using nil model:", ok)

	// Output:
	// Only check items #1 & #3: true
	// Only check items #0 & #3: true
	// Only check items #0 & #3 of a slice pointer: true
	// Only check items #0 & #3 of a slice pointer, using nil model: true

```{{% /expand%}}
