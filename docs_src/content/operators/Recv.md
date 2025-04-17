---
title: "Recv"
weight: 10
---

```go
func Recv(expectedValue any, timeout ...time.Duration) TestDeep
```

Recv is a [smuggler operator]({{% ref "operators#smuggler-operators" %}}). It reads from a channel or a pointer
to a channel and compares the read value to *expectedValue*.

*expectedValue* can be `any` value including a [TestDeep operator]({{% ref "operators" %}}). It
can also be [`RecvNothing`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#RecvNothing) to test nothing can be read from the
channel or [`RecvClosed`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#RecvClosed) to check the channel is closed.

If *timeout* is passed it should be only one item. It means: try to
read the channel during this duration to get a value before giving
up. If *timeout* is missing or ≤ 0, it defaults to 0 meaning Recv
does not wait for a value but gives up instantly if no value is
available on the channel.

```go
ch := make(chan int, 6)
td.Cmp(t, ch, td.Recv(td.RecvNothing)) // succeeds
td.Cmp(t, ch, td.Recv(42))             // fails, nothing to receive
// recv(DATA): values differ
//      got: nothing received on channel
// expected: 42

ch <- 42
td.Cmp(t, ch, td.Recv(td.RecvNothing)) // fails, 42 received instead
// recv(DATA): values differ
//      got: 42
// expected: nothing received on channel

td.Cmp(t, ch, td.Recv(42)) // fails, nothing to receive anymore
// recv(DATA): values differ
//      got: nothing received on channel
// expected: 42

ch <- 666
td.Cmp(t, ch, td.Recv(td.Between(600, 700))) // succeeds

close(ch)
td.Cmp(t, ch, td.Recv(td.RecvNothing)) // fails as channel is closed
// recv(DATA): values differ
//      got: channel is closed
// expected: nothing received on channel

td.Cmp(t, ch, td.Recv(td.RecvClosed)) // succeeds
```

Note that for convenience Recv accepts pointer on channel:

```go
ch := make(chan int, 6)
ch <- 42
td.Cmp(t, &ch, td.Recv(42)) // succeeds
```

Each time Recv is called, it tries to consume one item from the
channel, immediately or, if given, before *timeout* duration. To
consume several items in a same [`Cmp`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Cmp) call, one can use [`All`]({{% ref "All" %}})
operator as in:

```go
ch := make(chan int, 6)
ch <- 1
ch <- 2
ch <- 3
close(ch)
td.Cmp(t, ch, td.All( // succeeds
  td.Recv(1),
  td.Recv(2),
  td.Recv(3),
  td.Recv(td.RecvClosed),
))
```

To check nothing can be received during 100ms on channel ch (if
something is received before, including a close, it fails):

```go
td.Cmp(t, ch, td.Recv(td.RecvNothing, 100*time.Millisecond))
```

note that in case of success, the above [`Cmp`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Cmp) call always lasts 100ms.

To check 42 can be received from channel ch during the next 100ms
(if nothing is received during these 100ms or something different
from 42, including a close, it fails):

```go
td.Cmp(t, ch, td.Recv(42, 100*time.Millisecond))
```

note that in case of success, the above [`Cmp`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Cmp) call lasts less than 100ms.

A `nil` channel is not handled specifically, so it “is never ready
for communication” as specification says:

```go
var ch chan int
td.Cmp(t, ch, td.Recv(td.RecvNothing)) // always succeeds
td.Cmp(t, ch, td.Recv(42))             // or any other value, always fails
td.Cmp(t, ch, td.Recv(td.RecvClosed))  // always fails
```

so to check if a channel is not `nil` before reading from it, one can
either do:

```go
td.Cmp(t, ch, td.All(
  td.NotNil(),
  td.Recv(42),
))
// or
if td.Cmp(t, ch, td.NotNil()) {
  td.Cmp(t, ch, td.Recv(42))
}
```

[`TypeBehind`]({{% ref "operators#typebehind-method" %}}) method returns the [`reflect.Type`](https://pkg.go.dev/reflect#Type) of *expectedValue*,
except if *expectedValue* is a [TestDeep operator]({{% ref "operators" %}}). In this case, it
delegates [`TypeBehind()`]({{% ref "operators#typebehind-method" %}}) to the operator.

> See also [`Cap`]({{% ref "Cap" %}}) and [`Len`]({{% ref "Len" %}}).


> See also [<i class='fas fa-book'></i> Recv godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Recv).

### Examples

{{%expand "Basic example" %}}```go
	t := &testing.T{}

	got := make(chan int, 3)

	ok := td.Cmp(t, got, td.Recv(td.RecvNothing))
	fmt.Println("nothing to receive:", ok)

	got <- 1
	got <- 2
	got <- 3
	close(got)

	ok = td.Cmp(t, got, td.Recv(1))
	fmt.Println("1st receive is 1:", ok)

	ok = td.Cmp(t, got, td.All(
		td.Recv(2),
		td.Recv(td.Between(3, 4)),
		td.Recv(td.RecvClosed),
	))
	fmt.Println("next receives are 2, 3 then closed:", ok)

	ok = td.Cmp(t, got, td.Recv(td.RecvNothing))
	fmt.Println("nothing to receive:", ok)

	// Output:
	// nothing to receive: true
	// 1st receive is 1: true
	// next receives are 2, 3 then closed: true
	// nothing to receive: false

```{{% /expand%}}
{{%expand "ChannelPointer example" %}}```go
	t := &testing.T{}

	got := make(chan int, 3)

	ok := td.Cmp(t, got, td.Recv(td.RecvNothing))
	fmt.Println("nothing to receive:", ok)

	got <- 1
	got <- 2
	got <- 3
	close(got)

	ok = td.Cmp(t, &got, td.Recv(1))
	fmt.Println("1st receive is 1:", ok)

	ok = td.Cmp(t, &got, td.All(
		td.Recv(2),
		td.Recv(td.Between(3, 4)),
		td.Recv(td.RecvClosed),
	))
	fmt.Println("next receives are 2, 3 then closed:", ok)

	ok = td.Cmp(t, got, td.Recv(td.RecvNothing))
	fmt.Println("nothing to receive:", ok)

	// Output:
	// nothing to receive: true
	// 1st receive is 1: true
	// next receives are 2, 3 then closed: true
	// nothing to receive: false

```{{% /expand%}}
{{%expand "WithTimeout example" %}}```go
	t := &testing.T{}

	got := make(chan int, 1)
	tick := make(chan struct{})

	go func() {
		// ①
		<-tick
		time.Sleep(100 * time.Millisecond)
		got <- 0

		// ②
		<-tick
		time.Sleep(100 * time.Millisecond)
		got <- 1

		// ③
		<-tick
		time.Sleep(100 * time.Millisecond)
		close(got)
	}()

	td.Cmp(t, got, td.Recv(td.RecvNothing))

	// ①
	tick <- struct{}{}
	ok := td.Cmp(t, got, td.Recv(td.RecvNothing))
	fmt.Println("① RecvNothing:", ok)
	ok = td.Cmp(t, got, td.Recv(0, 150*time.Millisecond))
	fmt.Println("① receive 0 w/150ms timeout:", ok)
	ok = td.Cmp(t, got, td.Recv(td.RecvNothing))
	fmt.Println("① RecvNothing:", ok)

	// ②
	tick <- struct{}{}
	ok = td.Cmp(t, got, td.Recv(td.RecvNothing))
	fmt.Println("② RecvNothing:", ok)
	ok = td.Cmp(t, got, td.Recv(1, 150*time.Millisecond))
	fmt.Println("② receive 1 w/150ms timeout:", ok)
	ok = td.Cmp(t, got, td.Recv(td.RecvNothing))
	fmt.Println("② RecvNothing:", ok)

	// ③
	tick <- struct{}{}
	ok = td.Cmp(t, got, td.Recv(td.RecvNothing))
	fmt.Println("③ RecvNothing:", ok)
	ok = td.Cmp(t, got, td.Recv(td.RecvClosed, 150*time.Millisecond))
	fmt.Println("③ check closed w/150ms timeout:", ok)

	// Output:
	// ① RecvNothing: true
	// ① receive 0 w/150ms timeout: true
	// ① RecvNothing: true
	// ② RecvNothing: true
	// ② receive 1 w/150ms timeout: true
	// ② RecvNothing: true
	// ③ RecvNothing: true
	// ③ check closed w/150ms timeout: true

```{{% /expand%}}
{{%expand "NilChannel example" %}}```go
	t := &testing.T{}

	var ch chan int

	ok := td.Cmp(t, ch, td.Recv(td.RecvNothing))
	fmt.Println("nothing to receive from nil channel:", ok)

	ok = td.Cmp(t, ch, td.Recv(42))
	fmt.Println("something to receive from nil channel:", ok)

	ok = td.Cmp(t, ch, td.Recv(td.RecvClosed))
	fmt.Println("is a nil channel closed:", ok)

	// Output:
	// nothing to receive from nil channel: true
	// something to receive from nil channel: false
	// is a nil channel closed: false

```{{% /expand%}}
## CmpRecv shortcut

```go
func CmpRecv(t TestingT, got, expectedValue any, timeout time.Duration, args ...any) bool
```

CmpRecv is a shortcut for:

```go
td.Cmp(t, got, td.Recv(expectedValue, timeout), args...)
```

See above for details.

[`Recv`]({{% ref "Recv" %}}) optional parameter *timeout* is here mandatory.
0 value should be passed to mimic its absence in
original [`Recv`]({{% ref "Recv" %}}) call.

Returns true if the test is OK, false if it fails.

If *t* is a [`*T`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T) then its Config field is inherited.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpRecv godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpRecv).

### Examples

{{%expand "Basic example" %}}```go
	t := &testing.T{}

	got := make(chan int, 3)

	ok := td.CmpRecv(t, got, td.RecvNothing, 0)
	fmt.Println("nothing to receive:", ok)

	got <- 1
	got <- 2
	got <- 3
	close(got)

	ok = td.CmpRecv(t, got, 1, 0)
	fmt.Println("1st receive is 1:", ok)

	ok = td.Cmp(t, got, td.All(
		td.Recv(2),
		td.Recv(td.Between(3, 4)),
		td.Recv(td.RecvClosed),
	))
	fmt.Println("next receives are 2, 3 then closed:", ok)

	ok = td.CmpRecv(t, got, td.RecvNothing, 0)
	fmt.Println("nothing to receive:", ok)

	// Output:
	// nothing to receive: true
	// 1st receive is 1: true
	// next receives are 2, 3 then closed: true
	// nothing to receive: false

```{{% /expand%}}
{{%expand "ChannelPointer example" %}}```go
	t := &testing.T{}

	got := make(chan int, 3)

	ok := td.CmpRecv(t, got, td.RecvNothing, 0)
	fmt.Println("nothing to receive:", ok)

	got <- 1
	got <- 2
	got <- 3
	close(got)

	ok = td.CmpRecv(t, &got, 1, 0)
	fmt.Println("1st receive is 1:", ok)

	ok = td.Cmp(t, &got, td.All(
		td.Recv(2),
		td.Recv(td.Between(3, 4)),
		td.Recv(td.RecvClosed),
	))
	fmt.Println("next receives are 2, 3 then closed:", ok)

	ok = td.CmpRecv(t, got, td.RecvNothing, 0)
	fmt.Println("nothing to receive:", ok)

	// Output:
	// nothing to receive: true
	// 1st receive is 1: true
	// next receives are 2, 3 then closed: true
	// nothing to receive: false

```{{% /expand%}}
{{%expand "WithTimeout example" %}}```go
	t := &testing.T{}

	got := make(chan int, 1)
	tick := make(chan struct{})

	go func() {
		// ①
		<-tick
		time.Sleep(100 * time.Millisecond)
		got <- 0

		// ②
		<-tick
		time.Sleep(100 * time.Millisecond)
		got <- 1

		// ③
		<-tick
		time.Sleep(100 * time.Millisecond)
		close(got)
	}()

	td.CmpRecv(t, got, td.RecvNothing, 0)

	// ①
	tick <- struct{}{}
	ok := td.CmpRecv(t, got, td.RecvNothing, 0)
	fmt.Println("① RecvNothing:", ok)
	ok = td.CmpRecv(t, got, 0, 150*time.Millisecond)
	fmt.Println("① receive 0 w/150ms timeout:", ok)
	ok = td.CmpRecv(t, got, td.RecvNothing, 0)
	fmt.Println("① RecvNothing:", ok)

	// ②
	tick <- struct{}{}
	ok = td.CmpRecv(t, got, td.RecvNothing, 0)
	fmt.Println("② RecvNothing:", ok)
	ok = td.CmpRecv(t, got, 1, 150*time.Millisecond)
	fmt.Println("② receive 1 w/150ms timeout:", ok)
	ok = td.CmpRecv(t, got, td.RecvNothing, 0)
	fmt.Println("② RecvNothing:", ok)

	// ③
	tick <- struct{}{}
	ok = td.CmpRecv(t, got, td.RecvNothing, 0)
	fmt.Println("③ RecvNothing:", ok)
	ok = td.CmpRecv(t, got, td.RecvClosed, 150*time.Millisecond)
	fmt.Println("③ check closed w/150ms timeout:", ok)

	// Output:
	// ① RecvNothing: true
	// ① receive 0 w/150ms timeout: true
	// ① RecvNothing: true
	// ② RecvNothing: true
	// ② receive 1 w/150ms timeout: true
	// ② RecvNothing: true
	// ③ RecvNothing: true
	// ③ check closed w/150ms timeout: true

```{{% /expand%}}
{{%expand "NilChannel example" %}}```go
	t := &testing.T{}

	var ch chan int

	ok := td.CmpRecv(t, ch, td.RecvNothing, 0)
	fmt.Println("nothing to receive from nil channel:", ok)

	ok = td.CmpRecv(t, ch, 42, 0)
	fmt.Println("something to receive from nil channel:", ok)

	ok = td.CmpRecv(t, ch, td.RecvClosed, 0)
	fmt.Println("is a nil channel closed:", ok)

	// Output:
	// nothing to receive from nil channel: true
	// something to receive from nil channel: false
	// is a nil channel closed: false

```{{% /expand%}}
## T.Recv shortcut

```go
func (t *T) Recv(got, expectedValue any, timeout time.Duration, args ...any) bool
```

Recv is a shortcut for:

```go
t.Cmp(got, td.Recv(expectedValue, timeout), args...)
```

See above for details.

[`Recv`]({{% ref "Recv" %}}) optional parameter *timeout* is here mandatory.
0 value should be passed to mimic its absence in
original [`Recv`]({{% ref "Recv" %}}) call.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.Recv godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.Recv).

### Examples

{{%expand "Basic example" %}}```go
	t := td.NewT(&testing.T{})

	got := make(chan int, 3)

	ok := t.Recv(got, td.RecvNothing, 0)
	fmt.Println("nothing to receive:", ok)

	got <- 1
	got <- 2
	got <- 3
	close(got)

	ok = t.Recv(got, 1, 0)
	fmt.Println("1st receive is 1:", ok)

	ok = t.Cmp(got, td.All(
		td.Recv(2),
		td.Recv(td.Between(3, 4)),
		td.Recv(td.RecvClosed),
	))
	fmt.Println("next receives are 2, 3 then closed:", ok)

	ok = t.Recv(got, td.RecvNothing, 0)
	fmt.Println("nothing to receive:", ok)

	// Output:
	// nothing to receive: true
	// 1st receive is 1: true
	// next receives are 2, 3 then closed: true
	// nothing to receive: false

```{{% /expand%}}
{{%expand "ChannelPointer example" %}}```go
	t := td.NewT(&testing.T{})

	got := make(chan int, 3)

	ok := t.Recv(got, td.RecvNothing, 0)
	fmt.Println("nothing to receive:", ok)

	got <- 1
	got <- 2
	got <- 3
	close(got)

	ok = t.Recv(&got, 1, 0)
	fmt.Println("1st receive is 1:", ok)

	ok = t.Cmp(&got, td.All(
		td.Recv(2),
		td.Recv(td.Between(3, 4)),
		td.Recv(td.RecvClosed),
	))
	fmt.Println("next receives are 2, 3 then closed:", ok)

	ok = t.Recv(got, td.RecvNothing, 0)
	fmt.Println("nothing to receive:", ok)

	// Output:
	// nothing to receive: true
	// 1st receive is 1: true
	// next receives are 2, 3 then closed: true
	// nothing to receive: false

```{{% /expand%}}
{{%expand "WithTimeout example" %}}```go
	t := td.NewT(&testing.T{})

	got := make(chan int, 1)
	tick := make(chan struct{})

	go func() {
		// ①
		<-tick
		time.Sleep(100 * time.Millisecond)
		got <- 0

		// ②
		<-tick
		time.Sleep(100 * time.Millisecond)
		got <- 1

		// ③
		<-tick
		time.Sleep(100 * time.Millisecond)
		close(got)
	}()

	t.Recv(got, td.RecvNothing, 0)

	// ①
	tick <- struct{}{}
	ok := t.Recv(got, td.RecvNothing, 0)
	fmt.Println("① RecvNothing:", ok)
	ok = t.Recv(got, 0, 150*time.Millisecond)
	fmt.Println("① receive 0 w/150ms timeout:", ok)
	ok = t.Recv(got, td.RecvNothing, 0)
	fmt.Println("① RecvNothing:", ok)

	// ②
	tick <- struct{}{}
	ok = t.Recv(got, td.RecvNothing, 0)
	fmt.Println("② RecvNothing:", ok)
	ok = t.Recv(got, 1, 150*time.Millisecond)
	fmt.Println("② receive 1 w/150ms timeout:", ok)
	ok = t.Recv(got, td.RecvNothing, 0)
	fmt.Println("② RecvNothing:", ok)

	// ③
	tick <- struct{}{}
	ok = t.Recv(got, td.RecvNothing, 0)
	fmt.Println("③ RecvNothing:", ok)
	ok = t.Recv(got, td.RecvClosed, 150*time.Millisecond)
	fmt.Println("③ check closed w/150ms timeout:", ok)

	// Output:
	// ① RecvNothing: true
	// ① receive 0 w/150ms timeout: true
	// ① RecvNothing: true
	// ② RecvNothing: true
	// ② receive 1 w/150ms timeout: true
	// ② RecvNothing: true
	// ③ RecvNothing: true
	// ③ check closed w/150ms timeout: true

```{{% /expand%}}
{{%expand "NilChannel example" %}}```go
	t := td.NewT(&testing.T{})

	var ch chan int

	ok := t.Recv(ch, td.RecvNothing, 0)
	fmt.Println("nothing to receive from nil channel:", ok)

	ok = t.Recv(ch, 42, 0)
	fmt.Println("something to receive from nil channel:", ok)

	ok = t.Recv(ch, td.RecvClosed, 0)
	fmt.Println("is a nil channel closed:", ok)

	// Output:
	// nothing to receive from nil channel: true
	// something to receive from nil channel: false
	// is a nil channel closed: false

```{{% /expand%}}
