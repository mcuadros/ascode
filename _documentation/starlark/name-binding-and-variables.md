---
title: 'Name binding and variables'
weight: 4
---

After a Starlark file is parsed, but before its execution begins, the
Starlark interpreter checks statically that the program is well formed.
For example, `break` and `continue` statements may appear only within
a loop; a `return` statement may appear only within a
function; and `load` statements may appear only outside any function.

_Name resolution_ is the static checking process that
resolves names to variable bindings.
During execution, names refer to variables.  Statically, names denote
places in the code where variables are created; these places are
called _bindings_.  A name may denote different bindings at different
places in the program.  The region of text in which a particular name
refers to the same binding is called that binding's _scope_.

Four Starlark constructs bind names, as illustrated in the example below:
`load` statements (`a` and `b`),
`def` statements (`c`),
function parameters (`d`),
and assignments (`e`, `h`, including the augmented assignment `e += 1`).
Variables may be assigned or re-assigned explicitly (`e`, `h`), or implicitly, as
in a `for`-loop (`f`) or comprehension (`g`, `i`).

```python
load("lib.star", "a", b="B")

def c(d):
  e = 0
  for f in d:
     print([True for g in f])
     e += 1

h = [2*i for i in a]
```

The environment of a Starlark program is structured as a tree of
_lexical blocks_, each of which may contain name bindings.
The tree of blocks is parallel to the syntax tree.
Blocks are of five kinds.

<!-- Avoid the term "built-in" block since that's also a type. -->
At the root of the tree is the _predeclared_ block,
which binds several names implicitly.
The set of predeclared names includes the universal
constant values `None`, `True`, and `False`, and
various built-in functions such as `len` and `list`;
these functions are immutable and stateless.
An application may pre-declare additional names
to provide domain-specific functions to that file, for example.
These additional functions may have side effects on the application.
Starlark programs cannot change the set of predeclared bindings
or assign new values to them.

Nested beneath the predeclared block is the _module_ block,
which contains the bindings of the current module.
Bindings in the module block (such as `c`, and `h` in the
example) are called _global_ and may be visible to other modules.
The module block is empty at the start of the file
and is populated by top-level binding statements.

Nested beneath the module block is the _file_ block,
which contains bindings local to the current file.
Names in this block (such as `a` and `b` in the example)
are bound only by `load` statements.
The sets of names bound in the file block and in the module block do not overlap:
it is an error for a load statement to bind the name of a global,
or for a top-level statement to assign to a name bound by a load statement.

A file block contains a _function_ block for each top-level
function, and a _comprehension_ block for each top-level comprehension.
Bindings in either of these kinds of block,
and in the file block itself, are called _local_.
(In the example, the bindings for `e`, `f`, `g`, and `i` are all local.)
Additional functions and comprehensions, and their blocks, may be
nested in any order, to any depth.

If name is bound anywhere within a block, all uses of the name within
the block are treated as references to that binding,
even if the use appears before the binding.
This is true even at the top level, unlike Python.
The binding of `y` on the last line of the example below makes `y`
local to the function `hello`, so the use of `y` in the print
statement also refers to the local `y`, even though it appears
earlier.

```python
y = "goodbye"

def hello():
  for x in (1, 2):
    if x == 2:
      print(y) # prints "hello"
    if x == 1:
      y = "hello"
```

It is a dynamic error to evaluate a reference to a local variable
before it has been bound:

```python
def f():
  print(x)              # dynamic error: local variable x referenced before assignment
  x = "hello"
```

The same is true for global variables:

```python
print(x)                # dynamic error: global variable x referenced before assignment
x = "hello"
```

It is a static error to bind a global variable already explicitly bound in the file:

```python
x = 1
x = 2                   # static error: cannot reassign global x declared on line 1
```

<!-- The above rule, and the rule that forbids if-statements and loops at
     top level, exist to ensure that there is exactly one statement
     that binds each global variable, which makes cross-referenced
     documentation more useful, the designers assure me, but
     I am skeptical that it's worth the trouble. -->

If a name was pre-bound by the application, the Starlark program may
explicitly bind it, but only once.

An augmented assignment statement such as `x += y` is considered both a
reference to `x` and a binding use of `x`, so it may not be used at
top level.

<b>Implementation note:</b>
The Go implementation of Starlark permits augmented assignments to appear
at top level if the `-globalreassign` flag is enabled.

A function may refer to variables defined in an enclosing function.
In this example, the inner function `f` refers to a variable `x`
that is local to the outer function `squarer`.
`x` is a _free variable_ of `f`.
The function value (`f`) created by a `def` statement holds a
reference to each of its free variables so it may use
them even after the enclosing function has returned.

```python
def squarer():
    x = [0]
    def f():
      x[0] += 1
      return x[0]*x[0]
    return f

sq = squarer()
print(sq(), sq(), sq(), sq()) # "1 4 9 16"
```

An inner function cannot assign to a variable bound in an enclosing
function, because the assignment would bind the variable in the
inner function.
In the example below, the `x += 1` statement binds `x` within `f`,
hiding the outer `x`.
Execution fails because the inner `x` has not been assigned before the
attempt to increment it.

```python
def squarer():
    x = 0
    def f():
      x += 1            # dynamic error: local variable x referenced before assignment
      return x*x
    return f

sq = squarer()
```

(Starlark has no equivalent of Python's `nonlocal` or `global`
declarations, but as the first version of `squarer` showed, this
omission can be worked around by using a list of a single element.)


A name appearing after a dot, such as `split` in
`get_filename().split('/')`, is not resolved statically.
The [dot expression](/docs/starlark/expressions/#dot-expressions) `.split` is a dynamic operation
on the value returned by `get_filename()`.

