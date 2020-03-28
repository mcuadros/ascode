---
title: 'Built-in constants and functions'
weight: 9
toc: true
---

## Overview

The outermost block of the Starlark environment is known as the "predeclared" block.
It defines a number of fundamental values and functions needed by all Starlark programs,
such as `None`, `True`, `False`, and `len`, and possibly additional
application-specific names.

These names are not reserved words so Starlark programs are free to
redefine them in a smaller block such as a function body or even at
the top level of a module.  However, doing so may be confusing to the
reader.  Nonetheless, this rule permits names to be added to the
predeclared block in later versions of the language (or
application-specific dialect) without breaking existing programs.


## None

`None` is the distinguished value of the type `NoneType`.

## True and False

`True` and `False` are the two values of type `bool`.

## any

`any(x)` returns `True` if any element of the iterable sequence x has a truth value of true.
If the iterable is empty, it returns `False`.

## all

`all(x)` returns `False` if any element of the iterable sequence x has a truth value of false.
If the iterable is empty, it returns `True`.

## bool

`bool(x)` interprets `x` as a Boolean value---`True` or `False`.
With no argument, `bool()` returns `False`.


## chr

`chr(i)` returns a string that encodes the single Unicode code point
whose value is specified by the integer `i`. `chr` fails unless 0 ≤
`i` ≤ 0x10FFFF.

Example:

```python
chr(65)                         # "A",
chr(1049)                       # "Й", CYRILLIC CAPITAL LETTER SHORT I
chr(0x1F63F)                    # "😿", CRYING CAT FACE
```

See also: `ord`.

<b>Implementation note:</b> `chr` is not provided by the Java implementation.

## dict

`dict` creates a dictionary.  It accepts up to one positional
argument, which is interpreted as an iterable of two-element
sequences (pairs), each specifying a key/value pair in
the resulting dictionary.

`dict` also accepts any number of keyword arguments, each of which
specifies a key/value pair in the resulting dictionary;
each keyword is treated as a string.

```python
dict()                          # {}, empty dictionary
dict([(1, 2), (3, 4)])          # {1: 2, 3: 4}
dict([(1, 2), ["a", "b"]])      # {1: 2, "a": "b"}
dict(one=1, two=2)              # {"one": 1, "two", 1}
dict([(1, 2)], x=3)             # {1: 2, "x": 3}
```

With no arguments, `dict()` returns a new empty dictionary.

`dict(x)` where x is a dictionary returns a new copy of x.

## dir

`dir(x)` returns a new sorted list of the names of the attributes (fields and methods) of its operand.
The attributes of a value `x` are the names `f` such that `x.f` is a valid expression.

For example,

```python
dir("hello")                    # ['capitalize', 'count', ...], the methods of a string
```

Several types known to the interpreter, such as list, string, and dict, have methods, but none have fields.
However, an application may define types with fields that may be read or set by statements such as these:

```text
y = x.f
x.f = y
```

## enumerate

`enumerate(x)` returns a list of (index, value) pairs, each containing
successive values of the iterable sequence xand the index of the value
within the sequence.

The optional second parameter, `start`, specifies an integer value to
add to each index.

```python
enumerate(["zero", "one", "two"])               # [(0, "zero"), (1, "one"), (2, "two")]
enumerate(["one", "two"], 1)                    # [(1, "one"), (2, "two")]
```

## fail

The `fail(*args, sep=" ")` function causes execution to fail
with the specified error message.
Like `print`, arguments are formatted as if by `str(x)` and
separated by a space, unless an alternative separator is
specified by a `sep` named argument.

```python
fail("oops")				# "fail: oops"
fail("oops", 1, False, sep='/')		# "fail: oops/1/False"
```

## float

`float(x)` interprets its argument as a floating-point number.

If x is a `float`, the result is x.
if x is an `int`, the result is the nearest floating point value to x.
If x is a string, the string is interpreted as a floating-point literal.
With no arguments, `float()` returns `0.0`.

<b>Implementation note:</b>
Floating-point numbers are an optional feature.
The Go implementation of Starlark requires the `-float` flag to
enable support for floating-point literals, the `float` built-in
function, and the real division operator `/`.
The Java implementation does not yet support floating-point numbers.


## getattr

