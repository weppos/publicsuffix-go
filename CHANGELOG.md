# Changelog


## 0.50.2 - 2026-01-03

### Changed

- Definition updates.


## 0.50.1 - 2024-11-17

### Changed

- Definition updates.


## 0.50.0 - 2024-08-22

### Changed

- Removed dependency from go-github #1067
- Supported minimum version is now Go 1.23.
- Definition updates.


## 0.40.2 - 2024-07-23

### Changed

- Definition updates.


## 0.40.0 - 2024-06-15

### Changed

- Supported minimum version is now Go 1.21.
- Definition updates.


## 0.30.2 - 2024-03-18

### Changed

- Definition updates.


## 0.30.1 - 2023-07-11

### Changed

- Definition updates.


## 0.30.0 - 2023-02-20

### Changed

- Dropped Go < 1.16 from the list of supported versions to fix "package embed is not in GOROOT". "embed" is an indirect dependency and it's only available since Go 1.16.
- Exported defaultListVersion as ListVersion #334, #880


## 0.20.0 - 2022-07-26

### Changed

- Definition updates.


## 0.15.0 - 2021-04-14

### Changed

- Definition updates.
- Dropped Go 1.8 from the list of supported versions. "math/bits" is an indirect dependency and it's only available since Go 1.9.
- Improved performances by using rune instead of strings single char comparison #484, #485


## 0.14.0 - 2020-02-16

### Changed

- Added go modules #240.


## 0.13.0 - 2020-02-16

### Changed

- Rollback changes of v0.12.0. It turns out it is actually causing more issues.


## 0.12.0 - 2020-02-16

### Changed

- Extracted generator into its own package.


## 0.11.0 - 2020-02-16

### Changed

- Definition updates.


## 0.10.0 - 2019-08-08

### Changed

- Internal refactoring to use go gen when building definition list.


## 0.5.0 - 2019-06-04

### Fixed

- Added a DefaultRules() function that can be used to create a new list without modifying the default one #141, #170. (Thanks @guliyevemil1)
- Fixed nil pointer dereference when can't find a rule #16

### Changed

- Removed unreachable code #167


## 0.4.0 - 2018-03-30

### Changed

- Definition updates.
- gen tool now uses GitHub API instead of scraping GitHub UI #93.


## 0.3.2 - 2017-02-07

### Changed

- Definition updates.


## 0.3.1 - 2017-01-02

### Changed

- Definition updates.


## 0.3.0 - 2016-11-21

### Changed

- Definition updates.
- Changed internal representation of PSL rules to be A-label encoded, as well the public interface of the library to use ASCII-encoded names by default #31, #40.


## 0.2.0 - 2016-11-20

### Changed

- Definition updates.
- List.Select() is no longer exported. This was an experimental method and it's now kept private as the Find() implementation may change in the future.
- List.Find() now returns a pointer to a Rule, and not a Rule. That's because Find() can actually return `nil` if the DefaultRule find option is set. This is useful if you need to avoid the fallback to the default rule "*".


## 0.1.0 - 2016-06-25

Initial version
