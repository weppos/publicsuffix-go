# Public Suffix <small>Go</small>

Package <tt>publicsuffix</tt> provides a Go domain name parser based on the [Public Suffix List](http://publicsuffix.org/).

[![Build Status](https://travis-ci.org/weppos/publicsuffix-go.svg?branch=master)](https://travis-ci.org/weppos/publicsuffix-go)
[![GoDoc](https://godoc.org/github.com/weppos/publicsuffix-go/publicsuffix?status.svg)](https://godoc.org/github.com/weppos/publicsuffix-go/publicsuffix)


## :warning: Development Warning

This library is currently under development, therefore the methods and the implementation should be considered a work-in-progress. Changes in the method naming, method signatures, public or internal APIs may happen at any time.

Use this library at your own risk. :boom:


## Getting started

```shell
$ git clone git@github.com:weppos/publicsuffix-go.git
$ cd publicsuffix-go
```

Run the test suite.


## Testing

The following command runs the entire test suite.

```shell
$ go test ./...
```

There are 3 different test suites built into this library:

- Acceptance: the acceptance test suite contains some high level tests to ensure the library behaves as expected
- PSL: the PSL test suite runs the library against the []official Public Suffix test cases](https://github.com/publicsuffix/list/blob/master/tests/tests.txt)
- Unit: the unit test suite stresses the various single components of this package


## Installation

```shell
$ go get github.com/weppos/publicsuffix-go
```


## Usage

This is a simple example that demonstrates how to use the package with the default options and the default Public Suffix list packaged with the library.

```go
package main

import (
    "fmt"

    "github.com/weppos/publicsuffix-go/publicsuffix"
)

func main() {
    // Extract the domain from a string
    // using the default list
    fmt.Println(publicsuffix.Domain("example.com"))             // example.com
    fmt.Println(publicsuffix.Domain("www.example.com"))         // example.com
    fmt.Println(publicsuffix.Domain("example.co.uk"))           // example.co.uk
    fmt.Println(publicsuffix.Domain("www.example.co.uk"))       // example.co.uk

    // Parse the domain from a string
    // using the default list
    fmt.Println(publicsuffix.Parse("example.com"))             // &DomainName{"com", "example", ""}
    fmt.Println(publicsuffix.Parse("www.example.com"))         // &DomainName{"com", "example", "www"}
    fmt.Println(publicsuffix.Parse("example.co.uk"))           // &DomainName{"co.uk", "example", ""}
    fmt.Println(publicsuffix.Parse("www.example.co.uk"))       // &DomainName{"co.uk", "example", "www"}
}
```

#### Ignoring Private Domains

The PSL is composed by two list of suffixes: IANA suffixes, and Private Domains.

Private domains are submitted by private organizations. By default, private domains are not ignored.
Sometimes, you want to ignore these domains and only query against the IANA suffixes. You have two options:

1. Ignore the domains at runtime
2. Create a custom list without the private domains

In the first case, the private domains are ignored at runtime: they will still be included in the lists but the lookup will skip them when found.

```go
publicsuffix.DomainFromListWithOptions(publicsuffix.DefaultList(), "google.blogspot.com", nil)
// google.blogspot.com

publicsuffix.DomainFromListWithOptions(publicsuffix.DefaultList(), "google.blogspot.com", &FindOptions{IgnorePrivate: true})
// blogspot.com

// Note that the DefaultFindOptions includes the private domains by default
publicsuffix.DomainFromListWithOptions(publicsuffix.DefaultList(), "google.blogspot.com", DefaultFindOptions)
// google.blogspot.com
```

This solution is easy, but slower. If you find yourself ignoring the private domains in all cases (or in most cases), you may want to create a custom list without the private domains.

```go
list := NewListFromFile("path/to/list.txt", &ParserOption{PrivateDomains: false})
publicsuffix.DomainFromListWithOptions(list, "google.blogspot.com", nil)
// blogspot.com
```


## License

Copyright (c) 2016 Simone Carletti. This is Free Software distributed under the MIT license.
