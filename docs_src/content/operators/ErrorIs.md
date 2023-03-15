---
title: "ErrorIs"
weight: 10
---

```go
func ErrorIs(expected error) TestDeep
```

ErrorIs operator reports whether `any` [`error`](https://pkg.go.dev/builtin#error) in an [`error`](https://pkg.go.dev/builtin#error)'s chain
matches *expected*.

```go
_, err := os.Open("/unknown/file")
td.Cmp(t, err, os.ErrNotExist)             // fails
td.Cmp(t, err, td.ErrorIs(os.ErrNotExist)) // succeeds

err1 := fmt.Errorf("failure1")
err2 := fmt.Errorf("failure2: %w", err1)
err3 := fmt.Errorf("failure3: %w", err2)
err := fmt.Errorf("failure4: %w", err3)
td.Cmp(t, err, td.ErrorIs(err))  // succeeds
td.Cmp(t, err, td.ErrorIs(err1)) // succeeds
td.Cmp(t, err1, td.ErrorIs(err)) // fails
```

Behind the scene it uses [`errors.Is`](https://pkg.go.dev/errors#Is) function.

Note that like [`errors.Is`](https://pkg.go.dev/errors#Is), *expected* can be `nil`: in this case the
comparison succeeds when got is `nil` too.

> See also [`CmpError`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpError) and [`CmpNoError`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpNoError).


> See also [<i class='fas fa-book'></i> ErrorIs godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#ErrorIs).

### Example

{{%expand "Base example" %}}```go
	t := &testing.T{}

	err1 := fmt.Errorf("failure1")
	err2 := fmt.Errorf("failure2: %w", err1)
	err3 := fmt.Errorf("failure3: %w", err2)
	err := fmt.Errorf("failure4: %w", err3)

	ok := td.Cmp(t, err, td.ErrorIs(err))
	fmt.Println("error is itself:", ok)

	ok = td.Cmp(t, err, td.ErrorIs(err1))
	fmt.Println("error is also err1:", ok)

	ok = td.Cmp(t, err1, td.ErrorIs(err))
	fmt.Println("err1 is err:", ok)

	// Output:
	// error is itself: true
	// error is also err1: true
	// err1 is err: false

```{{% /expand%}}
## CmpErrorIs shortcut

```go
func CmpErrorIs(t TestingT, got any, expected error, args ...any) bool
```

CmpErrorIs is a shortcut for:

```go
td.Cmp(t, got, td.ErrorIs(expected), args...)
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


> See also [<i class='fas fa-book'></i> CmpErrorIs godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpErrorIs).

### Example

{{%expand "Base example" %}}```go
	t := &testing.T{}

	err1 := fmt.Errorf("failure1")
	err2 := fmt.Errorf("failure2: %w", err1)
	err3 := fmt.Errorf("failure3: %w", err2)
	err := fmt.Errorf("failure4: %w", err3)

	ok := td.CmpErrorIs(t, err, err)
	fmt.Println("error is itself:", ok)

	ok = td.CmpErrorIs(t, err, err1)
	fmt.Println("error is also err1:", ok)

	ok = td.CmpErrorIs(t, err1, err)
	fmt.Println("err1 is err:", ok)

	// Output:
	// error is itself: true
	// error is also err1: true
	// err1 is err: false

```{{% /expand%}}
## T.CmpErrorIs shortcut

```go
func (t *T) CmpErrorIs(got any, expected error, args ...any) bool
```

CmpErrorIs is a shortcut for:

```go
t.Cmp(got, td.ErrorIs(expected), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.CmpErrorIs godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.CmpErrorIs).

### Example

{{%expand "Base example" %}}```go
	t := td.NewT(&testing.T{})

	err1 := fmt.Errorf("failure1")
	err2 := fmt.Errorf("failure2: %w", err1)
	err3 := fmt.Errorf("failure3: %w", err2)
	err := fmt.Errorf("failure4: %w", err3)

	ok := t.CmpErrorIs(err, err)
	fmt.Println("error is itself:", ok)

	ok = t.CmpErrorIs(err, err1)
	fmt.Println("error is also err1:", ok)

	ok = t.CmpErrorIs(err1, err)
	fmt.Println("err1 is err:", ok)

	// Output:
	// error is itself: true
	// error is also err1: true
	// err1 is err: false

```{{% /expand%}}
