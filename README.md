# errors [![Travis-CI](https://travis-ci.org/Hatch1fy/errors.svg)](https://travis-ci.org/Hatch1fy/errors) [![GoDoc](https://godoc.org/github.com/Hatch1fy/errors?status.svg)](http://godoc.org/github.com/Hatch1fy/errors) [![Report card](https://goreportcard.com/badge/github.com/Hatch1fy/errors)](https://goreportcard.com/report/github.com/Hatch1fy/errors) [![Sourcegraph](https://sourcegraph.com/github.com/Hatch1fy/errors/-/badge.svg)](https://sourcegraph.com/github.com/Hatch1fy/errors?badge)

Package errors provides simple error handling primitives. This package is a fork of the original [github.com/pkg/errors](https://github.com/pkg/errors) for backwards-compatibility with previous error handling package at Hatchify.

## Getting started

```
go get github.com/Hatch1fy/errors
```

The traditional error handling idiom in Go is roughly akin to:

```go
if err != nil {
    return err
}
```

which applied recursively up the call stack results in error reports without context or debugging information. The errors package allows programmers to add context to the failure path in their code in a way that does not destroy the original value of the error.

### Adding context to an error

The `errors.Wrap` function returns a new error that adds context to the original error. For example:

```go
_, err := ioutil.ReadAll(r)
if err != nil {
    return errors.Wrap(err, "read failed")
}
```

### Retrieving the cause of an error

Using `errors.Wrap` constructs a stack of errors, adding context to the preceding error. Depending on the nature of the error it may be necessary to reverse the operation of errors.Wrap to retrieve the original error for inspection. Any error value which implements this interface can be inspected by `errors.Cause`.

```go
type causer interface {
        Cause() error
}
```

`errors.Cause` will recursively retrieve the topmost error which does not implement `causer`, which is assumed to be the original cause. For example:

```go
switch err := errors.Cause(err).(type) {
case *MyError:
    // handle specifically
default:
    // unknown error
}
```

### Error List

This forked version of the package allows to create multiple wrapped errors, an example:

```go
var errs ErrorList

if err := doSomething(); err != nil {
    errs.Push(err)
}
if err := doMore(); err != nil {
    errs.Push(err)
}
// hurr durr
return errs.Err()
```

Here `ErrorList` implements the `ErrorListIterator` interface, that allows external tools to analyze the list of errors and their individual stacks, while the stack of the resulting error is pointing to the place of `errs.Err()` call.

[Read the package documentation for more information](https://godoc.org/github.com/Hatch1fy/errors).

## License

BSD-2-Clause, see [LICENSE](/LICENSE)
