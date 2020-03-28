---
title: 'Data types'
weight: 3
toc: true
---

## Overview

These are the main data types built in to the interpreter:

```python
NoneType                     # the type of None
bool                         # True or False
int                          # a signed integer of arbitrary magnitude
float                        # an IEEE 754 double-precision floating point number
string                       # a byte string
list                         # a modifiable sequence of values
tuple                        # an unmodifiable sequence of values
dict                         # a mapping from values to values
set                          # a set of values
function                     # a function implemented in Starlark
builtin_function_or_method   # a function or method implemented by the interpreter or host application
```

Some functions, such as the iteration methods of `string`, or the
`range` function, return instances of special-purpose types that don't
appear in this list.
Additional data types may be defined by the host application into
which the interpreter is embedded, and those data types may
participate in basic operations of the language such as arithmetic,
comparison, indexing, and function calls.

<!-- We needn't mention the stringIterable type here. -->

Some operations can be applied to any Starlark value.  For example,
every value has a type string that can be obtained with the expression
`type(x)`, and any value may be converted to a string using the
expression `str(x)`, or to a Boolean truth value using the expression
`bool(x)`.  Other operations apply only to certain types.  For
example, the indexing operation `a[i]` works only with strings, lists,
and tuples, and any application-defined types that are _indexable_.
The [_value concepts_](#value-concepts) section explains the groupings of
types by the operators they support.


## None

`None` is a distinguished value used to indicate the absence of any other value.
For example, the result of a call to a function that contains no return statement is `None`.

`None` is equal only to itself.  Its [type](#type) is `"NoneType"`.
The truth value of `None` is `False`.


## Booleans

There are two Boolean values, `True` and `False`, representing the
truth or falsehood of a predicate.  The [type](#type) of a Boolean is `"bool"`.

Boolean values are typically used as conditions in `if`-statements,
although any Starlark value used as a condition is implicitly
interpreted as a Boolean.
For example, the values `None`, `0`, `0.0`, and the empty sequences
`""`, `()`, `[]`, and `{}` have a truth value of `False`, whereas non-zero
numbers and non-empty sequences have a truth value of `True`.
Application-defined types determine their own truth value.
Any value may be explicitly converted to a Boolean using the built-in `bool`
function.

```python
1 + 1 == 2                              # True
2 + 2 == 5                              # False

if 1 + 1:
        print("True")
else:
        print("False")
```

## Integers

The Starlark integer type represents integers.  Its [type](#type) is `"int"`.

Integers may be positive or negative, and arbitrarily large.
Integer arithmetic is exact.
Integers are totally ordered; comparisons follow mathematical
tradition.

The `+` and `-` operators perform addition and subtraction, respectively.
The `*` operator performs multiplication.

The `//` and `%` operations on integers compute floored division and
remainder of floored division, respectively.
If the signs of the operands differ, the sign of the remainder `x % y`
matches that of the dividend, `x`.
For all finite x and y (y ≠ 0), `(x // y) * y + (x % y) == x`.
The `/` operator implements real division, and
yields a `float` result even when its operands are both of type `int`.

Integers, including negative values, may be interpreted as bit vectors.
The `|`, `&`, and `^` operators implement bitwise OR, AND, and XOR,
respectively. The unary `~` operator yields the bitwise inversion of its
integer argument. The `<<` and `>>` operators shift the first argument
to the left or right by the number of bits given by the second argument.

Any bool, number, or string may be interpreted as an integer by using
the `int` built-in function.

An integer used in a Boolean context is considered true if it is
non-zero.

```python
100 // 5 * 9 + 32               # 212
3 // 2                          # 1
3 / 2                           # 1.5
111111111 * 111111111           # 12345678987654321
"0x%x" % (0x1234 & 0xf00f)      # "0x1004"
int("ffff", 16)                 # 65535, 0xffff
```

<b>Implementation note:</b>
In the Go implementation of Starlark, integer representation and
arithmetic is exact, motivated by the need for lossless manipulation
of protocol messages which may contain signed and unsigned 64-bit
integers.
The Java implementation currently supports only signed 32-bit integers.


## Floating-point numbers

The Starlark floating-point data type represents an IEEE 754
double-precision floating-point number.  Its [type](#type) is `"float"`.

Arithmetic on floats using the `+`, `-`, `*`, `/`, `//`, and `%`
 operators follows the IEE 754 standard.
However, computing the division or remainder of division by zero is a dynamic error.

An arithmetic operation applied to a mixture of `float` and `int`
operands works as if the `int` operand is first converted to a
`float`.  For example, `3.141 + 1` is equivalent to `3.141 +
float(1)`.
There are two floating-point division operators:
`x / y ` yields the floating-point quotient of `x` and `y`,
whereas `x // y` yields `floor(x / y)`, that is, the largest
integer value not greater than `x / y`.
Although the resulting number is integral, it is represented as a
`float` if either operand is a `float`.

The infinite float values `+Inf` and `-Inf` represent numbers
greater/less than all finite float values.

The non-finite `NaN` value represents the result of dubious operations
such as `Inf/Inf`.  A NaN value compares neither less than, nor
greater than, nor equal to any value, including itself.

All floats other than NaN are totally ordered, so they may be compared
using operators such as `==` and `<`.

Any bool, number, or string may be interpreted as a floating-point
number by using the `float` built-in function.

A float used in a Boolean context is considered true if it is
non-zero.

```python
1.23e45 * 1.23e45                               # 1.5129e+90
1.111111111111111 * 1.111111111111111           # 1.23457
3.0 / 2                                         # 1.5
3 / 2.0                                         # 1.5
float(3) / 2                                    # 1.5
3.0 // 2.0                                      # 1
```

<b>Implementation note:</b>
The Go implementation of Starlark supports floating-point numbers as an
optional feature, motivated by the need for lossless manipulation of
protocol messages.
The `-float` flag enables support for floating-point literals,
the `float` built-in function, and the real division operator `/`.
The Java implementation does not yet support floating-point numbers.


## Strings

A string represents an immutable sequence of bytes.
The [type](#type) of a string is `"string"`.

Strings can represent arbitrary binary data, including zero bytes, but
most strings contain text, encoded by convention using UTF-8.

The built-in `len` function returns the number of bytes in a string.

Strings may be concatenated with the `+` operator.

The substring expression `s[i:j]` returns the substring of `s` from
index `i` up to index `j`.  The index expression `s[i]` returns the
1-byte substring `s[i:i+1]`.

Strings are hashable, and thus may be used as keys in a dictionary.

Strings are totally ordered lexicographically, so strings may be
compared using operators such as `==` and `<`.

Strings are _not_ iterable sequences, so they cannot be used as the operand of
a `for`-loop, list comprehension, or any other operation than requires
an iterable sequence.
To obtain a view of a string as an iterable sequence of numeric byte
values, 1-byte substrings, numeric Unicode code points, or 1-code
point substrings, you must explicitly call one of its four methods:
`elems`, `elem_ords`, `codepoints`, or `codepoint_ords`.

Any value may formatted as a string using the `str` or `repr` built-in
functions, the `str % tuple` operator, or the `str.format` method.

A string used in a Boolean context is considered true if it is
non-empty.

Strings have several built-in methods:

* [`capitalize`](#string·capitalize)
* [`codepoint_ords`](#string·codepoint_ords)
* [`codepoints`](#string·codepoints)
* [`count`](#string·count)
* [`elem_ords`](#string·elem_ords)
* [`elems`](#string·elems)
* [`endswith`](#string·endswith)
* [`find`](#string·find)
* [`format`](#string·format)
* [`index`](#string·index)
* [`isalnum`](#string·isalnum)
* [`isalpha`](#string·isalpha)
* [`isdigit`](#string·isdigit)
* [`islower`](#string·islower)
* [`isspace`](#string·isspace)
* [`istitle`](#string·istitle)
* [`isupper`](#string·isupper)
* [`join`](#string·join)
* [`lower`](#string·lower)
* [`lstrip`](#string·lstrip)
* [`partition`](#string·partition)
* [`replace`](#string·replace)
* [`rfind`](#string·rfind)
* [`rindex`](#string·rindex)
* [`rpartition`](#string·rpartition)
* [`rsplit`](#string·rsplit)
* [`rstrip`](#string·rstrip)
* [`split`](#string·split)
* [`splitlines`](#string·splitlines)
* [`startswith`](#string·startswith)
* [`strip`](#string·strip)
* [`title`](#string·title)
* [`upper`](#string·upper)

<b>Implementation note:</b>
The type of a string element varies across implementations.
There is agreement that byte strings, with text conventionally encoded
using UTF-8, is the ideal choice, but the Java implementation treats
strings as sequences of UTF-16 codes and changing it appears
intractible; see Google Issue b/36360490.

<b>Implementation note:</b>
The Java implementation does not consistently treat strings as
iterable; see `testdata/string.star` in the test suite and Google Issue
b/34385336 for further details.

## Lists

A list is a mutable sequence of values.
The [type](#type) of a list is `"list"`.

Lists are indexable sequences: the elements of a list may be iterated
over by `for`-loops, list comprehensions, and various built-in
functions.

List may be constructed using bracketed list notation:

```python
[]              # an empty list
[1]             # a 1-element list
[1, 2]          # a 2-element list
```

Lists can also be constructed from any iterable sequence by using the
built-in `list` function.

The built-in `len` function applied to a list returns the number of elements.
The index expression `list[i]` returns the element at index i,
and the slice expression `list[i:j]` returns a new list consisting of
the elements at indices from i to j.

List elements may be added using the `append` or `extend` methods,
removed using the `remove` method, or reordered by assignments such as
`list[i] = list[j]`.

The concatenation operation `x + y` yields a new list containing all
the elements of the two lists x and y.

For most types, `x += y` is equivalent to `x = x + y`, except that it
evaluates `x` only once, that is, it allocates a new list to hold
the concatenation of `x` and `y`.
However, if `x` refers to a list, the statement does not allocate a
new list but instead mutates the original list in place, similar to
`x.extend(y)`.

Lists are not hashable, so may not be used in the keys of a dictionary.

A list used in a Boolean context is considered true if it is
non-empty.

A [_list comprehension_](#comprehensions) creates a new list whose elements are the
result of some expression applied to each element of another sequence.

```python
[x*x for x in [1, 2, 3, 4]]      # [1, 4, 9, 16]
```

A list value has these methods:

* [`append`](#list·append)
* [`clear`](#list·clear)
* [`extend`](#list·extend)
* [`index`](#list·index)
* [`insert`](#list·insert)
* [`pop`](#list·pop)
* [`remove`](#list·remove)

## Tuples

A tuple is an immutable sequence of values.
The [type](#type) of a tuple is `"tuple"`.

Tuples are constructed using parenthesized list notation:

```python
()                      # the empty tuple
(1,)                    # a 1-tuple
(1, 2)                  # a 2-tuple ("pair")
(1, 2, 3)               # a 3-tuple
```

Observe that for the 1-tuple, the trailing comma is necessary to
distinguish it from the parenthesized expression `(1)`.
1-tuples are seldom used.

Starlark, unlike Python, does not permit a trailing comma to appear in
an unparenthesized tuple expression:

```python
for k, v, in dict.items(): pass                 # syntax error at 'in'
_ = [(v, k) for k, v, in dict.items()]          # syntax error at 'in'
f = lambda a, b, : None                         # syntax error at ':'

sorted(3, 1, 4, 1,)                             # ok
[1, 2, 3, ]                                     # ok
{1: 2, 3:4, }                                   # ok
```

Any iterable sequence may be converted to a tuple by using the
built-in `tuple` function.

Like lists, tuples are indexed sequences, so they may be indexed and
sliced.  The index expression `tuple[i]` returns the tuple element at
index i, and the slice expression `tuple[i:j]` returns a sub-sequence
of a tuple.

Tuples are iterable sequences, so they may be used as the operand of a
`for`-loop, a list comprehension, or various built-in functions.

Unlike lists, tuples cannot be modified.
However, the mutable elements of a tuple may be modified.

Tuples are hashable (assuming their elements are hashable),
so they may be used as keys of a dictionary.

Tuples may be concatenated using the `+` operator.

A tuple used in a Boolean context is considered true if it is
non-empty.


## Dictionaries

A dictionary is a mutable mapping from keys to values.
The [type](#type) of a dictionary is `"dict"`.

Dictionaries provide constant-time operations to insert an element, to
look up the value for a key, or to remove an element.  Dictionaries
are implemented using hash tables, so keys must be hashable.  Hashable
values include `None`, Booleans, numbers, and strings, and tuples
composed from hashable values.  Most mutable values, such as lists,
dictionaries, and sets, are not hashable, even when frozen.
Attempting to use a non-hashable value as a key in a dictionary
results in a dynamic error.

A [dictionary expression](#dictionary-expressions) specifies a
dictionary as a set of key/value pairs enclosed in braces:

```python
coins = {
  "penny": 1,
  "nickel": 5,
  "dime": 10,
  "quarter": 25,
}
```

The expression `d[k]`, where `d` is a dictionary and `k` is a key,
retrieves the value associated with the key.  If the dictionary
contains no such item, the operation fails:

```python
coins["penny"]          # 1
coins["dime"]           # 10
coins["silver dollar"]  # error: key not found
```

The number of items in a dictionary `d` is given by `len(d)`.
A key/value item may be added to a dictionary, or updated if the key
is already present, by using `d[k]` on the left side of an assignment:

```python
len(coins)				# 4
coins["shilling"] = 20
len(coins)				# 5, item was inserted
coins["shilling"] = 5
len(coins)				# 5, existing item was updated
```

A dictionary can also be constructed using a [dictionary
comprehension](#comprehension), which evaluates a pair of expressions,
the _key_ and the _value_, for every element of another iterable such
as a list.  This example builds a mapping from each word to its length
in bytes:

```python
words = ["able", "baker", "charlie"]
{x: len(x) for x in words}	# {"charlie": 7, "baker": 5, "able": 4}
```

Dictionaries are iterable sequences, so they may be used as the
operand of a `for`-loop, a list comprehension, or various built-in
functions.
Iteration yields the dictionary's keys in the order in which they were
inserted; updating the value associated with an existing key does not
affect the iteration order.

```python
x = dict([("a", 1), ("b", 2)])          # {"a": 1, "b": 2}
x.update([("a", 3), ("c", 4)])          # {"a": 3, "b": 2, "c": 4}
```

```python
for name in coins:
  print(name, coins[name])	# prints "quarter 25", "dime 10", ...
```

Like all mutable values in Starlark, a dictionary can be frozen, and
once frozen, all subsequent operations that attempt to update it will
fail.

A dictionary used in a Boolean context is considered true if it is
non-empty.

Dictionaries may be compared for equality using `==` and `!=`.  Two
dictionaries compare equal if they contain the same number of items
and each key/value item (k, v) found in one dictionary is also present
in the other.  Dictionaries are not ordered; it is an error to compare
two dictionaries with `<`.


A dictionary value has these methods:

* [`clear`](#dict·clear)
* [`get`](#dict·get)
* [`items`](#dict·items)
* [`keys`](#dict·keys)
* [`pop`](#dict·pop)
* [`popitem`](#dict·popitem)
* [`setdefault`](#dict·setdefault)
* [`update`](#dict·update)
* [`values`](#dict·values)

## Sets

A set is a mutable set of values.
The [type](#type) of a set is `"set"`.

Like dictionaries, sets are implemented using hash tables, so the
elements of a set must be hashable.

Sets may be compared for equality or inequality using `==` and `!=`.
Two sets compare equal if they contain the same elements.

Sets are iterable sequences, so they may be used as the operand of a
`for`-loop, a list comprehension, or various built-in functions.
Iteration yields the set's elements in the order in which they were
inserted.

The binary `|` and `&` operators compute union and intersection when
applied to sets.  The right operand of the `|` operator may be any
iterable value.  The binary `in` operator performs a set membership
test when its right operand is a set.

The binary `^` operator performs symmetric difference of two sets.

Sets are instantiated by calling the built-in `set` function, which
returns a set containing all the elements of its optional argument,
which must be an iterable sequence.  Sets have no literal syntax.

The only method of a set is `union`, which is equivalent to the `|` operator.

A set used in a Boolean context is considered true if it is non-empty.

<b>Implementation note:</b>
The Go implementation of Starlark requires the `-set` flag to
enable support for sets.
The Java implementation does not support sets.


## Functions

A function value represents a function defined in Starlark.
Its [type](#type) is `"function"`.
A function value used in a Boolean context is always considered true.

Functions defined by a [`def` statement](#function-definitions) are named;
functions defined by a [`lambda` expression](#lambda-expressions) are anonymous.

Function definitions may be nested, and an inner function may refer to a local variable of an outer function.

A function definition defines zero or more named parameters.
Starlark has a rich mechanism for passing arguments to functions.

<!-- TODO break up this explanation into caller-side and callee-side
     parts, and put the former under function calls and the latter
     under function definitions. Also try to convey that the Callable
     interface sees the flattened-out args and kwargs and that's what
     built-ins get.
-->

The example below shows a definition and call of a function of two
required parameters, `x` and `y`.

```python
def idiv(x, y):
  return x // y

idiv(6, 3)		# 2
```

A call may provide arguments to function parameters either by
position, as in the example above, or by name, as in first two calls
below, or by a mixture of the two forms, as in the third call below.
All the positional arguments must precede all the named arguments.
Named arguments may improve clarity, especially in functions of
several parameters.

```python
idiv(x=6, y=3)		# 2
idiv(y=3, x=6)		# 2

idiv(6, y=3)		# 2
```

<b>Optional parameters:</b> A parameter declaration may specify a
default value using `name=value` syntax; such a parameter is
_optional_.  The default value expression is evaluated during
execution of the `def` statement or evaluation of the `lambda`
expression, and the default value forms part of the function value.
All optional parameters must follow all non-optional parameters.
A function call may omit arguments for any suffix of the optional
parameters; the effective values of those arguments are supplied by
the function's parameter defaults.

```python
def f(x, y=3):
  return x, y

f(1, 2)	# (1, 2)
f(1)	# (1, 3)
```

If a function parameter's default value is a mutable expression,
modifications to the value during one call may be observed by
subsequent calls.
Beware of this when using lists or dicts as default values.
If the function becomes frozen, its parameters' default values become
frozen too.

```python
# module a.star
def f(x, list=[]):
  list.append(x)
  return list

f(4, [1,2,3])           # [1, 2, 3, 4]
f(1)                    # [1]
f(2)                    # [1, 2], not [2]!

# module b.star
load("a.star", "f")
f(3)                    # error: cannot append to frozen list
```

<b>Variadic functions:</b> Some functions allow callers to provide an
arbitrary number of arguments.
After all required and optional parameters, a function definition may
specify a _variadic arguments_ or _varargs_ parameter, indicated by a
star preceding the parameter name: `*args`.
Any surplus positional arguments provided by the caller are formed
into a tuple and assigned to the `args` parameter.

```python
def f(x, y, *args):
  return x, y, args

f(1, 2)                 # (1, 2, ())
f(1, 2, 3, 4)           # (1, 2, (3, 4))
```

<b>Keyword-variadic functions:</b> Some functions allow callers to
provide an arbitrary sequence of `name=value` keyword arguments.
A function definition may include a final _keyword arguments_ or
_kwargs_ parameter, indicated by a double-star preceding the parameter
name: `**kwargs`.
Any surplus named arguments that do not correspond to named parameters
are collected in a new dictionary and assigned to the `kwargs` parameter:

```python
def f(x, y, **kwargs):
  return x, y, kwargs

f(1, 2)                 # (1, 2, {})
f(x=2, y=1)             # (2, 1, {})
f(x=2, y=1, z=3)        # (2, 1, {"z": 3})
```

It is a static error if any two parameters of a function have the same name.

Just as a function definition may accept an arbitrary number of
positional or named arguments, a function call may provide an
arbitrary number of positional or named arguments supplied by a
list or dictionary:

```python
def f(a, b, c=5):
  return a * b + c

f(*[2, 3])              # 11
f(*[2, 3, 7])           # 13
f(*[2])                 # error: f takes at least 2 arguments (1 given)

f(**dict(b=3, a=2))             # 11
f(**dict(c=7, a=2, b=3))        # 13
f(**dict(a=2))                  # error: f takes at least 2 arguments (1 given)
f(**dict(d=4))                  # error: f got unexpected keyword argument "d"
```

Once the parameters have been successfully bound to the arguments
supplied by the call, the sequence of statements that comprise the
function body is executed.

It is a static error if a function call has two named arguments of the
same name, such as `f(x=1, x=2)`. A call that provides a `**kwargs`
argument may yet have two values for the same name, such as
`f(x=1, **dict(x=2))`. This results in a dynamic error.

Function arguments are evaluated in the order they appear in the call.
<!-- see https://github.com/bazelbuild/starlark/issues/13 -->

Unlike Python, Starlark does not allow more than one `*args` argument in a
call, and if a `*args` argument is present it must appear after all
positional and named arguments.

The final argument to a function call may be followed by a trailing comma.

A function call completes normally after the execution of either a
`return` statement, or of the last statement in the function body.
The result of the function call is the value of the return statement's
operand, or `None` if the return statement had no operand or if the
function completeted without executing a return statement.

```python
def f(x):
  if x == 0:
    return
  if x < 0:
    return -x
  print(x)

f(1)            # returns None after printing "1"
f(0)            # returns None without printing
f(-1)           # returns 1 without printing
```

<b>Implementation note:</b>
The Go implementation of Starlark requires the `-recursion`
flag to allow recursive functions.


If the `-recursion` flag is not specified it is a dynamic error for a
function to call itself or another function value with the same
declaration.

```python
def fib(x):
  if x < 2:
    return x
  return fib(x-2) + fib(x-1)	# dynamic error: function fib called recursively

fib(5)
```

This rule, combined with the invariant that all loops are iterations
over finite sequences, implies that Starlark programs can not be
Turing complete unless the `-recursion` flag is specified.

<!-- This rule is supposed to deter people from abusing Starlark for
     inappropriate uses, especially in the build system.
     It may work for that purpose, but it doesn't stop Starlark programs
     from consuming too much time or space.  Perhaps it should be a
     dialect option.
-->



## Built-in functions

A built-in function is a function or method implemented in Go by the interpreter
or the application into which the interpreter is embedded.

The [type](#type) of a built-in function is `"builtin_function_or_method"`.
<b>Implementation note:</b>
The Java implementation of `type(x)` returns `"function"` for all
functions, whether built in or defined in Starlark,
even though applications distinguish these two types.

A built-in function value used in a Boolean context is always considered true.

Many built-in functions are predeclared in the environment
(see [Name Resolution](#name-resolution)).
Some built-in functions such as `len` are _universal_, that is,
available to all Starlark programs.
The host application may predeclare additional built-in functions
in the environment of a specific module.

Except where noted, built-in functions accept only positional arguments.
The parameter names serve merely as documentation.

Most built-in functions that have a Boolean parameter require its
argument to be `True` or `False`. Unlike `if` statements, other values
are not implicitly converted to their truth value and instead cause a
dynamic error.

