---
title: 'Lexical elements'
weight: 2
---

A Starlark program consists of one or more modules.
Each module is defined by a single UTF-8-encoded text file.

A complete grammar of Starlark can be found in [grammar.txt](https://github.com/google/starlark-go/blob/master/syntax/grammar.txt).
That grammar is presented piecemeal throughout this document
in boxes such as this one, which explains the notation:

```grammar {.good}
Grammar notation

- lowercase and 'quoted' items are lexical tokens.
- Capitalized names denote grammar productions.
- (...) implies grouping.
- x | y means either x or y.
- [x] means x is optional.
- {x} means x is repeated zero or more times.
- The end of each declaration is marked with a period.
```

The contents of a Starlark file are broken into a sequence of tokens of
five kinds: white space, punctuation, keywords, identifiers, and literals.
Each token is formed from the longest sequence of characters that
would form a valid token of each kind.

```grammar {.good}
File = {Statement | newline} eof .
```

*White space* consists of spaces (U+0020), tabs (U+0009), carriage
returns (U+000D), and newlines (U+000A).  Within a line, white space
has no effect other than to delimit the previous token, but newlines,
and spaces at the start of a line, are significant tokens.

*Comments*: A hash character (`#`) appearing outside of a string
literal marks the start of a comment; the comment extends to the end
of the line, not including the newline character.
Comments are treated like other white space.

*Punctuation*: The following punctuation characters or sequences of
characters are tokens:

```text
+    -    *    /    //   %    =
+=   -=   *=   /=   //=  %=   ==   !=
^    <    >    <<   >>   &    |
^=   <=   >=   <<=  >>=  &=   |=
.    ,    ;    :    ~    **
(    )    [    ]    {    }
```

*Keywords*: The following tokens are keywords and may not be used as
identifiers:

```text
and            elif           in             or
break          else           lambda         pass
continue       for            load           return
def            if             not            while
```

The tokens below also may not be used as identifiers although they do not
appear in the grammar; they are reserved as possible future keywords:

<!-- and to remain a syntactic subset of Python -->

```text
as             finally        nonlocal
assert         from           raise
class          global         try
del            import         with
except         is             yield
```

<b>Implementation note:</b>
The Go implementation permits `assert` to be used as an identifier,
and this feature is widely used in its tests.

*Identifiers*: an identifier is a sequence of Unicode letters, decimal
 digits, and underscores (`_`), not starting with a digit.
Identifiers are used as names for values.

Examples:

```text
None    True    len
x       index   starts_with     arg0
```

*Literals*: literals are tokens that denote specific values.  Starlark
has string, integer, and floating-point literals.

```text
0                               # int
123                             # decimal int
0x7f                            # hexadecimal int
0o755                           # octal int
0b1011                          # binary int

0.0     0.       .0             # float
1e10    1e+10    1e-10
1.1e10  1.1e+10  1.1e-10

"hello"      'hello'            # string
'''hello'''  """hello"""        # triple-quoted string
r'hello'     r"hello"           # raw string literal
```

Integer and floating-point literal tokens are defined by the following grammar:

```grammar {.good}
int         = decimal_lit | octal_lit | hex_lit | binary_lit .
decimal_lit = ('1' … '9') {decimal_digit} | '0' .
octal_lit   = '0' ('o'|'O') octal_digit {octal_digit} .
hex_lit     = '0' ('x'|'X') hex_digit {hex_digit} .
binary_lit  = '0' ('b'|'B') binary_digit {binary_digit} .

float     = decimals '.' [decimals] [exponent]
          | decimals exponent
          | '.' decimals [exponent]
          .
decimals  = decimal_digit {decimal_digit} .
exponent  = ('e'|'E') ['+'|'-'] decimals .

decimal_digit = '0' … '9' .
octal_digit   = '0' … '7' .
hex_digit     = '0' … '9' | 'A' … 'F' | 'a' … 'f' .
binary_digit  = '0' | '1' .
```

## String literals

A Starlark string literal denotes a string value. 
In its simplest form, it consists of the desired text 
surrounded by matching single- or double-quotation marks:

```python
"abc"
'abc'
```

Literal occurrences of the chosen quotation mark character must be
escaped by a preceding backslash. So, if a string contains several
of one kind of quotation mark, it may be convenient to quote the string
using the other kind, as in these examples:

```python
'Have you read "To Kill a Mockingbird?"'
"Yes, it's a classic."

"Have you read \"To Kill a Mockingbird?\""
'Yes, it\'s a classic.'
```

### String escapes

Within a string literal, the backslash character `\` indicates the
start of an _escape sequence_, a notation for expressing things that
are impossible or awkward to write directly.

The following *traditional escape sequences* represent the ASCII control
codes 7-13:

```
\a   \x07 alert or bell
\b   \x08 backspace
\f   \x0C form feed
\n   \x0A line feed
\r   \x0D carriage return
\t   \x09 horizontal tab
\v   \x0B vertical tab
```

A *literal backslash* is written using the escape `\\`.

An *escaped newline*---that is, a backslash at the end of a line---is ignored,
allowing a long string to be split across multiple lines of the source file.

```python
"abc\
def"			# "abcdef"
```

An *octal escape* encodes a single byte using its octal value.
It consists of a backslash followed by one, two, or three octal digits [0-7].
It is error if the value is greater than decimal 255.

```python
'\0'			# "\x00"  a string containing a single NUL byte
'\12'			# "\n"    octal 12 = decimal 10
'\101-\132'		# "A-Z"
'\119'			# "\t9"   = "\11" + "9"
```

<b>Implementation note:</b>
The Java implementation encodes strings using UTF-16,
so an octal escape encodes a single UTF-16 code unit.
Octal escapes for values above 127 are therefore not portable across implementations.
There is little reason to use octal escapes in new code.

A *hex escape* encodes a single byte using its hexadecimal value.
It consists of `\x` followed by exactly two hexadecimal digits [0-9A-Fa-f].

```python
"\x00"			# "\x00"  a string containing a single NUL byte
"(\x20)"		# "( )"   ASCII 0x20 = 32 = space

red, reset = "\x1b[31m", "\x1b[0m"	# ANSI terminal control codes for color
"(" + red + "hello" + reset + ")"	# "(hello)" with red text, if on a terminal
```

<b>Implementation note:</b>
The Java implementation does not support hex escapes.

An ordinary string literal may not contain an unescaped newline,
but a *multiline string literal* may spread over multiple source lines.
It is denoted using three quotation marks at start and end.
Within it, unescaped newlines and quotation marks (or even pairs of
quotation marks) have their literal meaning, but three quotation marks
end the literal. This makes it easy to quote large blocks of text with
few escapes.

```
haiku = '''
Yesterday it worked.
Today it is not working.
That's computers. Sigh.
'''
```

Regardless of the platform's convention for text line endings---for
example, a linefeed (\n) on UNIX, or a carriage return followed by a
linefeed (\r\n) on Microsoft Windows---an unescaped line ending in a
multiline string literal always denotes a line feed (\n).

Starlark also supports *raw string literals*, which look like an
ordinary single- or double-quotation preceded by `r`. Within a raw
string literal, there is no special processing of backslash escapes,
other than an escaped quotation mark (which denotes a literal
quotation mark), or an escaped newline (which denotes a backslash
followed by a newline). This form of quotation is typically used when
writing strings that contain many quotation marks or backslashes (such
as regular expressions or shell commands) to reduce the burden of
escaping:

```python
"a\nb"		# "a\nb"  = 'a' + '\n' + 'b'
r"a\nb"		# "a\\nb" = 'a' + '\\' + '\n' + 'b'

"a\
b"		# "ab"
r"a\
b"		# "a\\\nb"
```

It is an error for a backslash to appear within a string literal other
than as part of one of the escapes described above.

TODO: define indent, outdent, semicolon, newline, eof
