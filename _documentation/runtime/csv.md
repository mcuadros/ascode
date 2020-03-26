---
title: 'encoding/csv'
---

csv reads comma-separated values files

## Index


* [def <b>read_all</b>(source, comma=",", comment="", lazy_quotes=False, trim_leading_space=False, fields_per_record=0, skip=0) [][]string](#def-icsvibread_allb)
* [def <b>write_all</b>(source,comma=",") string](#def-icsvibwrite_allb)


## Functions


#### def <i>csv</i>.<b>read_all</b>
```go
csv.read_all(source, comma=",", comment="", lazy_quotes=False, trim_leading_space=False, fields_per_record=0, skip=0) [][]string
```
read all rows from a source string, returning a list of string lists

###### Arguments

| name | type | description |
|------|------|-------------|
| `source` | `string` | input string of csv data |
| `comma` | `string` | comma is the field delimiter, defaults to "," (a comma).comma must be a valid character and must not be \r, \n,or the Unicode replacement character (0xFFFD). |
| `comment` | `string` | comment, if not "", is the comment character. Lines beginning with thecomment character without preceding whitespace are ignored.With leading whitespace the comment character becomes part of thefield, even if trim_leading_space is True.comment must be a valid character and must not be \r, \n,or the Unicode replacement character (0xFFFD).It must also not be equal to comma. |
| `lazy_quotes` | `bool` | If lazy_quotes is True, a quote may appear in an unquoted field anda non-doubled quote may appear in a quoted field. |
| `trim_leading_space` | `bool` | If trim_leading_space is True, leading white space in a field is ignored.This is done even if the field delimiter, comma, is white space. |
| `fields_per_record` | `int` | fields_per_record is the number of expected fields per record.If fields_per_record is positive, read_all requires each record tohave the given number of fields. If fields_per_record is 0, read_all sets it tothe number of fields in the first record, so that future records musthave the same field count. If fields_per_record is negative, no check ismade and records may have a variable number of fields. |
| `skip` | `int` | number of rows to skip, omitting from returned rows |



#### def <i>csv</i>.<b>write_all</b>
```go
csv.write_all(source,comma=",") string
```
write all rows from source to a csv-encoded string

###### Arguments

| name | type | description |
|------|------|-------------|
| `source` | `[][]string` | array of arrays of strings to write to csv |
| `comma` | `string` | comma is the field delimiter, defaults to "," (a comma).comma must be a valid character and must not be \r, \n,or the Unicode replacement character (0xFFFD). |



