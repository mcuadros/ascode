---
title: 'Expressions'
weight: 6
toc: true
---

## Overview

An expression specifies the computation of a value.

The Starlark grammar defines several categories of expression.
An _operand_ is an expression consisting of a single token (such as an
identifier or a literal), or a bracketed expression.
Operands are self-delimiting.
An operand may be followed by any number of dot, call, or slice
suffixes, to form a _primary_ expression.
In some places in the Starlark grammar where an expression is expected,
it is legal to provide a comma-separated list of expressions denoting
a tuple.
The grammar uses `Expression` where a multiple-component expression is allowed,
and `Test` where it accepts an expression of only a single component.

```grammar {.good}
Expression = Test {',' Test} .

Test = LambdaExpr | IfExpr | PrimaryExpr | UnaryExpr | BinaryExpr .

PrimaryExpr = Operand
            | PrimaryExpr DotSuffix
            | PrimaryExpr CallSuffix
            | PrimaryExpr SliceSuffix
            .

Operand = identifier
        | int | float | string
        | ListExpr | ListComp
        | DictExpr | DictComp
        | '(' [Expression] [,] ')'
        | ('-' | '+') PrimaryExpr
        .

DotSuffix   = '.' identifier .
CallSuffix  = '(' [Arguments [',']] ')' .
SliceSuffix = '[' [Expression] [':' Test [':' Test]] ']' .
```

TODO: resolve position of +x, -x, and 'not x' in grammar: Operand or UnaryExpr?

## Identifiers

```grammar {.good} {.good}
Primary = identifier
```

An identifier is a name that identifies a value.

Lookup of locals and globals may fail if not yet defined.

## Literals

Starlark supports literals of three different kinds:

```grammar {.good}
Primary = int | float | string
```

