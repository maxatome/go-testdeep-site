---
title: "Code"
weight: 10
---

```go
func Code(fn any) TestDeep
```

Code operator allows to check data using a custom function. So
*fn* is a function that must take one parameter whose type must be
the same as the type of the compared value.

*fn* can return a single `bool` kind value, telling that yes or no
the custom test is successful:

```go
td.Cmp(t, gotTime,
  td.Code(func(date time.Time) bool {
    return date.Year() == 2018
  }))
```

or two values (`bool`, `string`) kinds. The `bool` value has the same
meaning as above, and the `string` value is used to describe the
test when it fails:

```go
td.Cmp(t, gotTime,
  td.Code(func(date time.Time) (bool, string) {
    if date.Year() == 2018 {
      return true, ""
    }
    return false, "year must be 2018"
  }))
```

or a single [`error`](https://pkg.go.dev/builtin#error) value. If the returned [`error`](https://pkg.go.dev/builtin#error) is `nil`, the test
succeeded, else the [`error`](https://pkg.go.dev/builtin#error) contains the reason of failure:

```go
td.Cmp(t, gotJsonRawMesg,
  td.Code(func(b json.RawMessage) error {
    var c map[string]int
    err := json.Unmarshal(b, &c)
    if err != nil {
      return err
    }
    if c["test"] != 42 {
      return fmt.Errorf(`key "test" does not match 42`)
    }
    return nil
  }))
```

This operator allows to handle `any` specific comparison not handled
by standard operators.

It is not recommended to call [`Cmp`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Cmp) (or `any` other Cmp*
functions or [`*T`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T) methods) inside the body of *fn*, because of
confusion produced by output in case of failure. When the data
needs to be transformed before being compared again, [`Smuggle`]({{% ref "Smuggle" %}})
operator should be used instead.

But in some cases it can be better to handle yourself the
comparison than to chain [TestDeep operators]({{% ref "operators" %}}). In this case, *fn* can
be a function receiving one or two [`*T`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T) as first parameters and
returning no values.

When *fn* expects one [`*T`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T) parameter, it is directly derived from the
[`testing.TB`](https://pkg.go.dev/testing#TB) instance passed originally to [`Cmp`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Cmp) (or its derivatives)
using [`NewT`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#NewT):

```go
td.Cmp(t, httpRequest, td.Code(func(t *td.T, r *http.Request) {
  token, err := DecodeToken(r.Header.Get("X-Token-1"))
  if t.CmpNoError(err) {
    t.True(token.OK())
  }
}))
```

When *fn* expects two [`*T`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T) parameters, they are directly derived from
the [`testing.TB`](https://pkg.go.dev/testing#TB) instance passed originally to [`Cmp`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Cmp) (or its derivatives)
using [`AssertRequire`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#AssertRequire):

```go
td.Cmp(t, httpRequest, td.Code(func(assert, require *td.T, r *http.Request) {
  token, err := DecodeToken(r.Header.Get("X-Token-1"))
  require.CmpNoError(err)
  assert.True(token.OK())
}))
```

Note that these forms do not work when there is no initial
[`testing.TB`](https://pkg.go.dev/testing#TB) instance, like when using [`EqDeeplyError`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#EqDeeplyError) or
[`EqDeeply`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#EqDeeply) functions, or when the Code operator is called behind
the following operators, as they just check if a match occurs
without raising an [`error`](https://pkg.go.dev/builtin#error): [`Any`]({{% ref "Any" %}}), [`Bag`]({{% ref "Bag" %}}), [`Contains`]({{% ref "Contains" %}}), [`ContainsKey`]({{% ref "ContainsKey" %}}),
[`None`]({{% ref "None" %}}), [`Not`]({{% ref "Not" %}}), [`NotAny`]({{% ref "NotAny" %}}), [`Set`]({{% ref "Set" %}}), [`SubBagOf`]({{% ref "SubBagOf" %}}), [`SubSetOf`]({{% ref "SubSetOf" %}}),
[`SuperBagOf`]({{% ref "SuperBagOf" %}}) and [`SuperSetOf`]({{% ref "SuperSetOf" %}}).

RootName is inherited but not the current path, but it can be
recovered if needed:

```go
got := map[string]int{"foo": 123}
td.NewT(t).
  RootName("PIPO").
  Cmp(got, td.Map(map[string]int{}, td.MapEntries{
    "foo": td.Code(func(t *td.T, n int) {
      t.Cmp(n, 124)                                   // inherit only RootName
      t.RootName(t.Config.OriginalPath()).Cmp(n, 125) // recover current path
      t.RootName("").Cmp(n, 126)                      // undo RootName inheritance
    }),
  }))
```

produces the following errors:

```
--- FAIL: TestCodeCustom (0.00s)
    td_code_test.go:339: Failed test
        PIPO: values differ             ← inherit only RootName
               got: 123
          expected: 124
    td_code_test.go:338: Failed test
        PIPO["foo"]: values differ      ← recover current path
               got: 123
          expected: 125
    td_code_test.go:342: Failed test
        DATA: values differ             ← undo RootName inheritance
               got: 123
          expected: 126
```

[`TypeBehind`]({{% ref "operators#typebehind-method" %}}) method returns the [`reflect.Type`](https://pkg.go.dev/reflect#Type) of last parameter of *fn*.


> See also [<i class='fas fa-book'></i> Code godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Code).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := "12"

	ok := td.Cmp(t, got,
		td.Code(func(num string) bool {
			n, err := strconv.Atoi(num)
			return err == nil && n > 10 && n < 100
		}),
		"checks string `%s` contains a number and this number is in ]10 .. 100[",
		got)
	fmt.Println(ok)

	// Same with failure reason
	ok = td.Cmp(t, got,
		td.Code(func(num string) (bool, string) {
			n, err := strconv.Atoi(num)
			if err != nil {
				return false, "not a number"
			}
			if n > 10 && n < 100 {
				return true, ""
			}
			return false, "not in ]10 .. 100["
		}),
		"checks string `%s` contains a number and this number is in ]10 .. 100[",
		got)
	fmt.Println(ok)

	// Same with failure reason thanks to error
	ok = td.Cmp(t, got,
		td.Code(func(num string) error {
			n, err := strconv.Atoi(num)
			if err != nil {
				return err
			}
			if n > 10 && n < 100 {
				return nil
			}
			return fmt.Errorf("%d not in ]10 .. 100[", n)
		}),
		"checks string `%s` contains a number and this number is in ]10 .. 100[",
		got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true

```{{% /expand%}}
{{%expand "Custom example" %}}```go
	t := &testing.T{}

	got := 123

	ok := td.Cmp(t, got, td.Code(func(t *td.T, num int) {
		t.Cmp(num, 123)
	}))
	fmt.Println("with one *td.T:", ok)

	ok = td.Cmp(t, got, td.Code(func(assert, require *td.T, num int) {
		assert.Cmp(num, 123)
		require.Cmp(num, 123)
	}))
	fmt.Println("with assert & require *td.T:", ok)

	// Output:
	// with one *td.T: true
	// with assert & require *td.T: true

```{{% /expand%}}
## CmpCode shortcut

```go
func CmpCode(t TestingT, got, fn any, args ...any) bool
```

CmpCode is a shortcut for:

```go
td.Cmp(t, got, td.Code(fn), args...)
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


> See also [<i class='fas fa-book'></i> CmpCode godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpCode).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := "12"

	ok := td.CmpCode(t, got, func(num string) bool {
		n, err := strconv.Atoi(num)
		return err == nil && n > 10 && n < 100
	},
		"checks string `%s` contains a number and this number is in ]10 .. 100[",
		got)
	fmt.Println(ok)

	// Same with failure reason
	ok = td.CmpCode(t, got, func(num string) (bool, string) {
		n, err := strconv.Atoi(num)
		if err != nil {
			return false, "not a number"
		}
		if n > 10 && n < 100 {
			return true, ""
		}
		return false, "not in ]10 .. 100["
	},
		"checks string `%s` contains a number and this number is in ]10 .. 100[",
		got)
	fmt.Println(ok)

	// Same with failure reason thanks to error
	ok = td.CmpCode(t, got, func(num string) error {
		n, err := strconv.Atoi(num)
		if err != nil {
			return err
		}
		if n > 10 && n < 100 {
			return nil
		}
		return fmt.Errorf("%d not in ]10 .. 100[", n)
	},
		"checks string `%s` contains a number and this number is in ]10 .. 100[",
		got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true

```{{% /expand%}}
{{%expand "Custom example" %}}```go
	t := &testing.T{}

	got := 123

	ok := td.CmpCode(t, got, func(t *td.T, num int) {
		t.Cmp(num, 123)
	})
	fmt.Println("with one *td.T:", ok)

	ok = td.CmpCode(t, got, func(assert, require *td.T, num int) {
		assert.Cmp(num, 123)
		require.Cmp(num, 123)
	})
	fmt.Println("with assert & require *td.T:", ok)

	// Output:
	// with one *td.T: true
	// with assert & require *td.T: true

```{{% /expand%}}
## T.Code shortcut

```go
func (t *T) Code(got, fn any, args ...any) bool
```

Code is a shortcut for:

```go
t.Cmp(got, td.Code(fn), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.Code godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.Code).

### Examples

{{%expand "Base example" %}}```go
	t := td.NewT(&testing.T{})

	got := "12"

	ok := t.Code(got, func(num string) bool {
		n, err := strconv.Atoi(num)
		return err == nil && n > 10 && n < 100
	},
		"checks string `%s` contains a number and this number is in ]10 .. 100[",
		got)
	fmt.Println(ok)

	// Same with failure reason
	ok = t.Code(got, func(num string) (bool, string) {
		n, err := strconv.Atoi(num)
		if err != nil {
			return false, "not a number"
		}
		if n > 10 && n < 100 {
			return true, ""
		}
		return false, "not in ]10 .. 100["
	},
		"checks string `%s` contains a number and this number is in ]10 .. 100[",
		got)
	fmt.Println(ok)

	// Same with failure reason thanks to error
	ok = t.Code(got, func(num string) error {
		n, err := strconv.Atoi(num)
		if err != nil {
			return err
		}
		if n > 10 && n < 100 {
			return nil
		}
		return fmt.Errorf("%d not in ]10 .. 100[", n)
	},
		"checks string `%s` contains a number and this number is in ]10 .. 100[",
		got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true

```{{% /expand%}}
{{%expand "Custom example" %}}```go
	t := td.NewT(&testing.T{})

	got := 123

	ok := t.Code(got, func(t *td.T, num int) {
		t.Cmp(num, 123)
	})
	fmt.Println("with one *td.T:", ok)

	ok = t.Code(got, func(assert, require *td.T, num int) {
		assert.Cmp(num, 123)
		require.Cmp(num, 123)
	})
	fmt.Println("with assert & require *td.T:", ok)

	// Output:
	// with one *td.T: true
	// with assert & require *td.T: true

```{{% /expand%}}
