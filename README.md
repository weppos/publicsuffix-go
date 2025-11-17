# PublicSuffix for Go

The `publicsuffix` package provides a Go domain name parser based on the [Public Suffix List](https://publicsuffix.org/).

[![Tests](https://github.com/weppos/publicsuffix-go/workflows/Tests/badge.svg)](https://github.com/weppos/publicsuffix-go/actions?query=workflow%3ATests)
[![GoDoc](https://godoc.org/github.com/weppos/publicsuffix-go/publicsuffix?status.svg)](https://pkg.go.dev/github.com/weppos/publicsuffix-go/publicsuffix)


## Requirements

`publicsuffix-go` requires **Go >= 1.21**. We do our best not to break older versions of Go if we don't have to, but due to tooling constraints, we don't always test older versions.


## Getting started

Clone the repository [in your workspace](https://golang.org/doc/code.html#Organization) and move into it:

```shell
mkdir -p $GOPATH/src/github.com/weppos && cd $_
git clone git@github.com:weppos/publicsuffix-go.git
cd publicsuffix-go
```

Fetch the dependencies:

```shell
go get ./...
```

Run the test suite.

```shell
go test ./...
```


## Installation

```shell
go get github.com/weppos/publicsuffix-go
```


## Testing

The following command runs the entire test suite.

```shell
go test ./...
```

There are 3 different test suites built into this library:

- Acceptance: the acceptance test suite contains some high level tests to ensure the library behaves as expected
- PSL: the PSL test suite runs the library against the [official Public Suffix test cases](https://github.com/publicsuffix/list/blob/master/tests/tests.txt)
- Unit: the unit test suite stresses the various single components of this package


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

### Private domains

The PSL is composed by two list of suffixes: IANA suffixes, and Private Domains.

Private domains are submitted by private organizations. By default, private domains are not ignored.
Sometimes, you want to ignore these domains and only query against the IANA suffixes. You have two options:

1. Ignore the domains at runtime
2. Create a custom list without the private domains

In the first case, the private domains are ignored at runtime: they will still be included in the lists but the lookup will skip them when found.

```go
publicsuffix.DomainFromListWithOptions(publicsuffix.DefaultList(), "google.blogspot.com", nil)
// google.blogspot.com

publicsuffix.DomainFromListWithOptions(publicsuffix.DefaultList(), "google.blogspot.com", &publicsuffix.FindOptions{IgnorePrivate: true})
// blogspot.com

// Note that the DefaultFindOptions includes the private domains by default
publicsuffix.DomainFromListWithOptions(publicsuffix.DefaultList(), "google.blogspot.com", publicsuffix.DefaultFindOptions)
// google.blogspot.com
```

This solution is easy, but slower. If you find yourself ignoring the private domains in all cases (or in most cases), you may want to create a custom list without the private domains.

```go
list := NewListFromFile("path/to/list.txt", &publicsuffix.ParserOption{PrivateDomains: false})
publicsuffix.DomainFromListWithOptions(list, "google.blogspot.com", nil)
// blogspot.com
```

## IDN domains, A-labels and U-labels

[A-label and U-label](https://tools.ietf.org/html/rfc5890#section-2.3.2.1) are two different ways to represent IDN domain names. These two encodings are also known as ASCII (A-label) or Pynucode vs Unicode (U-label). Conversions between U-labels and A-labels are performed according to the ["Punycode" specification](https://tools.ietf.org/html/rfc3492), adding or removing the ACE prefix as needed.

IDNA-aware applications generally use the A-label form for storing and manipulating data, whereas the U-labels can appear in presentation and user interface forms.

Although the PSL list has been traditionally U-label encoded, this library follows the common industry standards and stores the rules in their A-label form. Therefore, unless explicitly mentioned, any method call, comparison or internal representation is expected to be ASCII-compatible encoded (ACE).

Passing Unicode names to the library may either result in error or unexpected behaviors.

If you are interested in the details of this decision, you can read the full discussion [here](https://github.com/weppos/publicsuffix-go/issues/31).


## Differences with `golang.org/x/net/publicsuffix`

The [`golang.org/x/net/publicsuffix`](https://godoc.org/golang.org/x/net/publicsuffix) is a package part of the Golang `x/net` package, that provides a public suffix list implementation.

The main difference is that the `x/net` package is optimized for speed, but it's less flexible. The list is compiled and embedded into the package itself. However, this is also the main downside.
The [list is not frequently refreshed](https://github.com/letsencrypt/boulder/issues/1374#issuecomment-182429297), hence the results may be inaccurate, in particular if you heavily rely on the private domain section of the list. Changes in the IANA section are less frequent, whereas changes in the Private Domains section happens weekly.

This package provides the following extra features:

- Ability to load an arbitrary list at runtime (e.g. you can feed your own list, or create multiple lists)
- Ability to create multiple lists
- Ability to parse a domain using a previously defined list
- Ability to add custom rules to an existing list, or merge/load rules from other lists (provided as file or string)
- Advanced access to the list rules
- Ability to ignore private domains at runtime, or when the list is parsed

This package also aims for 100% compatibility with the `x/net` package. A special adapter is provided as a drop-in replacement. Simply change the include statement from

```go
import (
    "golang.org/x/net/publicsuffix"
)
```

to

```go
import (
    "github.com/weppos/publicsuffix-go/net/publicsuffix"
)
```

The `github.com/weppos/publicsuffix-go/net/publicsuffix` package defines the same methods defined in `golang.org/x/net/publicsuffix`, but these methods are implemented using the `github.com/weppos/publicsuffix-go/publicsuffix` package.

Note that the adapter doesn't offer the flexibility of `github.com/weppos/publicsuffix-go/publicsuffix`, such as the ability to use multiple lists or disable private domains at runtime.


## Using with `cookiejar.PublicSuffixList`

This package implements the [`cookiejar.PublicSuffixList` interface](https://godoc.org/net/http/cookiejar#PublicSuffixList). It means it can be used as a value for the `PublicSuffixList` option when creating a `net/http/cookiejar`.

```go
import (
    "net/http/cookiejar"
    "github.com/weppos/publicsuffix-go/publicsuffix"
)

deliciousJar := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.CookieJarList})
```


## What is the public suffix list?

The [Public Suffix List](https://publicsuffix.org) is a cross-vendor initiative to provide an accurate list of domain name suffixes.

The Public Suffix List is an initiative of the Mozilla Project, but is maintained as a community resource. It is available for use in any software, but was originally created to meet the needs of browser manufacturers.

A "public suffix" is one under which Internet users can directly register names. Some examples of public suffixes are ".com", ".co.uk" and "pvt.k12.wy.us". The Public Suffix List is a list of all known public suffixes.


## Why use the public suffix list instead of regular expressions?

Previously, browsers used an algorithm which basically only denied setting wide-ranging cookies for top-level domains with no dots (e.g. com or org). However, this did not work for top-level domains where only third-level registrations are allowed (e.g. co.uk). In these cases, websites could set a cookie for co.uk which will be passed onto every website registered under co.uk.

Clearly, this was a security risk as it allowed websites other than the one setting the cookie to read it, and therefore potentially extract sensitive information.

Since there is no algorithmic method of finding the highest level at which a domain may be registered for a particular top-level domain (the policies differ with each registry), the only method is to create a list of all top-level domains and the level at which domains can be registered. This is the aim of the effective TLD list.

As well as being used to prevent cookies from being set where they shouldn't be, the list can also potentially be used for other applications where the registry controlled and privately controlled parts of a domain name need to be known, for example when grouping by top-level domains.

Source: https://wiki.mozilla.org/Public_Suffix_List

Not convinced yet? Check out [this real world example](https://stackoverflow.com/q/288810/123527).


## Does PublicSuffix make network requests?

No. `PublicSuffix` comes with a bundled list. It does not make any HTTP requests to parse or validate a domain.


## Terminology

- **TLD** (Top-Level Domain): The last segment of a domain name. For example, in `mozilla.org`, the `.org` portion is the TLD.

- **SLD** (Second-Level Domain): A domain directly below a top-level domain. For example, in `https://www.mozilla.org/en-US/`, `mozilla` is the second-level domain of the `.org` TLD.

- **TRD** (Third-Level Domain): Also known as a subdomain, this is the part of the domain before the SLD or root domain. For example, in `https://www.mozilla.org/en-US/`, `www` is the TRD.

- **FQDN** (Fully Qualified Domain Name): A complete domain name that includes the hostname, domain, and top-level domain, ending with a trailing dot. The format is `[hostname].[domain].[tld].` (e.g., `www.mozilla.org.`).


## Documentation and support

### Documentation

Library documentation is available at https://pkg.go.dev/github.com/weppos/publicsuffix-go/publicsuffix.

### Bug reports and contributions

- **Bug Tracker**: https://github.com/weppos/publicsuffix-go/issues
- **Code Repository**: https://github.com/weppos/publicsuffix-go

Contributions are welcome! Please include tests and/or feature coverage for every patch, and create a topic branch for every separate change you make.


## Security and vulnerability reporting

For full information and details about our security policy, please visit [`SECURITY.md`](SECURITY.md).


## Changelog

See [CHANGELOG.md](CHANGELOG.md) for details.


## License

Copyright (c) 2016-2025 Simone Carletti. This is Free Software distributed under the MIT license.

The [Public Suffix List source](https://publicsuffix.org/list/) is subject to the terms of the Mozilla Public License, v. 2.0.
