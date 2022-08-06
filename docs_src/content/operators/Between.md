---
title: "Between"
weight: 10
---

```go
func Between(from, to any, bounds ...BoundsKind) TestDeep
```

Between operator checks that data is between *from* and
*to*. *from* and *to* can be any numeric, `string`, [`time.Time`](https://pkg.go.dev/time#Time) (or
assignable) value or implement at least one of the two following
methods:

```go
func (a T) Less(b T) bool   // returns true if a < b
func (a T) Compare(b T) int // returns -1 if a < b, 1 if a > b, 0 if a == b
```

*from* and *to* must be the same type as the compared value, except
if [`BeLax` config flag](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#ContextConfig.BeLax) is true. [`time.Duration`](https://pkg.go.dev/time#Duration) type is accepted as
*to* when *from* is [`time.Time`](https://pkg.go.dev/time#Time) or convertible. *bounds* allows *to*
specify whether *bounds* are included or not:

- [`BoundsInIn`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#BoundsInIn) (default): between *from* and *to* both included
- [`BoundsInOut`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#BoundsInOut): between *from* included and *to* excluded
- [`BoundsOutIn`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#BoundsOutIn): between *from* excluded and *to* included
- [`BoundsOutOut`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#BoundsOutOut): between *from* and *to* both excluded


If *bounds* is missing, it defaults *to* [`BoundsInIn`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#BoundsInIn).

```go
tc.Cmp(t, 17, td.Between(17, 20))               // succeeds, BoundsInIn by default
tc.Cmp(t, 17, td.Between(10, 17, BoundsInOut))  // fails
tc.Cmp(t, 17, td.Between(10, 17, BoundsOutIn))  // succeeds
tc.Cmp(t, 17, td.Between(17, 20, BoundsOutOut)) // fails
tc.Cmp(t,                                       // succeeds
  netip.MustParse("127.0.0.1"),
  td.Between(netip.MustParse("127.0.0.0"), netip.MustParse("127.255.255.255")))
```

[`TypeBehind`]({{< ref "operators#typebehind-method" >}}) method returns the [`reflect.Type`](https://pkg.go.dev/reflect#Type) of *from*.


> See also [<i class='fas fa-book'></i> Between godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Between).

### Examples

{{%expand "Int example" %}}```go
	t := &testing.T{}

	got := 156

	ok := td.Cmp(t, got, td.Between(154, 156),
		"checks %v is in [154 .. 156]", got)
	fmt.Println(ok)

	// BoundsInIn is implicit
	ok = td.Cmp(t, got, td.Between(154, 156, td.BoundsInIn),
		"checks %v is in [154 .. 156]", got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Between(154, 156, td.BoundsInOut),
		"checks %v is in [154 .. 156[", got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Between(154, 156, td.BoundsOutIn),
		"checks %v is in ]154 .. 156]", got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Between(154, 156, td.BoundsOutOut),
		"checks %v is in ]154 .. 156[", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
	// true
	// false

```{{% /expand%}}
{{%expand "String example" %}}```go
	t := &testing.T{}

	got := "abc"

	ok := td.Cmp(t, got, td.Between("aaa", "abc"),
		`checks "%v" is in ["aaa" .. "abc"]`, got)
	fmt.Println(ok)

	// BoundsInIn is implicit
	ok = td.Cmp(t, got, td.Between("aaa", "abc", td.BoundsInIn),
		`checks "%v" is in ["aaa" .. "abc"]`, got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Between("aaa", "abc", td.BoundsInOut),
		`checks "%v" is in ["aaa" .. "abc"[`, got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Between("aaa", "abc", td.BoundsOutIn),
		`checks "%v" is in ]"aaa" .. "abc"]`, got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Between("aaa", "abc", td.BoundsOutOut),
		`checks "%v" is in ]"aaa" .. "abc"[`, got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
	// true
	// false

```{{% /expand%}}
{{%expand "Time example" %}}```go
	t := &testing.T{}

	before := time.Now()
	occurredAt := time.Now()
	after := time.Now()

	ok := td.Cmp(t, occurredAt, td.Between(before, after))
	fmt.Println("It occurred between before and after:", ok)

	type MyTime time.Time
	ok = td.Cmp(t, MyTime(occurredAt), td.Between(MyTime(before), MyTime(after)))
	fmt.Println("Same for convertible MyTime type:", ok)

	ok = td.Cmp(t, MyTime(occurredAt), td.Between(before, after))
	fmt.Println("MyTime vs time.Time:", ok)

	ok = td.Cmp(t, occurredAt, td.Between(before, 10*time.Second))
	fmt.Println("Using a time.Duration as TO:", ok)

	ok = td.Cmp(t, MyTime(occurredAt), td.Between(MyTime(before), 10*time.Second))
	fmt.Println("Using MyTime as FROM and time.Duration as TO:", ok)

	// Output:
	// It occurred between before and after: true
	// Same for convertible MyTime type: true
	// MyTime vs time.Time: false
	// Using a time.Duration as TO: true
	// Using MyTime as FROM and time.Duration as TO: true

```{{% /expand%}}
## CmpBetween shortcut

```go
func CmpBetween(t TestingT, got, from , to any, bounds BoundsKind, args ...any) bool
```

CmpBetween is a shortcut for:

```go
td.Cmp(t, got, td.Between(from, to, bounds), args...)
```

See above for details.

[`Between`]({{< ref "Between" >}}) optional parameter *bounds* is here mandatory.
[`BoundsInIn`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#BoundsInIn) value should be passed *to* mimic its absence in
original [`Between`]({{< ref "Between" >}}) call.

Returns true if the test is OK, false if it fails.

If *t* is a [`*T`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T) then its Config field is inherited.

*args...* are optional and allow *to* name the test. This name is
used in case of failure *to* qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt#Fprintf) is used *to* compose the name, else *args* are passed *to*
[`fmt.Fprint`](https://pkg.go.dev/fmt#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpBetween godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpBetween).

### Examples

{{%expand "Int example" %}}```go
	t := &testing.T{}

	got := 156

	ok := td.CmpBetween(t, got, 154, 156, td.BoundsInIn,
		"checks %v is in [154 .. 156]", got)
	fmt.Println(ok)

	// BoundsInIn is implicit
	ok = td.CmpBetween(t, got, 154, 156, td.BoundsInIn,
		"checks %v is in [154 .. 156]", got)
	fmt.Println(ok)

	ok = td.CmpBetween(t, got, 154, 156, td.BoundsInOut,
		"checks %v is in [154 .. 156[", got)
	fmt.Println(ok)

	ok = td.CmpBetween(t, got, 154, 156, td.BoundsOutIn,
		"checks %v is in ]154 .. 156]", got)
	fmt.Println(ok)

	ok = td.CmpBetween(t, got, 154, 156, td.BoundsOutOut,
		"checks %v is in ]154 .. 156[", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
	// true
	// false

```{{% /expand%}}
{{%expand "String example" %}}```go
	t := &testing.T{}

	got := "abc"

	ok := td.CmpBetween(t, got, "aaa", "abc", td.BoundsInIn,
		`checks "%v" is in ["aaa" .. "abc"]`, got)
	fmt.Println(ok)

	// BoundsInIn is implicit
	ok = td.CmpBetween(t, got, "aaa", "abc", td.BoundsInIn,
		`checks "%v" is in ["aaa" .. "abc"]`, got)
	fmt.Println(ok)

	ok = td.CmpBetween(t, got, "aaa", "abc", td.BoundsInOut,
		`checks "%v" is in ["aaa" .. "abc"[`, got)
	fmt.Println(ok)

	ok = td.CmpBetween(t, got, "aaa", "abc", td.BoundsOutIn,
		`checks "%v" is in ]"aaa" .. "abc"]`, got)
	fmt.Println(ok)

	ok = td.CmpBetween(t, got, "aaa", "abc", td.BoundsOutOut,
		`checks "%v" is in ]"aaa" .. "abc"[`, got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
	// true
	// false

```{{% /expand%}}
{{%expand "Time example" %}}```go
	t := &testing.T{}

	before := time.Now()
	occurredAt := time.Now()
	after := time.Now()

	ok := td.CmpBetween(t, occurredAt, before, after, td.BoundsInIn)
	fmt.Println("It occurred between before and after:", ok)

	type MyTime time.Time
	ok = td.CmpBetween(t, MyTime(occurredAt), MyTime(before), MyTime(after), td.BoundsInIn)
	fmt.Println("Same for convertible MyTime type:", ok)

	ok = td.CmpBetween(t, MyTime(occurredAt), before, after, td.BoundsInIn)
	fmt.Println("MyTime vs time.Time:", ok)

	ok = td.CmpBetween(t, occurredAt, before, 10*time.Second, td.BoundsInIn)
	fmt.Println("Using a time.Duration as TO:", ok)

	ok = td.CmpBetween(t, MyTime(occurredAt), MyTime(before), 10*time.Second, td.BoundsInIn)
	fmt.Println("Using MyTime as FROM and time.Duration as TO:", ok)

	// Output:
	// It occurred between before and after: true
	// Same for convertible MyTime type: true
	// MyTime vs time.Time: false
	// Using a time.Duration as TO: true
	// Using MyTime as FROM and time.Duration as TO: true

```{{% /expand%}}
## T.Between shortcut

```go
func (t *T) Between(got, from , to any, bounds BoundsKind, args ...any) bool
```

Between is a shortcut for:

```go
t.Cmp(got, td.Between(from, to, bounds), args...)
```

See above for details.

[`Between`]({{< ref "Between" >}}) optional parameter *bounds* is here mandatory.
[`BoundsInIn`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#BoundsInIn) value should be passed *to* mimic its absence in
original [`Between`]({{< ref "Between" >}}) call.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow *to* name the test. This name is
used in case of failure *to* qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt#Fprintf) is used *to* compose the name, else *args* are passed *to*
[`fmt.Fprint`](https://pkg.go.dev/fmt#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.Between godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.Between).

### Examples

{{%expand "Int example" %}}```go
	t := td.NewT(&testing.T{})

	got := 156

	ok := t.Between(got, 154, 156, td.BoundsInIn,
		"checks %v is in [154 .. 156]", got)
	fmt.Println(ok)

	// BoundsInIn is implicit
	ok = t.Between(got, 154, 156, td.BoundsInIn,
		"checks %v is in [154 .. 156]", got)
	fmt.Println(ok)

	ok = t.Between(got, 154, 156, td.BoundsInOut,
		"checks %v is in [154 .. 156[", got)
	fmt.Println(ok)

	ok = t.Between(got, 154, 156, td.BoundsOutIn,
		"checks %v is in ]154 .. 156]", got)
	fmt.Println(ok)

	ok = t.Between(got, 154, 156, td.BoundsOutOut,
		"checks %v is in ]154 .. 156[", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
	// true
	// false

```{{% /expand%}}
{{%expand "String example" %}}```go
	t := td.NewT(&testing.T{})

	got := "abc"

	ok := t.Between(got, "aaa", "abc", td.BoundsInIn,
		`checks "%v" is in ["aaa" .. "abc"]`, got)
	fmt.Println(ok)

	// BoundsInIn is implicit
	ok = t.Between(got, "aaa", "abc", td.BoundsInIn,
		`checks "%v" is in ["aaa" .. "abc"]`, got)
	fmt.Println(ok)

	ok = t.Between(got, "aaa", "abc", td.BoundsInOut,
		`checks "%v" is in ["aaa" .. "abc"[`, got)
	fmt.Println(ok)

	ok = t.Between(got, "aaa", "abc", td.BoundsOutIn,
		`checks "%v" is in ]"aaa" .. "abc"]`, got)
	fmt.Println(ok)

	ok = t.Between(got, "aaa", "abc", td.BoundsOutOut,
		`checks "%v" is in ]"aaa" .. "abc"[`, got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
	// true
	// false

```{{% /expand%}}
{{%expand "Time example" %}}```go
	t := td.NewT(&testing.T{})

	before := time.Now()
	occurredAt := time.Now()
	after := time.Now()

	ok := t.Between(occurredAt, before, after, td.BoundsInIn)
	fmt.Println("It occurred between before and after:", ok)

	type MyTime time.Time
	ok = t.Between(MyTime(occurredAt), MyTime(before), MyTime(after), td.BoundsInIn)
	fmt.Println("Same for convertible MyTime type:", ok)

	ok = t.Between(MyTime(occurredAt), before, after, td.BoundsInIn)
	fmt.Println("MyTime vs time.Time:", ok)

	ok = t.Between(occurredAt, before, 10*time.Second, td.BoundsInIn)
	fmt.Println("Using a time.Duration as TO:", ok)

	ok = t.Between(MyTime(occurredAt), MyTime(before), 10*time.Second, td.BoundsInIn)
	fmt.Println("Using MyTime as FROM and time.Duration as TO:", ok)

	// Output:
	// It occurred between before and after: true
	// Same for convertible MyTime type: true
	// MyTime vs time.Time: false
	// Using a time.Duration as TO: true
	// Using MyTime as FROM and time.Duration as TO: true

```{{% /expand%}}
