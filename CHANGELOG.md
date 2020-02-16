# CHANGELOG

## 0.13.0

- CHANGED: Rollback changes of v0.12.0. It turns out it is actually causing more issues.

## 0.12.0

- CHANGED: Extracted generator into its own package.

## 0.11.0

- CHANGED: Definition updates.

## 0.10.0

- ENHANCEMENT: Internal refactoring to use go gen when building definition list.

## 0.5.0

- FIXED: Added a DefaultRules() function that can be used to create a new list without modifying the default one (GH-141, GH-170). Thanks @guliyevemil1

- FIXED: Fixed nil pointer dereference when can't find a rule (GH-16) 

- CHANGED: Removed unreachable code (GH-167)

## 0.4.0

- CHANGED: Definition updates.

- ENHANCEMENT: gen tool now uses GitHub API instead of scraping GitHub UI (GH-93).

## 0.3.2

- CHANGED: Definition updates.

## 0.3.1

- CHANGED: Definition updates.

## 0.3.0

- CHANGED: Definition updates.

- ENHANCEMENT: Changed internal representation of PSL rules to be A-label encoded, as well the public interface of the library to use ASCII-encoded names by default (GH-31, GH-40).

## 0.2.0

- CHANGED: Definition updates.

- ENHANCEMENT: List.Select() is no longer exported. This was an experimental method and it's now kept private as the Find() implementation may change in the future.

- ENHANCEMENT: List.Find() now returns a pointer to a Rule, and not a Rule. That's because Find() can actually return `nil` if the DefaultRule find option is set. This is useful if you need to avoid the fallback to the default rule "*".

## 0.1.0

Initial version
