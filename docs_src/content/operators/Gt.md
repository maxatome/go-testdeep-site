---
title: "Gt"
weight: 10
---

```go
func Gt(minExpectedValue any) TestDeep
```

Gt operator checks that data is greater than
*minExpectedValue*. *minExpectedValue* can be any numeric, `string`,
[`time.Time`](https://pkg.go.dev/time#Time) (or assignable) value or implements at least one of the
two following methods:

```go
func (a T) Less(b T) bool   // returns true if a < b
func (a T) Compare(b T) int // returns -1 if a < b, 1 if a > b, 0 if a == b
```

*minExpectedValue* must be the same type as the compared value,
except if [`BeLax` config flag](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#ContextConfig.BeLax) is true.

```go
td.Cmp(t, 17, td.Gt(15))
before := time.Now()
td.Cmp(t, time.Now(), td.Gt(before))
```

[`TypeBehind`]({{% ref "operators#typebehind-method" %}}) method returns the [`reflect.Type`](https://pkg.go.dev/reflect#Type) of *minExpectedValue*.


> See also [<i class='fas fa-book'></i> Gt godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Gt).

### Examples

{{%expand "Int example" %}}```go
	t := &testing.T{}

	got := 156

	ok := td.Cmp(t, got, td.Gt(155), "checks %v is > 155", got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Gt(156), "checks %v is > 156", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
{{%expand "String example" %}}```go
	t := &testing.T{}

	got := "abc"

	ok := td.Cmp(t, got, td.Gt("abb"), `checks "%v" is > "abb"`, got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Gt("abc"), `checks "%v" is > "abc"`, got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
## CmpGt shortcut

```go
func CmpGt(t TestingT, got, minExpectedValue any, args ...any) bool
```

CmpGt is a shortcut for:

```go
td.Cmp(t, got, td.Gt(minExpectedValue), args...)
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


> See also [<i class='fas fa-book'></i> CmpGt godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpGt).

### Examples

{{%expand "Int example" %}}```go
	t := &testing.T{}

	got := 156

	ok := td.CmpGt(t, got, 155, "checks %v is > 155", got)
	fmt.Println(ok)

	ok = td.CmpGt(t, got, 156, "checks %v is > 156", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
{{%expand "String example" %}}```go
	t := &testing.T{}

	got := "abc"

	ok := td.CmpGt(t, got, "abb", `checks "%v" is > "abb"`, got)
	fmt.Println(ok)

	ok = td.CmpGt(t, got, "abc", `checks "%v" is > "abc"`, got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
## T.Gt shortcut

```go
func (t *T) Gt(got, minExpectedValue any, args ...any) bool
```

Gt is a shortcut for:

```go
t.Cmp(got, td.Gt(minExpectedValue), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.Gt godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.Gt).

### Examples

{{%expand "Int example" %}}```go
	t := td.NewT(&testing.T{})

	got := 156

	ok := t.Gt(got, 155, "checks %v is > 155", got)
	fmt.Println(ok)

	ok = t.Gt(got, 156, "checks %v is > 156", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
{{%expand "String example" %}}```go
	t := td.NewT(&testing.T{})

	got := "abc"

	ok := t.Gt(got, "abb", `checks "%v" is > "abb"`, got)
	fmt.Println(ok)

	ok = t.Gt(got, "abc", `checks "%v" is > "abc"`, got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