Evaluation of a literal yields a value of the given type (string, int,
or float) with the given value.
See [Literals](#lexical-elements) for details.

## Parenthesized expressions

```grammar {.good}
Primary = '(' [Expression] ')'
```

A single expression enclosed in parentheses yields the result of that expression.
Explicit parentheses may be used for clarity,
or to override the default association of subexpressions.

```python
1 + 2 * 3 + 4                   # 11
(1 + 2) * (3 + 4)               # 21
```

If the parentheses are empty, or contain a single expression followed
by a comma, or contain two or more expressions, the expression yields a tuple.

```python
()                              # (), the empty tuple
(1,)                            # (1,), a tuple of length 1
(1, 2)                          # (1, 2), a 2-tuple or pair
(1, 2, 3)                       # (1, 2, 3), a 3-tuple or triple
```

In some contexts, such as a `return` or assignment statement or the
operand of a `for` statement, a tuple may be expressed without
parentheses.

```python
x, y = 1, 2

return 1, 2

for x in 1, 2:
   print(x)
```

Starlark (like Python 3) does not accept an unparenthesized tuple
expression as the operand of a list comprehension:

```python
[2*x for x in 1, 2, 3]	       	# parse error: unexpected ','
```

## Dictionary expressions

A dictionary expression is a comma-separated list of colon-separated
key/value expression pairs, enclosed in curly brackets, and it yields
a new dictionary object.
An optional comma may follow the final pair.

```grammar {.good}
DictExpr = '{' [Entries [',']] '}' .
Entries  = Entry {',' Entry} .
Entry    = Test ':' Test .
```

Examples:


```python
{}
{"one": 1}
{"one": 1, "two": 2,}
```

The key and value expressions are evaluated in left-to-right order.
Evaluation fails if the same key is used multiple times.

Only [hashable](#hashing) values may be used as the keys of a dictionary.
This includes all built-in types except dictionaries, sets, and lists;
a tuple is hashable only if its elements are hashable.


## List expressions

A list expression is a comma-separated list of element expressions,
enclosed in square brackets, and it yields a new list object.
An optional comma may follow the last element expression.

```grammar {.good}
ListExpr = '[' [Expression [',']] ']' .
```

Element expressions are evaluated in left-to-right order.

Examples:

```python
[]                      # [], empty list
[1]                     # [1], a 1-element list
[1, 2, 3,]              # [1, 2, 3], a 3-element list
```

## Unary operators

There are three unary operators, all appearing before their operand:
`+`, `-`, `~`, and `not`.

```grammar {.good}
UnaryExpr = '+' PrimaryExpr
          | '-' PrimaryExpr
          | '~' PrimaryExpr
          | 'not' Test
          .
```

```text
+ number        unary positive          (int, float)
- number        unary negation          (int, float)
~ number        unary bitwise inversion (int)
not x           logical negation        (any type)
```

The `+` and `-` operators may be applied to any number
(`int` or `float`) and return the number unchanged.
Unary `+` is never necessary in a correct program,
but may serve as an assertion that its operand is a number,
or as documentation.

```python
if x > 0:
	return +1
else if x < 0:
	return -1
else:
	return 0
```

The `not` operator returns the negation of the truth value of its
operand.

```python
not True                        # False
not False                       # True
not [1, 2, 3]                   # False
not ""                          # True
not 0                           # True
```

The `~` operator yields the bitwise inversion of its integer argument.
The bitwise inversion of x is defined as -(x+1).

```python
~1                              # -2
~-1                             # 0
~0                              # -1
```


## Binary operators

Starlark has the following binary operators, arranged in order of increasing precedence:

```text
or
and
==   !=   <    >   <=   >=   in   not in
|
^
&
<<   >>
-    +
*    /    //   %
```

Comparison operators, `in`, and `not in` are non-associative,
so the parser will not accept `0 <= i < n`.
All other binary operators of equal precedence associate to the left.

```grammar {.good}
BinaryExpr = Test {Binop Test} .

Binop = 'or'
      | 'and'
      | '==' | '!=' | '<' | '>' | '<=' | '>=' | 'in' | 'not' 'in'
      | '|'
      | '^'
      | '&'
      | '-' | '+'
      | '*' | '%' | '/' | '//'
      | '<<' | '>>'
      .
```

### `or` and `and`

The `or` and `and` operators yield, respectively, the logical disjunction and
conjunction of their arguments, which need not be Booleans.
The expression `x or y` yields the value of `x` if its truth value is `True`,
or the value of `y` otherwise.

```starlark
False or False		# False
False or True		# True
True  or False		# True
True  or True		# True

0 or "hello"		# "hello"
1 or "hello"		# 1
```

Similarly, `x and y` yields the value of `x` if its truth value is
`False`, or the value of `y` otherwise.

```starlark
False and False		# False
False and True		# False
True  and False		# False
True  and True		# True

0 and "hello"		# 0
1 and "hello"		# "hello"
```

These operators use "short circuit" evaluation, so the second
expression is not evaluated if the value of the first expression has
already determined the result, allowing constructions like these:

```python
len(x) > 0 and x[0] == 1		# x[0] is not evaluated if x is empty
x and x[0] == 1
len(x) == 0 or x[0] == ""
not x or not x[0]
```

### Comparisons

The `==` operator reports whether its operands are equal; the `!=`
operator is its negation.

The operators `<`, `>`, `<=`, and `>=` perform an ordered comparison
of their operands.  It is an error to apply these operators to
operands of unequal type, unless one of the operands is an `int` and
the other is a `float`.  Of the built-in types, only the following
support ordered comparison, using the ordering relation shown:

```shell
NoneType        # None <= None
bool            # False < True
int             # mathematical
float           # as defined by IEEE 754
string          # lexicographical
tuple           # lexicographical
list            # lexicographical
```

Comparison of floating point values follows the IEEE 754 standard,
which breaks several mathematical identities.  For example, if `x` is
a `NaN` value, the comparisons `x < y`, `x == y`, and `x > y` all
yield false for all values of `y`.

Applications may define additional types that support ordered
comparison.

The remaining built-in types support only equality comparisons.
Values of type `dict` or `set` compare equal if their elements compare
equal, and values of type `function` or `builtin_function_or_method` are equal only to
themselves.

```shell
dict                            # equal contents
set                             # equal contents
function                        # identity
builtin_function_or_method      # identity
```

### Arithmetic operations

The following table summarizes the binary arithmetic operations
available for built-in types:

```shell
Arithmetic (int or float; result has type float unless both operands have type int)
   number + number              # addition
   number - number              # subtraction
   number * number              # multiplication
   number / number              # real division  (result is always a float)
   number // number             # floored division
   number % number              # remainder of floored division
   number ^ number              # bitwise XOR
   number << number             # bitwise left shift
   number >> number             # bitwise right shift

Concatenation
   string + string
     list + list
    tuple + tuple

Repetition (string/list/tuple)
      int * sequence
 sequence * int

String interpolation
   string % any                 # see String Interpolation

Sets
      int | int                 # bitwise union (OR)
      set | set                 # set union
      int & int                 # bitwise intersection (AND)
      set & set                 # set intersection
      set ^ set                 # set symmetric difference
```

The operands of the arithmetic operators `+`, `-`, `*`, `//`, and
`%` must both be numbers (`int` or `float`) but need not have the same type.
The type of the result has type `int` only if both operands have that type.
The result of real division `/` always has type `float`.

The `+` operator may be applied to non-numeric operands of the same
type, such as two lists, two tuples, or two strings, in which case it
computes the concatenation of the two operands and yields a new value of
the same type.

```python
"Hello, " + "world"		# "Hello, world"
(1, 2) + (3, 4)			# (1, 2, 3, 4)
[1, 2] + [3, 4]			# [1, 2, 3, 4]
```

The `*` operator may be applied to an integer _n_ and a value of type
`string`, `list`, or `tuple`, in which case it yields a new value
of the same sequence type consisting of _n_ repetitions of the original sequence.
The order of the operands is immaterial.
Negative values of _n_ behave like zero.

```python
'mur' * 2               # 'murmur'
3 * range(3)            # [0, 1, 2, 0, 1, 2, 0, 1, 2]
```

Applications may define additional types that support any subset of
these operators.

The `&` operator requires two operands of the same type, either `int` or `set`.
For integers, it yields the bitwise intersection (AND) of its operands.
For sets, it yields a new set containing the intersection of the
elements of the operand sets, preserving the element order of the left
operand.

The `|` operator likewise computes bitwise or set unions.
The result of `set | set` is a new set whose elements are the
union of the operands, preserving the order of the elements of the
operands, left before right.

The `^` operator accepts operands of either `int` or `set` type.
For integers, it yields the bitwise XOR (exclusive OR) of its operands.
For sets, it yields a new set containing elements of either first or second
operand but not both (symmetric difference).

The `<<` and `>>` operators require operands of `int` type both. They shift
the first operand to the left or right by the number of bits given by the
second operand. It is a dynamic error if the second operand is negative.
Implementations may impose a limit on the second operand of a left shift.

```python
0x12345678 & 0xFF               # 0x00000078
0x12345678 | 0xFF               # 0x123456FF
0b01011101 ^ 0b110101101        # 0b111110000
0b01011101 >> 2                 # 0b010111
0b01011101 << 2                 # 0b0101110100

set([1, 2]) & set([2, 3])       # set([2])
set([1, 2]) | set([2, 3])       # set([1, 2, 3])
set([1, 2]) ^ set([2, 3])       # set([1, 3])
```

<b>Implementation note:</b>
The Go implementation of Starlark requires the `-set` flag to
enable support for sets.
The Java implementation does not support sets.


### Membership tests

```text
      any in     sequence		(list, tuple, dict, set, string)
      any not in sequence
```

The `in` operator reports whether its first operand is a member of its
second operand, which must be a list, tuple, dict, set, or string.
The `not in` operator is its negation.
Both return a Boolean.

The meaning of membership varies by the type of the second operand:
the members of a list, tuple, or set are its elements;
the members of a dict are its keys;
the members of a string are all its substrings.

```python
1 in [1, 2, 3]                  # True
4 in (1, 2, 3)                  # False
4 not in set([1, 2, 3])         # True

d = {"one": 1, "two": 2}
"one" in d                      # True
"three" in d                    # False
1 in d                          # False
[] in d				# False

"nasty" in "dynasty"            # True
"a" in "banana"                 # True
"f" not in "way"                # True
```

### String interpolation

The expression `format % args` performs _string interpolation_, a
simple form of template expansion.
The `format` string is interpreted as a sequence of literal portions
and _conversions_.
Each conversion, which starts with a `%` character, is replaced by its
corresponding value from `args`.
The characters following `%` in each conversion determine which
argument it uses and how to convert it to a string.

Each `%` character marks the start of a conversion specifier, unless
it is immediately followed by another `%`, in which case both
characters together denote a literal percent sign.

If the `"%"` is immediately followed by `"(key)"`, the parenthesized
substring specifies the key of the `args` dictionary whose
corresponding value is the operand to convert.
Otherwise, the conversion's operand is the next element of `args`,
which must be a tuple with exactly one component per conversion,
unless the format string contains only a single conversion, in which
case `args` itself is its operand.

Starlark does not support the flag, width, and padding specifiers
supported by Python's `%` and other variants of C's `printf`.

After the optional `(key)` comes a single letter indicating what
operand types are valid and how to convert the operand `x` to a string:

```text
%       none            literal percent sign
s       any             as if by str(x)
r       any             as if by repr(x)
d       number          signed integer decimal
i       number          signed integer decimal
o       number          signed octal
x       number          signed hexadecimal, lowercase
X       number          signed hexadecimal, uppercase
e       number          float exponential format, lowercase
E       number          float exponential format, uppercase
f       number          float decimal format, lowercase
F       number          float decimal format, uppercase
g       number          like %e for large exponents, %f otherwise
G       number          like %E for large exponents, %F otherwise
c       string          x (string must encode a single Unicode code point)
        int             as if by chr(x)
```

It is an error if the argument does not have the type required by the
conversion specifier.  A Boolean argument is not considered a number.

Examples:

```python
"Hello %s, your score is %d" % ("Bob", 75)      # "Hello Bob, your score is 75"

"%d %o %x %c" % (65, 65, 65, 65)                # "65 101 41 A" (decimal, octal, hexadecimal, Unicode)

"%(greeting)s, %(audience)s" % dict(            # "Hello, world"
  greeting="Hello",
  audience="world",
)

"rate = %g%% APR" % 3.5                         # "rate = 3.5% APR"
```

One subtlety: to use a tuple as the operand of a conversion in format
string containing only a single conversion, you must wrap the tuple in
a singleton tuple:

```python
"coordinates=%s" % (40.741491, -74.003680)	# error: too many arguments for format string
"coordinates=%s" % ((40.741491, -74.003680),)	# "coordinates=(40.741491, -74.003680)"
```

TODO: specify `%e` and `%f` more precisely.

## Conditional expressions

A conditional expression has the form `a if cond else b`.
It first evaluates the condition `cond`.
If it's true, it evaluates `a` and yields its value;
otherwise it yields the value of `b`.

```grammar {.good}
IfExpr = Test 'if' Test 'else' Test .
```

Example:

```python
"yes" if enabled else "no"
```

## Comprehensions

A comprehension constructs new list or dictionary value by looping
over one or more iterables and evaluating a _body_ expression that produces
successive elements of the result.

A list comprehension consists of a single expression followed by one
or more _clauses_, the first of which must be a `for` clause.
Each `for` clause resembles a `for` statement, and specifies an
iterable operand and a set of variables to be assigned by successive
values of the iterable.
An `if` cause resembles an `if` statement, and specifies a condition
that must be met for the body expression to be evaluated.
A sequence of `for` and `if` clauses acts like a nested sequence of
`for` and `if` statements.

```grammar {.good}
ListComp = '[' Test {CompClause} ']'.
DictComp = '{' Entry {CompClause} '}' .

CompClause = 'for' LoopVariables 'in' Test
           | 'if' Test .

LoopVariables = PrimaryExpr {',' PrimaryExpr} .
```

Examples:

```python
[x*x for x in range(5)]                 # [0, 1, 4, 9, 16]
[x*x for x in range(5) if x%2 == 0]     # [0, 4, 16]
[(x, y) for x in range(5)
        if x%2 == 0
        for y in range(5)
        if y > x]                       # [(0, 1), (0, 2), (0, 3), (0, 4), (2, 3), (2, 4)]
```

A dict comprehension resembles a list comprehension, but its body is a
pair of expressions, `key: value`, separated by a colon,
and its result is a dictionary containing the key/value pairs
for which the body expression was evaluated.
Evaluation fails if the value of any key is unhashable.

As with a `for` loop, the loop variables may exploit compound
assignment:

```python
[x*y+z for (x, y), z in [((2, 3), 5), (("o", 2), "!")]]         # [11, 'oo!']
```

Starlark, following Python 3, does not accept an unparenthesized
tuple or lambda expression as the operand of a `for` clause:

```python
[x*x for x in 1, 2, 3]		# parse error: unexpected comma
[x*x for x in lambda: 0]	# parse error: unexpected lambda
```

Comprehensions in Starlark, again following Python 3, define a new lexical
block, so assignments to loop variables have no effect on variables of
the same name in an enclosing block:

```python
x = 1
_ = [x for x in [2]]            # new variable x is local to the comprehension
print(x)                        # 1
```

The operand of a comprehension's first clause (always a `for`) is
resolved in the lexical block enclosing the comprehension.
In the examples below, identifiers referring to the outer variable
named `x` have been distinguished by subscript.

```python
x₀ = (1, 2, 3)
[x*x for x in x₀]               # [1, 4, 9]
[x*x for x in x₀ if x%2 == 0]   # [4]
```

All subsequent `for` and `if` expressions are resolved within the
comprehension's lexical block, as in this rather obscure example:

```python
x₀ = ([1, 2], [3, 4], [5, 6])
[x*x for x in x₀ for x in x if x%2 == 0]     # [4, 16, 36]
```

which would be more clearly rewritten as:

```python
x = ([1, 2], [3, 4], [5, 6])
[z*z for y in x for z in y if z%2 == 0]     # [4, 16, 36]
```


## Function and method calls

```grammar {.good}
CallSuffix = '(' [Arguments [',']] ')' .

Arguments = Argument {',' Argument} .
Argument  = Test | identifier '=' Test | '*' Test | '**' Test .
```

A value `f` of type `function` or `builtin_function_or_method` may be called using the expression `f(...)`.
Applications may define additional types whose values may be called in the same way.

A method call such as `filename.endswith(".star")` is the composition
of two operations, `m = filename.endswith` and `m(".star")`.
The first, a dot operation, yields a _bound method_, a function value
that pairs a receiver value (the `filename` string) with a choice of
method ([string·endswith](#string·endswith)).

Only built-in or application-defined types may have methods.

See [Functions](#functions) for an explanation of function parameter passing.

## Dot expressions

A dot expression `x.f` selects the attribute `f` (a field or method)
of the value `x`.

Fields are possessed by none of the main Starlark [data types](#data-types),
but some application-defined types have them.
Methods belong to the built-in types `string`, `list`, `dict`, and
`set`, and to many application-defined types.

```grammar {.good}
DotSuffix = '.' identifier .
```

A dot expression fails if the value does not have an attribute of the
specified name.

Use the built-in function `hasattr(x, "f")` to ascertain whether a
value has a specific attribute, or `dir(x)` to enumerate all its
attributes.  The `getattr(x, "f")` function can be used to select an
attribute when the name `"f"` is not known statically.

A dot expression that selects a method typically appears within a call
expression, as in these examples:

```python
["able", "baker", "charlie"].index("baker")     # 1
"banana".count("a")                             # 3
"banana".reverse()                              # error: string has no .reverse field or method
```

But when not called immediately, the dot expression evaluates to a
_bound method_, that is, a method coupled to a specific receiver
value.  A bound method can be called like an ordinary function,
without a receiver argument:

```python
f = "banana".count
f                                               # <built-in method count of string value>
f("a")                                          # 3
f("n")                                          # 2
```

## Index expressions

An index expression `a[i]` yields the `i`th element of an _indexable_
type such as a string, tuple, or list.  The index `i` must be an `int`
value in the range -`n` ≤ `i` < `n`, where `n` is `len(a)`; any other
index results in an error.

```grammar {.good}
SliceSuffix = '[' [Expression] [':' Test [':' Test]] ']' .
```

A valid negative index `i` behaves like the non-negative index `n+i`,
allowing for convenient indexing relative to the end of the
sequence.

```python
"abc"[0]                        # "a"
"abc"[1]                        # "b"
"abc"[-1]                       # "c"

("zero", "one", "two")[0]       # "zero"
("zero", "one", "two")[1]       # "one"
("zero", "one", "two")[-1]      # "two"
```

An index expression `d[key]` may also be applied to a dictionary `d`,
to obtain the value associated with the specified key.  It is an error
if the dictionary contains no such key.

An index expression appearing on the left side of an assignment causes
the specified list or dictionary element to be updated:

```starlark
a = range(3)            # a == [0, 1, 2]
a[2] = 7                # a == [0, 1, 7]

coins["suzie b"] = 100
```

It is a dynamic error to attempt to update an element of an immutable
type, such as a tuple or string, or a frozen value of a mutable type.

## Slice expressions

A slice expression `a[start:stop:stride]` yields a new value containing a
sub-sequence of `a`, which must be a string, tuple, or list.

```grammar {.good}
SliceSuffix = '[' [Expression] [':' Test [':' Test]] ']' .
```

Each of the `start`, `stop`, and `stride` operands is optional;
if present, and not `None`, each must be an integer.
The `stride` value defaults to 1.
If the stride is not specified, the colon preceding it may be omitted too.
It is an error to specify a stride of zero.

Conceptually, these operands specify a sequence of values `i` starting
at `start` and successively adding `stride` until `i` reaches or
passes `stop`. The result consists of the concatenation of values of
`a[i]` for which `i` is valid.`

The effective start and stop indices are computed from the three
operands as follows.  Let `n` be the length of the sequence.

<b>If the stride is positive:</b>
If the `start` operand was omitted, it defaults to -infinity.
If the `end` operand was omitted, it defaults to +infinity.
For either operand, if a negative value was supplied, `n` is added to it.
The `start` and `end` values are then "clamped" to the
nearest value in the range 0 to `n`, inclusive.

<b>If the stride is negative:</b>
If the `start` operand was omitted, it defaults to +infinity.
If the `end` operand was omitted, it defaults to -infinity.
For either operand, if a negative value was supplied, `n` is added to it.
The `start` and `end` values are then "clamped" to the
nearest value in the range -1 to `n`-1, inclusive.

```python
"abc"[1:]               # "bc"  (remove first element)
"abc"[:-1]              # "ab"  (remove last element)
"abc"[1:-1]             # "b"   (remove first and last element)
"banana"[1::2]          # "aaa" (select alternate elements starting at index 1)
"banana"[4::-2]         # "nnb" (select alternate elements in reverse, starting at index 4)
```

Unlike Python, Starlark does not allow a slice expression on the left
side of an assignment.

Slicing a tuple or string may be more efficient than slicing a list
because tuples and strings are immutable, so the result of the
operation can share the underlying representation of the original
operand (when the stride is 1). By contrast, slicing a list requires
the creation of a new list and copying of the necessary elements.

<!-- TODO tighten up this section -->

## Lambda expressions

A `lambda` expression yields a new function value.

```grammar {.good}
LambdaExpr = 'lambda' [Parameters] ':' Test .

Parameters = Parameter {',' Parameter} .
Parameter  = identifier
           | identifier '=' Test
           | '*'
           | '*' identifier
           | '**' identifier
           .
```

Syntactically, a lambda expression consists of the keyword `lambda`,
followed by a parameter list like that of a `def` statement but
unparenthesized, then a colon `:`, and a single expression, the
_function body_.

Example:

```python
def map(f, list):
    return [f(x) for x in list]

map(lambda x: 2*x, range(3))    # [2, 4, 6]
```

As with functions created by a `def` statement, a lambda function
captures the syntax of its body, the default values of any optional
parameters, the value of each free variable appearing in its body, and
the global dictionary of the current module.

The name of a function created by a lambda expression is `"lambda"`.

The two statements below are essentially equivalent, but the
function created by the `def` statement is named `twice` and the
function created by the lambda expression is named `lambda`.

```python
def twice(x):
   return x * 2

twice = lambda x: x * 2
```

<b>Implementation note:</b>
The Go implementation of Starlark requires the `-lambda` flag
to enable support for lambda expressions.
The Java implementation does not support them.
See Google Issue b/36358844.

