---
title: 'Built-in methods'
weight: 10
toc: true
---

## Overview

This section lists the methods of built-in types.  Methods are selected
using [dot expressions](#dot-expressions).
For example, strings have a `count` method that counts
occurrences of a substring; `"banana".count("a")` yields `3`.

As with built-in functions, built-in methods accept only positional
arguments except where noted.
The parameter names serve merely as documentation.


## dict·clear

`D.clear()` removes all the entries of dictionary D and returns `None`.
It fails if the dictionary is frozen or if there are active iterators.

```python
x = {"one": 1, "two": 2}
x.clear()                               # None
print(x)                                # {}
```

## dict·get

`D.get(key[, default])` returns the dictionary value corresponding to the given key.
If the dictionary contains no such value, `get` returns `None`, or the
value of the optional `default` parameter if present.

`get` fails if `key` is unhashable, or the dictionary is frozen or has active iterators.

```python
x = {"one": 1, "two": 2}
x.get("one")                            # 1
x.get("three")                          # None
x.get("three", 0)                       # 0
```

## dict·items

`D.items()` returns a new list of key/value pairs, one per element in
dictionary D, in the same order as they would be returned by a `for` loop.

```python
x = {"one": 1, "two": 2}
x.items()                               # [("one", 1), ("two", 2)]
```

## dict·keys

`D.keys()` returns a new list containing the keys of dictionary D, in the
same order as they would be returned by a `for` loop.

```python
x = {"one": 1, "two": 2}
x.keys()                               # ["one", "two"]
```

## dict·pop

`D.pop(key[, default])` returns the value corresponding to the specified
key, and removes it from the dictionary.  If the dictionary contains no
such value, and the optional `default` parameter is present, `pop`
returns that value; otherwise, it fails.

`pop` fails if `key` is unhashable, or the dictionary is frozen or has active iterators.

```python
x = {"one": 1, "two": 2}
x.pop("one")                            # 1
x                                       # {"two": 2}
x.pop("three", 0)                       # 0
x.pop("four")                           # error: missing key
```

## dict·popitem

`D.popitem()` returns the first key/value pair, removing it from the dictionary.

`popitem` fails if the dictionary is empty, frozen, or has active iterators.

```python
x = {"one": 1, "two": 2}
x.popitem()                             # ("one", 1)
x.popitem()                             # ("two", 2)
x.popitem()                             # error: empty dict
```

## dict·setdefault

`D.setdefault(key[, default])` returns the dictionary value corresponding to the given key.
If the dictionary contains no such value, `setdefault`, like `get`,
returns `None` or the value of the optional `default` parameter if
present; `setdefault` additionally inserts the new key/value entry into the dictionary.

`setdefault` fails if the key is unhashable, or if the dictionary is frozen or has active iterators.

```python
x = {"one": 1, "two": 2}
x.setdefault("one")                     # 1
x.setdefault("three", 0)                # 0
x                                       # {"one": 1, "two": 2, "three": 0}
x.setdefault("four")                    # None
x                                       # {"one": 1, "two": 2, "three": None}
```

## dict·update

`D.update([pairs][, name=value[, ...])` makes a sequence of key/value
insertions into dictionary D, then returns `None.`

If the positional argument `pairs` is present, it must be `None`,
another `dict`, or some other iterable.
If it is another `dict`, then its key/value pairs are inserted into D.
If it is an iterable, it must provide a sequence of pairs (or other iterables of length 2),
each of which is treated as a key/value pair to be inserted into D.

For each `name=value` argument present, the name is converted to a
string and used as the key for an insertion into D, with its corresponding
value being `value`.

`update` fails if the dictionary is frozen or has active iterators.

```python
x = {}
x.update([("a", 1), ("b", 2)], c=3)
x.update({"d": 4})
x.update(e=5)
x                                       # {"a": 1, "b": "2", "c": 3, "d": 4, "e": 5}
```

## dict·values

`D.values()` returns a new list containing the dictionary's values, in the
same order as they would be returned by a `for` loop over the
dictionary.

```python
x = {"one": 1, "two": 2}
x.values()                              # [1, 2]
```

## list·append

`L.append(x)` appends `x` to the list L, and returns `None`.

`append` fails if the list is frozen or has active iterators.

```python
x = []
x.append(1)                             # None
x.append(2)                             # None
x.append(3)                             # None
x                                       # [1, 2, 3]
```

## list·clear

`L.clear()` removes all the elements of the list L and returns `None`.
It fails if the list is frozen or if there are active iterators.

```python
x = [1, 2, 3]
x.clear()                               # None
x                                       # []
```

## list·extend

`L.extend(x)` appends the elements of `x`, which must be iterable, to
the list L, and returns `None`.

`extend` fails if `x` is not iterable, or if the list L is frozen or has active iterators.

```python
x = []
x.extend([1, 2, 3])                     # None
x.extend(["foo"])                       # None
x                                       # [1, 2, 3, "foo"]
```

## list·index

`L.index(x[, start[, end]])` finds `x` within the list L and returns its index.

The optional `start` and `end` parameters restrict the portion of
list L that is inspected.  If provided and not `None`, they must be list
indices of type `int`. If an index is negative, `len(L)` is effectively
added to it, then if the index is outside the range `[0:len(L)]`, the
nearest value within that range is used; see [Indexing](#indexing).

`index` fails if `x` is not found in L, or if `start` or `end`
is not a valid index (`int` or `None`).

```python
x = list("banana".codepoints())
x.index("a")                            # 1 (bAnana)
x.index("a", 2)                         # 3 (banAna)
x.index("a", -2)                        # 5 (bananA)
```

## list·insert

`L.insert(i, x)` inserts the value `x` in the list L at index `i`, moving
higher-numbered elements along by one.  It returns `None`.

As usual, the index `i` must be an `int`. If its value is negative,
the length of the list is added, then its value is clamped to the
nearest value in the range `[0:len(L)]` to yield the effective index.

`insert` fails if the list is frozen or has active iterators.

```python
x = ["b", "c", "e"]
x.insert(0, "a")                        # None
x.insert(-1, "d")                       # None
x                                       # ["a", "b", "c", "d", "e"]
```

## list·pop

`L.pop([index])` removes and returns the last element of the list L, or,
if the optional index is provided, at that index.

`pop` fails if the index is not valid for `L[i]`,
or if the list is frozen or has active iterators.

```python
x = [1, 2, 3, 4, 5]
x.pop()                                 # 5
x                                       # [1, 2, 3, 4]
x.pop(-2)                               # 3
x                                       # [1, 2, 4]
x.pop(-3)                               # 1
x                                       # [2, 4]
x.pop()                                 # 4
x                                       # [2]
```

## list·remove

`L.remove(x)` removes the first occurrence of the value `x` from the list L, and returns `None`.

`remove` fails if the list does not contain `x`, is frozen, or has active iterators.

```python
x = [1, 2, 3, 2]
x.remove(2)                             # None (x == [1, 3, 2])
x.remove(2)                             # None (x == [1, 3])
x.remove(2)                             # error: element not found
```

## set·union

`S.union(iterable)` returns a new set into which have been inserted
all the elements of set S and all the elements of the argument, which
must be iterable.

`union` fails if any element of the iterable is not hashable.

```python
x = set([1, 2])
y = set([2, 3])
x.union(y)                              # set([1, 2, 3])
```

## string·elem_ords

`S.elem_ords()` returns an iterable value containing the
sequence of numeric bytes values in the string S.

To materialize the entire sequence of bytes, apply `list(...)` to the result.

Example:

```python
list("Hello, 世界".elem_ords())        # [72, 101, 108, 108, 111, 44, 32, 228, 184, 150, 231, 149, 140]
```

See also: `string·elems`.

<b>Implementation note:</b> `elem_ords` is not provided by the Java implementation.

## string·capitalize

`S.capitalize()` returns a copy of string S with its first code point
changed to its title case and all subsequent letters changed to their
lower case.

```python
"hello, world!".capitalize()		# "Hello, world!"
"hElLo, wOrLd!".capitalize()		# "Hello, world!"
"¿Por qué?".capitalize()		# "¿por qué?"
```

## string·codepoint_ords

`S.codepoint_ords()` returns an iterable value containing the
sequence of integer Unicode code points encoded by the string S.
Each invalid code within the string is treated as if it encodes the
Unicode replacement character, U+FFFD.

By returning an iterable, not a list, the cost of decoding the string
is deferred until actually needed; apply `list(...)` to the result to
materialize the entire sequence.

Example:

```python
list("Hello, 世界".codepoint_ords())        # [72, 101, 108, 108, 111, 44, 32, 19990, 30028]

for cp in "Hello, 世界".codepoint_ords():
   print(chr(cp))  # prints 'H', 'e', 'l', 'l', 'o', ',', ' ', '世', '界'
```

See also: `string·codepoints`.

<b>Implementation note:</b> `codepoint_ords` is not provided by the Java implementation.

## string·count

`S.count(sub[, start[, end]])` returns the number of occcurences of
`sub` within the string S, or, if the optional substring indices
`start` and `end` are provided, within the designated substring of S.
They are interpreted according to Starlark's [indexing conventions](#indexing).

```python
"hello, world!".count("o")              # 2
"hello, world!".count("o", 7, 12)       # 1  (in "world")
```

## string·endswith

`S.endswith(suffix[, start[, end]])` reports whether the string
`S[start:end]` has the specified suffix.

```python
"filename.star".endswith(".star")         # True
```

The `suffix` argument may be a tuple of strings, in which case the
function reports whether any one of them is a suffix.

```python
'foo.cc'.endswith(('.cc', '.h'))         # True
```


## string·find

`S.find(sub[, start[, end]])` returns the index of the first
occurrence of the substring `sub` within S.

If either or both of `start` or `end` are specified,
they specify a subrange of S to which the search should be restricted.
They are interpreted according to Starlark's [indexing conventions](#indexing).

If no occurrence is found, `found` returns -1.

```python
"bonbon".find("on")             # 1
"bonbon".find("on", 2)          # 4
"bonbon".find("on", 2, 5)       # -1
```

## string·format

`S.format(*args, **kwargs)` returns a version of the format string S
in which bracketed portions `{...}` are replaced
by arguments from `args` and `kwargs`.

Within the format string, a pair of braces `{{` or `}}` is treated as
a literal open or close brace.
Each unpaired open brace must be matched by a close brace `}`.
The optional text between corresponding open and close braces
specifies which argument to use and how to format it, and consists of
three components, all optional:
a field name, a conversion preceded by '`!`', and a format specifier
preceded by '`:`'.

```text
{field}
{field:spec}
{field!conv}
{field!conv:spec}
```

The *field name* may be either a decimal number or a keyword.
A number is interpreted as the index of a positional argument;
a keyword specifies the value of a keyword argument.
If all the numeric field names form the sequence 0, 1, 2, and so on,
they may be omitted and those values will be implied; however,
the explicit and implicit forms may not be mixed.

The *conversion* specifies how to convert an argument value `x` to a
string. It may be either `!r`, which converts the value using
`repr(x)`, or `!s`, which converts the value using `str(x)` and is
the default.

The *format specifier*, after a colon, specifies field width,
alignment, padding, and numeric precision.
Currently it must be empty, but it is reserved for future use.

```python
"a{x}b{y}c{}".format(1, x=2, y=3)               # "a2b3c1"
"a{}b{}c".format(1, 2)                          # "a1b2c"
"({1}, {0})".format("zero", "one")              # "(one, zero)"
"Is {0!r} {0!s}?".format('heterological')       # 'is "heterological" heterological?'
```

## string·index

`S.index(sub[, start[, end]])` returns the index of the first
occurrence of the substring `sub` within S, like `S.find`, except
that if the substring is not found, the operation fails.

```python
"bonbon".index("on")             # 1
"bonbon".index("on", 2)          # 4
"bonbon".index("on", 2, 5)       # error: substring not found  (in "nbo")
```

## string·isalnum

`S.isalnum()` reports whether the string S is non-empty and consists only
Unicode letters and digits.

```python
"base64".isalnum()              # True
"Catch-22".isalnum()            # False
```

## string·isalpha

`S.isalpha()` reports whether the string S is non-empty and consists only of Unicode letters.

```python
"ABC".isalpha()                 # True
"Catch-22".isalpha()            # False
"".isalpha()                    # False
```

## string·isdigit

`S.isdigit()` reports whether the string S is non-empty and consists only of Unicode digits.

```python
"123".isdigit()                 # True
"Catch-22".isdigit()            # False
"".isdigit()                    # False
```

## string·islower

`S.islower()` reports whether the string S contains at least one cased Unicode
letter, and all such letters are lowercase.

```python
"hello, world".islower()        # True
"Catch-22".islower()            # False
"123".islower()                 # False
```

## string·isspace

`S.isspace()` reports whether the string S is non-empty and consists only of Unicode spaces.

```python
"    ".isspace()                # True
"\r\t\n".isspace()              # True
"".isspace()                    # False
```

## string·istitle

`S.istitle()` reports whether the string S contains at least one cased Unicode
letter, and all such letters that begin a word are in title case.

```python
"Hello, World!".istitle()       # True
"Catch-22".istitle()            # True
"HAL-9000".istitle()            # False
"ǅenan".istitle()		# True
"Ǆenan".istitle()		# False ("Ǆ" is a single Unicode letter)
"123".istitle()                 # False
```

## string·isupper

`S.isupper()` reports whether the string S contains at least one cased Unicode
letter, and all such letters are uppercase.

```python
"HAL-9000".isupper()            # True
"Catch-22".isupper()            # False
"123".isupper()                 # False
```

## string·join

`S.join(iterable)` returns the string formed by concatenating each
element of its argument, with a copy of the string S between
successive elements. The argument must be an iterable whose elements
are strings.

```python
", ".join(["one", "two", "three"])      # "one, two, three"
"a".join("ctmrn".codepoints())          # "catamaran"
```

## string·lower

`S.lower()` returns a copy of the string S with letters converted to lowercase.

```python
"Hello, World!".lower()                 # "hello, world!"
```

## string·lstrip

`S.lstrip()` returns a copy of the string S with leading whitespace removed.

Like `strip`, it accepts an optional string parameter that specifies an
alternative set of Unicode code points to remove.

```python
"  hello  ".lstrip()                    # "hello  "
"  hello  ".lstrip("h o")               # "ello  "
```

## string·partition

`S.partition(x)` splits string S into three parts and returns them as
a tuple: the portion before the first occurrence of string `x`, `x` itself,
and the portion following it.
If S does not contain `x`, `partition` returns `(S, "", "")`.

`partition` fails if `x` is not a string, or is the empty string.

```python
"one/two/three".partition("/")		# ("one", "/", "two/three")
```

## string·replace

`S.replace(old, new[, count])` returns a copy of string S with all
occurrences of substring `old` replaced by `new`. If the optional
argument `count`, which must be an `int`, is non-negative, it
specifies a maximum number of occurrences to replace.

```python
"banana".replace("a", "o")		# "bonono"
"banana".replace("a", "o", 2)		# "bonona"
```

## string·rfind

`S.rfind(sub[, start[, end]])` returns the index of the substring `sub` within
S, like `S.find`, except that `rfind` returns the index of the substring's
_last_ occurrence.

```python
"bonbon".rfind("on")             # 4
"bonbon".rfind("on", None, 5)    # 1
"bonbon".rfind("on", 2, 5)       # -1
```

## string·rindex

`S.rindex(sub[, start[, end]])` returns the index of the substring `sub` within
S, like `S.index`, except that `rindex` returns the index of the substring's
_last_ occurrence.

```python
"bonbon".rindex("on")             # 4
"bonbon".rindex("on", None, 5)    # 1                           (in "bonbo")
"bonbon".rindex("on", 2, 5)       # error: substring not found  (in "nbo")
```

## string·rpartition

`S.rpartition(x)` is like `partition`, but splits `S` at the last occurrence of `x`.

```python
"one/two/three".partition("/")		# ("one/two", "/", "three")
```

## string·rsplit

`S.rsplit([sep[, maxsplit]])` splits a string into substrings like `S.split`,
except that when a maximum number of splits is specified, `rsplit` chooses the
rightmost splits.

```python
"banana".rsplit("n")                         # ["ba", "a", "a"]
"banana".rsplit("n", 1)                      # ["bana", "a"]
"one two  three".rsplit(None, 1)             # ["one two", "three"]
"".rsplit("n")                               # [""]
```

## string·rstrip

`S.rstrip()` returns a copy of the string S with trailing whitespace removed.

Like `strip`, it accepts an optional string parameter that specifies an
alternative set of Unicode code points to remove.

```python
"  hello  ".rstrip()                    # "  hello"
"  hello  ".rstrip("h o")               # "  hell"
```

## string·split

`S.split([sep [, maxsplit]])` returns the list of substrings of S,
splitting at occurrences of the delimiter string `sep`.

Consecutive occurrences of `sep` are considered to delimit empty
strings, so `'food'.split('o')` returns `['f', '', 'd']`.
Splitting an empty string with a specified separator returns `['']`.
If `sep` is the empty string, `split` fails.

If `sep` is not specified or is `None`, `split` uses a different
algorithm: it removes all leading spaces from S
(or trailing spaces in the case of `rsplit`),
then splits the string around each consecutive non-empty sequence of
Unicode white space characters.
If S consists only of white space, `S.split()` returns the empty list.

If `maxsplit` is given and non-negative, it specifies a maximum number of splits.

```python
"one two  three".split()                    # ["one", "two", "three"]
"one two  three".split(" ")                 # ["one", "two", "", "three"]
"one two  three".split(None, 1)             # ["one", "two  three"]
"banana".split("n")                         # ["ba", "a", "a"]
"banana".split("n", 1)                      # ["ba", "ana"]
"".split("n")                               # [""]
```

## string·elems

`S.elems()` returns an iterable value containing successive
1-byte substrings of S.
To materialize the entire sequence, apply `list(...)` to the result.

Example:

```python
list('Hello, 世界'.elems())  # ["H", "e", "l", "l", "o", ",", " ", "\xe4", "\xb8", "\x96", "\xe7", "\x95", "\x8c"]
```

See also: `string·elem_ords`.


## string·codepoints

`S.codepoints()` returns an iterable value containing the sequence of
substrings of S that each encode a single Unicode code point.
Each invalid code within the string is treated as if it encodes the
Unicode replacement character, U+FFFD.

By returning an iterable, not a list, the cost of decoding the string
is deferred until actually needed; apply `list(...)` to the result to
materialize the entire sequence.

Example:

```python
list('Hello, 世界'.codepoints())  # ['H', 'e', 'l', 'l', 'o', ',', ' ', '世', '界']

for cp in 'Hello, 世界'.codepoints():
   print(cp)  # prints 'H', 'e', 'l', 'l', 'o', ',', ' ', '世', '界'
```

See also: `string·codepoint_ords`.

<b>Implementation note:</b> `codepoints` is not provided by the Java implementation.

## string·splitlines

`S.splitlines([keepends])` returns a list whose elements are the
successive lines of S, that is, the strings formed by splitting S at
line terminators (currently assumed to be a single newline, `\n`,
regardless of platform).

The optional argument, `keepends`, is interpreted as a Boolean.
If true, line terminators are preserved in the result, though
the final element does not necessarily end with a line terminator.

As a special case, if S is the empty string,
`splitlines` returns the empty list.

```python
"one\n\ntwo".splitlines()       # ["one", "", "two"]
"one\n\ntwo".splitlines(True)   # ["one\n", "\n", "two"]
"".splitlines()                 # [] -- a special case
```

## string·startswith

`S.startswith(prefix[, start[, end]])` reports whether the string
`S[start:end]` has the specified prefix.

```python
"filename.star".startswith("filename")         # True
```

The `prefix` argument may be a tuple of strings, in which case the
function reports whether any one of them is a prefix.

```python
'abc'.startswith(('a', 'A'))                  # True
'ABC'.startswith(('a', 'A'))                  # True
'def'.startswith(('a', 'A'))                  # False
```

## string·strip

`S.strip()` returns a copy of the string S with leading and trailing whitespace removed.

It accepts an optional string argument:
`S.strip(cutset)` instead removes all leading
and trailing Unicode code points contained in `cutset`.

```python
"  hello  ".strip()                     # "hello"
"  hello  ".strip("h o")                # "ell"
```

## string·title

`S.title()` returns a copy of the string S with letters converted to title case.

Letters are converted to upper case at the start of words, lower case elsewhere.

```python
"hElLo, WoRlD!".title()                 # "Hello, World!"
"ǆenan".title()                        # "ǅenan" ("ǅ" is a single Unicode letter)
```

## string·upper

`S.upper()` returns a copy of the string S with letters converted to uppercase.

```python
"Hello, World!".upper()                 # "HELLO, WORLD!"
```
