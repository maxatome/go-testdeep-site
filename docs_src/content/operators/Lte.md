---
title: "Lte"
weight: 10
---

```go
func Lte(maxExpectedValue any) TestDeep
```

Lte operator checks that data is lesser or equal than
*maxExpectedValue*. *maxExpectedValue* can be any numeric, `string`,
[`time.Time`](https://pkg.go.dev/time#Time) (or assignable) value or implements at least one of the
two following methods:

```go
func (a T) Less(b T) bool   // returns true if a < b
func (a T) Compare(b T) int // returns -1 if a < b, 1 if a > b, 0 if a == b
```

*maxExpectedValue* must be the same type as the compared value,
except if [`BeLax` config flag](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#ContextConfig.BeLax) is true.

```go
td.Cmp(t, 17, td.Lte(17))
before := time.Now()
td.Cmp(t, before, td.Lt(time.Now()))
```

[`TypeBehind`]({{% ref "operators#typebehind-method" %}}) method returns the [`reflect.Type`](https://pkg.go.dev/reflect#Type) of *maxExpectedValue*.


> See also [<i class='fas fa-book'></i> Lte godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Lte).

### Examples

{{%expand "Int example" %}}```go
	t := &testing.T{}

	got := 156

	ok := td.Cmp(t, got, td.Lte(156), "checks %v is ≤ 156", got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Lte(157), "checks %v is ≤ 157", got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Lte(155), "checks %v is ≤ 155", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false

```{{% /expand%}}
{{%expand "String example" %}}```go
	t := &testing.T{}

	got := "abc"

	ok := td.Cmp(t, got, td.Lte("abc"), `checks "%v" is ≤ "abc"`, got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Lte("abd"), `checks "%v" is ≤ "abd"`, got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Lte("abb"), `checks "%v" is ≤ "abb"`, got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false

```{{% /expand%}}
## CmpLte shortcut

```go
func CmpLte(t TestingT, got, maxExpectedValue any, args ...any) bool
```

CmpLte is a shortcut for:

```go
td.Cmp(t, got, td.Lte(maxExpectedValue), args...)
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


> See also [<i class='fas fa-book'></i> CmpLte godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpLte).

### Examples

{{%expand "Int example" %}}```go
	t := &testing.T{}

	got := 156

	ok := td.CmpLte(t, got, 156, "checks %v is ≤ 156", got)
	fmt.Println(ok)

	ok = td.CmpLte(t, got, 157, "checks %v is ≤ 157", got)
	fmt.Println(ok)

	ok = td.CmpLte(t, got, 155, "checks %v is ≤ 155", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false

```{{% /expand%}}
{{%expand "String example" %}}```go
	t := &testing.T{}

	got := "abc"

	ok := td.CmpLte(t, got, "abc", `checks "%v" is ≤ "abc"`, got)
	fmt.Println(ok)

	ok = td.CmpLte(t, got, "abd", `checks "%v" is ≤ "abd"`, got)
	fmt.Println(ok)

	ok = td.CmpLte(t, got, "abb", `checks "%v" is ≤ "abb"`, got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false

```{{% /expand%}}
## T.Lte shortcut

```go
func (t *T) Lte(got, maxExpectedValue any, args ...any) bool
```

Lte is a shortcut for:

```go
t.Cmp(got, td.Lte(maxExpectedValue), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.Lte godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.Lte).

### Examples

{{%expand "Int example" %}}```go
	t := td.NewT(&testing.T{})

	got := 156

	ok := t.Lte(got, 156, "checks %v is ≤ 156", got)
	fmt.Println(ok)

	ok = t.Lte(got, 157, "checks %v is ≤ 157", got)
	fmt.Println(ok)

	ok = t.Lte(got, 155, "checks %v is ≤ 155", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false

```{{% /expand%}}
{{%expand "String example" %}}```go
	t := td.NewT(&testing.T{})

	got := "abc"

	ok := t.Lte(got, "abc", `checks "%v" is ≤ "abc"`, got)
	fmt.Println(ok)

	ok = t.Lte(got, "abd", `checks "%v" is ≤ "abd"`, got)
	fmt.Println(ok)

	ok = t.Lte(got, "abb", `checks "%v" is ≤ "abb"`, got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false

```{{% /expand%}}
