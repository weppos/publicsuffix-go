# CHANGELOG

#### 0.3.2

- CHANGED: Definition updates.

#### 0.3.1

- CHANGED: Definition updates.

#### 0.3.0

- CHANGED: Definition updates.

- CHANGED: Changed internal representation of PSL rules to be A-label encoded, as well the public interface of the library to use ASCII-encoded names by default (GH-31, GH-40).

#### 0.2.0

- CHANGED: Definition updates.

- CHANGED: List.Select() is no longer exported. This was an experimental method and it's now kept private as the Find() implementation may change in the future.

- CHANGED: List.Find() now returns a pointer to a Rule, and not a Rule. That's because Find() can actually return `nil` if the DefaultRule find option is set. This is useful if you need to avoid the fallback to the default rule "*".

#### 0.1.0

Initial version
