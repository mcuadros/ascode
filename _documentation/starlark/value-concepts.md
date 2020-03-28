---
title: 'Value concepts'
weight: 5
toc: true
---

## Overview

Starlark has eleven core [data types](#data-types).  An application
that embeds the Starlark intepreter may define additional types that
behave like Starlark values.  All values, whether core or
application-defined, implement a few basic behaviors:

```text
str(x)		-- return a string representation of x
type(x)		-- return a string describing the type of x
bool(x)		-- convert x to a Boolean truth value
```

## Identity and mutation

Starlark is an imperative language: programs consist of sequences of
statements executed for their side effects.
For example, an assignment statement updates the value held by a
variable, and calls to some built-in functions such as `print` change
the state of the application that embeds the interpreter.

Values of some data types, such as `NoneType`, `bool`, `int`, `float`, and
`string`, are _immutable_; they can never change.
Immutable values have no notion of _identity_: it is impossible for a
Starlark program to tell whether two integers, for instance, are
represented by the same object; it can tell only whether they are
equal.

Values of other data types, such as `list`, `dict`, and `set`, are
_mutable_: they may be modified by a statement such as `a[i] = 0` or
`items.clear()`.  Although `tuple` and `function` values are not
directly mutable, they may refer to mutable values indirectly, so for
this reason we consider them mutable too.  Starlark values of these
types are actually _references_ to variables.

Copying a reference to a variable, using an assignment statement for
instance, creates an _alias_ for the variable, and the effects of
operations applied to the variable through one alias are visible
through all others.

```python
x = []                          # x refers to a new empty list variable
y = x                           # y becomes an alias for x
x.append(1)                     # changes the variable referred to by x
print(y)                        # "[1]"; y observes the mutation
```

Starlark uses _call-by-value_ parameter passing: in a function call,
argument values are assigned to function parameters as if by
assignment statements.  If the values are references, the caller and
callee may refer to the same variables, so if the called function
changes the variable referred to by a parameter, the effect may also
be observed by the caller:

```python
def f(y):
    y.append(1)                 # changes the variable referred to by x

x = []                          # x refers to a new empty list variable
f(x)                            # f's parameter y becomes an alias for x
print(x)                        # "[1]"; x observes the mutation
```


As in all imperative languages, understanding _aliasing_, the
relationship between reference values and the variables to which they
refer, is crucial to writing correct programs.

## Freezing a value

Starlark has a feature unusual among imperative programming languages:
a mutable value may be _frozen_ so that all subsequent attempts to
mutate it fail with a dynamic error; the value, and all other values
reachable from it, become _immutable_.

Immediately after execution of a Starlark module, all values in its
top-level environment are frozen. Because all the global variables of
an initialized Starlark module are immutable, the module may be published to
and used by other threads in a parallel program without the need for
locks. For example, the Bazel build system loads and executes BUILD
and .bzl files in parallel, and two modules being executed
concurrently may freely access variables or call functions from a
third without the possibility of a race condition.

## Hashing

The `dict` and `set` data types are implemented using hash tables, so
only _hashable_ values are suitable as keys of a `dict` or elements of
a `set`. Attempting to use a non-hashable value as the key in a hash
table results in a dynamic error.

The hash of a value is an unspecified integer chosen so that two equal
values have the same hash, in other words, `x == y => hash(x) == hash(y)`.
A hashable value has the same hash throughout its lifetime.

Values of the types `NoneType`, `bool`, `int`, `float`, and `string`,
which are all immutable, are hashable.

Values of mutable types such as `list`, `dict`, and `set` are not
hashable. These values remain unhashable even if they have become
immutable due to _freezing_.

A `tuple` value is hashable only if all its elements are hashable.
Thus `("localhost", 80)` is hashable but `([127, 0, 0, 1], 80)` is not.

Values of the types `function` and `builtin_function_or_method` are also hashable.
Although functions are not necessarily immutable, as they may be
closures that refer to mutable variables, instances of these types
are compared by reference identity (see [Comparisons](#comparisons)),
so their hash values are derived from their identity.


## Sequence types

Many Starlark data types represent a _sequence_ of values: lists,
tuples, and sets are sequences of arbitrary values, and in many
contexts dictionaries act like a sequence of their keys.

We can classify different kinds of sequence types based on the
operations they support.
Each is listed below using the name of its corresponding interface in
the interpreter's Go API.

* `Iterable`: an _iterable_ value lets us process each of its elements in a fixed order.
  Examples: `dict`, `set`, `list`, `tuple`, but not `string`.
* `Sequence`: a _sequence of known length_ lets us know how many elements it
  contains without processing them.
  Examples: `dict`, `set`, `list`, `tuple`, but not `string`.
* `Indexable`: an _indexed_ type has a fixed length and provides efficient
  random access to its elements, which are identified by integer indices.
  Examples: `string`, `tuple`, and `list`.
* `SetIndexable`: a _settable indexed type_ additionally allows us to modify the
  element at a given integer index. Example: `list`.
* `Mapping`: a mapping is an association of keys to values. Example: `dict`.

Although all of Starlark's core data types for sequences implement at
least the `Sequence` contract, it's possible for an application
that embeds the Starlark interpreter to define additional data types
representing sequences of unknown length that implement only the `Iterable` contract.

Strings are not iterable, though they do support the `len(s)` and
`s[i]` operations. Starlark deviates from Python here to avoid a common
pitfall in which a string is used by mistake where a list containing a
single string was intended, resulting in its interpretation as a sequence
of bytes.

Most Starlark operators and built-in functions that need a sequence
of values will accept any iterable.

It is a dynamic error to mutate a sequence such as a list, set, or
dictionary while iterating over it.

```python
def increment_values(dict):
  for k in dict:
    dict[k] += 1			# error: cannot insert into hash table during iteration

dict = {"one": 1, "two": 2}
increment_values(dict)
```


## Indexing

Many Starlark operators and functions require an index operand `i`,
such as `a[i]` or `list.insert(i, x)`. Others require two indices `i`
and `j` that indicate the start and end of a sub-sequence, such as
`a[i:j]`, `list.index(x, i, j)`, or `string.find(x, i, j)`.
All such operations follow similar conventions, described here.

Indexing in Starlark is *zero-based*. The first element of a string
or list has index 0, the next 1, and so on. The last element of a
sequence of length `n` has index `n-1`.

```python
"hello"[0]			# "h"
"hello"[4]			# "o"
"hello"[5]			# error: index out of range
```

For sub-sequence operations that require two indices, the first is
_inclusive_ and the second _exclusive_. Thus `a[i:j]` indicates the
sequence starting with element `i` up to but not including element
`j`. The length of this sub-sequence is `j-i`. This convention is known
as *half-open indexing*.

```python
"hello"[1:4]			# "ell"
```

Either or both of the index operands may be omitted. If omitted, the
first is treated equivalent to 0 and the second is equivalent to the
length of the sequence:

```python
"hello"[1:]                     # "ello"
"hello"[:4]                     # "hell"
```

It is permissible to supply a negative integer to an indexing
operation. The effective index is computed from the supplied value by
the following two-step procedure. First, if the value is negative, the
length of the sequence is added to it. This provides a convenient way
to address the final elements of the sequence:

```python
"hello"[-1]                     # "o",  like "hello"[4]
"hello"[-3:-1]                  # "ll", like "hello"[2:4]
```

Second, for sub-sequence operations, if the value is still negative, it
is replaced by zero, or if it is greater than the length `n` of the
sequence, it is replaced by `n`. In effect, the index is "truncated" to
the nearest value in the range `[0:n]`.

```python
"hello"[-1000:+1000]		# "hello"
```

This truncation step does not apply to indices of individual elements:

```python
"hello"[-6]		# error: index out of range
"hello"[-5]		# "h"
"hello"[4]		# "o"
"hello"[5]		# error: index out of range
```

