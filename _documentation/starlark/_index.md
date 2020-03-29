---
title: 'Language definition'
weight: 10
---

Starlark is a dialect of Python intended for use as a configuration
language.  A Starlark interpreter is typically embedded within a larger
application, and this application may define additional
domain-specific functions and data types beyond those provided by the
core language.  For example, Starlark is embedded within (and was
originally developed for) the [Bazel build tool](https://bazel.build),
and [Bazel's build language](https://docs.bazel.build/versions/2.0.0/skylark/language.html) is based on Starlark.

This document describes the Go implementation of Starlark
at go.starlark.net/starlark.
The language it defines is similar but not identical to
[the Java-based implementation](https://github.com/bazelbuild/starlark)
used by Bazel.
We identify places where their behaviors differ, and an
[appendix](/docs/starlark/dialect-differences/) provides a summary of those
differences.
We plan to converge both implementations on a single specification.

This document is maintained by Alan Donovan <adonovan@google.com>.
It was influenced by the Python specification,
Copyright 1990&ndash;2017, Python Software Foundation,
and the Go specification, Copyright 2009&ndash;2017, The Go Authors.

Starlark was designed and implemented in Java by Laurent Le Brun,
Dmitry Lomov, Jon Brandvin, and Damien Martin-Guillerez, standing on
the shoulders of the Python community.
The Go implementation was written by Alan Donovan and Jay Conrod;
its scanner was derived from one written by Russ Cox.

## Overview

Starlark is an untyped dynamic language with high-level data types,
first-class functions with lexical scope, and automatic memory
management or _garbage collection_.

Starlark is strongly influenced by Python, and is almost a subset of
that language.  In particular, its data types and syntax for
statements and expressions will be very familiar to any Python
programmer.
However, Starlark is intended not for writing applications but for
expressing configuration: its programs are short-lived and have no
external side effects and their main result is structured data or side
effects on the host application.
As a result, Starlark has no need for classes, exceptions, reflection,
concurrency, and other such features of Python.

Starlark execution is _deterministic_: all functions and operators
in the core language produce the same execution each time the program
is run; there are no sources of random numbers, clocks, or unspecified
iterators. This makes Starlark suitable for use in applications where
reproducibility is paramount, such as build tools.
