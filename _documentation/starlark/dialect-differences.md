---
title: 'Dialect differences'
weight: 11
---

The list below summarizes features of the Go implementation that are
known to differ from the Java implementation of Starlark used by Bazel.
Some of these features may be controlled by global options to allow
applications to mimic the Bazel dialect more closely. Our goal is
eventually to eliminate all such differences on a case-by-case basis.
See [Starlark spec issue 20](https://github.com/bazelbuild/starlark/issues/20).

* Integers are represented with infinite precision.
* Integer arithmetic is exact.
* Floating-point literals are supported (option: `-float`).
* The `float` built-in function is provided (option: `-float`).
* Real division using `float / float` is supported (option: `-float`).
* String interpolation supports the `[ioxXeEfFgGc]` conversions.
* `def` statements may be nested (option: `-nesteddef`).
* `lambda` expressions are supported (option: `-lambda`).
* String elements are bytes.
* Non-ASCII strings are encoded using UTF-8.
* Strings support octal and hex byte escapes.
* Strings have the additional methods `elem_ords`, `codepoint_ords`, and `codepoints`.
* The `chr` and `ord` built-in functions are supported.
* The `set` built-in function is provided (option: `-set`).
* `set & set` and `set | set` compute set intersection and union, respectively.
* `assert` is a valid identifier.
* Dot expressions may appear on the left side of an assignment: `x.f = 1`.
* `type(x)` returns `"builtin_function_or_method"` for built-in functions.
* `if`, `for`, and `while` are permitted at top level (option: `-globalreassign`).
* top-level rebindings are permitted (option: `-globalreassign`).