`getattr(x, name)` returns the value of the attribute (field or method) of x named `name`.
It is a dynamic error if x has no such attribute.

`getattr(x, "f")` is equivalent to `x.f`.

```python
getattr("banana", "split")("a")	       # ["b", "n", "n", ""], equivalent to "banana".split("a")
```

The three-argument form `getattr(x, name, default)` returns the
provided `default` value instead of failing.

## hasattr

`hasattr(x, name)` reports whether x has an attribute (field or method) named `name`.

## hash

`hash(x)` returns an integer hash of a string x
such that two equal strings have the same hash.
In other words `x == y` implies `hash(x) == hash(y)`.

In the interests of reproducibility of Starlark program behavior over time and
across implementations, the specific hash function is the same as that implemented by
[java.lang.String.hashCode](https://docs.oracle.com/javase/7/docs/api/java/lang/String.html#hashCode),
a simple polynomial accumulator over the UTF-16 transcoding of the string:
 ```
s[0]*31^(n-1) + s[1]*31^(n-2) + ... + s[n-1]
```

`hash` fails if given a non-string operand,
even if the value is hashable and thus suitable as the key of dictionary.

## int

`int(x[, base])` interprets its argument as an integer.

If x is an `int`, the result is x.
If x is a `float`, the result is the integer value nearest to x,
truncating towards zero; it is an error if x is not finite (`NaN`,
`+Inf`, `-Inf`).
If x is a `bool`, the result is 0 for `False` or 1 for `True`.

If x is a string, it is interpreted as a sequence of digits in the
specified base, decimal by default.
If `base` is zero, x is interpreted like an integer literal, the base
being inferred from an optional base marker such as `0b`, `0o`, or
`0x` preceding the first digit.
Irrespective of base, the string may start with an optional `+` or `-`
sign indicating the sign of the result.

## len

`len(x)` returns the number of elements in its argument.

It is a dynamic error if its argument is not a sequence.

## list

`list` constructs a list.

`list(x)` returns a new list containing the elements of the
iterable sequence x.

With no argument, `list()` returns a new empty list.

## max

`max(x)` returns the greatest element in the iterable sequence x.

It is an error if any element does not support ordered comparison,
or if the sequence is empty.

The optional named parameter `key` specifies a function to be applied
to each element prior to comparison.

```python
max([3, 1, 4, 1, 5, 9])                         # 9
max("two", "three", "four")                     # "two", the lexicographically greatest
max("two", "three", "four", key=len)            # "three", the longest
```

## min

`min(x)` returns the least element in the iterable sequence x.

It is an error if any element does not support ordered comparison,
or if the sequence is empty.

```python
min([3, 1, 4, 1, 5, 9])                         # 1
min("two", "three", "four")                     # "four", the lexicographically least
min("two", "three", "four", key=len)            # "two", the shortest
```


## ord

`ord(s)` returns the integer value of the sole Unicode code point encoded by the string `s`.

If `s` does not encode exactly one Unicode code point, `ord` fails.
Each invalid code within the string is treated as if it encodes the
Unicode replacement character, U+FFFD.

Example:

```python
ord("A")				# 65
ord("Й")				# 1049
ord("😿")					# 0x1F63F
ord("Й"[1:])				# 0xFFFD (Unicode replacement character)
```

See also: `chr`.

<b>Implementation note:</b> `ord` is not provided by the Java implementation.

## print

`print(*args, sep=" ")` prints its arguments, followed by a newline.
Arguments are formatted as if by `str(x)` and separated with a space,
unless an alternative separator is specified by a `sep` named argument.

Example:

```python
print(1, "hi")		       		# "1 hi\n"
print("hello", "world")			# "hello world\n"
print("hello", "world", sep=", ")	# "hello, world\n"
```

Typically the formatted string is printed to the standard error file,
but the exact behavior is a property of the Starlark thread and is
determined by the host application.

## range

`range` returns an immutable sequence of integers defined by the specified interval and stride.

```python
range(stop)                             # equivalent to range(0, stop)
range(start, stop)                      # equivalent to range(start, stop, 1)
range(start, stop, step)
```

`range` requires between one and three integer arguments.
With one argument, `range(stop)` returns the ascending sequence of non-negative integers less than `stop`.
With two arguments, `range(start, stop)` returns only integers not less than `start`.

With three arguments, `range(start, stop, step)` returns integers
formed by successively adding `step` to `start` until the value meets or passes `stop`.
A call to `range` fails if the value of `step` is zero.

A call to `range` does not materialize the entire sequence, but
returns a fixed-size value of type `"range"` that represents the
parameters that define the sequence.
The `range` value is iterable and may be indexed efficiently.

```python
list(range(10))                         # [0, 1, 2, 3, 4, 5, 6, 7, 8, 9]
list(range(3, 10))                      # [3, 4, 5, 6, 7, 8, 9]
list(range(3, 10, 2))                   # [3, 5, 7, 9]
list(range(10, 3, -2))                  # [10, 8, 6, 4]
```

The `len` function applied to a `range` value returns its length.
The truth value of a `range` value is `True` if its length is non-zero.

Range values are comparable: two `range` values compare equal if they
denote the same sequence of integers, even if they were created using
different parameters.

Range values are not hashable.  <!-- should they be? -->

The `str` function applied to a `range` value yields a string of the
form `range(10)`, `range(1, 10)`, or `range(1, 10, 2)`.

The `x in y` operator, where `y` is a range, reports whether `x` is equal to
some member of the sequence `y`; the operation fails unless `x` is a
number.

## repr

`repr(x)` formats its argument as a string.

All strings in the result are double-quoted.

```python
repr(1)                 # '1'
repr("x")               # '"x"'
repr([1, "x"])          # '[1, "x"]'
```

## reversed

`reversed(x)` returns a new list containing the elements of the iterable sequence x in reverse order.

```python
reversed(range(5))                              # [4, 3, 2, 1, 0]
reversed("stressed".codepoints())               # ["d", "e", "s", "s", "e", "r", "t", "s"]
reversed({"one": 1, "two": 2}.keys())           # ["two", "one"]
```

## set

`set(x)` returns a new set containing the elements of the iterable x.
With no argument, `set()` returns a new empty set.

```python
set([3, 1, 4, 1, 5, 9])         # set([3, 1, 4, 5, 9])
```

<b>Implementation note:</b>
Sets are an optional feature of the Go implementation of Starlark,
enabled by the `-set` flag.


## sorted

`sorted(x)` returns a new list containing the elements of the iterable sequence x,
in sorted order.  The sort algorithm is stable.

The optional named parameter `reverse`, if true, causes `sorted` to
return results in reverse sorted order.

The optional named parameter `key` specifies a function of one
argument to apply to obtain the value's sort key.
The default behavior is the identity function.

```python
sorted(set("harbors".codepoints()))                             # ['a', 'b', 'h', 'o', 'r', 's']
sorted([3, 1, 4, 1, 5, 9])                                      # [1, 1, 3, 4, 5, 9]
sorted([3, 1, 4, 1, 5, 9], reverse=True)                        # [9, 5, 4, 3, 1, 1]

sorted(["two", "three", "four"], key=len)                       # ["two", "four", "three"], shortest to longest
sorted(["two", "three", "four"], key=len, reverse=True)         # ["three", "four", "two"], longest to shortest
```


## str

`str(x)` formats its argument as a string.

If x is a string, the result is x (without quotation).
All other strings, such as elements of a list of strings, are double-quoted.

```python
str(1)                          # '1'
str("x")                        # 'x'
str([1, "x"])                   # '[1, "x"]'
```

## tuple

`tuple(x)` returns a tuple containing the elements of the iterable x.

With no arguments, `tuple()` returns the empty tuple.

## type

type(x) returns a string describing the type of its operand.

```python
type(None)              # "NoneType"
type(0)                 # "int"
type(0.0)               # "float"
```

## zip

`zip()` returns a new list of n-tuples formed from corresponding
elements of each of the n iterable sequences provided as arguments to
`zip`.  That is, the first tuple contains the first element of each of
the sequences, the second element contains the second element of each
of the sequences, and so on.  The result list is only as long as the
shortest of the input sequences.

```python
zip()                                   # []
zip(range(5))                           # [(0,), (1,), (2,), (3,), (4,)]
zip(range(5), "abc")                    # [(0, "a"), (1, "b"), (2, "c")]
```
