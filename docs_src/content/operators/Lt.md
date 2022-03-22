---
title: "Lt"
weight: 10
---

```go
func Lt(maxExpectedValue any) TestDeep
```

[`Lt`]({{< ref "Lt" >}}) operator checks that data is lesser than
*maxExpectedValue*. *maxExpectedValue* can be `any` numeric, `string`,
[`time.Time`](https://pkg.go.dev/time/#Time) (or assignable) value or implements at least one of the
two following methods:
```go
func (a T) Less(b T) bool   // returns true if a < b
func (a T) Compare(b T) int // returns -1 if a < b, 1 if a > b, 0 if a == b
```

*maxExpectedValue* must be the same type as the compared value,
except if BeLax config flag is true.

```go
td.Cmp(t, 17, td.Lt(19))
before := time.Now()
td.Cmp(t, before, td.Lt(time.Now()))
```

[`TypeBehind`]({{< ref "operators#typebehind-method" >}}) method returns the [`reflect.Type`](https://pkg.go.dev/reflect/#Type) of *maxExpectedValue*.


> See also [<i class='fas fa-book'></i> Lt godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Lt).

### Examples

{{%expand "Int example" %}}```go
	t := &testing.T{}

	got := 156

	ok := td.Cmp(t, got, td.Lt(157), "checks %v is < 157", got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Lt(156), "checks %v is < 156", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
{{%expand "String example" %}}```go
	t := &testing.T{}

	got := "abc"

	ok := td.Cmp(t, got, td.Lt("abd"), `checks "%v" is < "abd"`, got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Lt("abc"), `checks "%v" is < "abc"`, got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
## CmpLt shortcut

```go
func CmpLt(t TestingT, got, maxExpectedValue any, args ...any) bool
```

CmpLt is a shortcut for:

```go
td.Cmp(t, got, td.Lt(maxExpectedValue), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

If "t" is a *T then its Config is inherited.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpLt godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpLt).

### Examples

{{%expand "Int example" %}}```go
	t := &testing.T{}

	got := 156

	ok := td.CmpLt(t, got, 157, "checks %v is < 157", got)
	fmt.Println(ok)

	ok = td.CmpLt(t, got, 156, "checks %v is < 156", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
{{%expand "String example" %}}```go
	t := &testing.T{}

	got := "abc"

	ok := td.CmpLt(t, got, "abd", `checks "%v" is < "abd"`, got)
	fmt.Println(ok)

	ok = td.CmpLt(t, got, "abc", `checks "%v" is < "abc"`, got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
## T.Lt shortcut

```go
func (t *T) Lt(got, maxExpectedValue any, args ...any) bool
```

[`Lt`]({{< ref "Lt" >}}) is a shortcut for:

```go
t.Cmp(got, td.Lt(maxExpectedValue), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.Lt godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.Lt).

### Examples

{{%expand "Int example" %}}```go
	t := td.NewT(&testing.T{})

	got := 156

	ok := t.Lt(got, 157, "checks %v is < 157", got)
	fmt.Println(ok)

	ok = t.Lt(got, 156, "checks %v is < 156", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
{{%expand "String example" %}}```go
	t := td.NewT(&testing.T{})

	got := "abc"

	ok := t.Lt(got, "abd", `checks "%v" is < "abd"`, got)
	fmt.Println(ok)

	ok = t.Lt(got, "abc", `checks "%v" is < "abc"`, got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
