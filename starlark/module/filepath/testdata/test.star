load('path/filepath', 'filepath')
load('assert.star', 'assert')

assert.eq(filepath.separator, "/")

matches = filepath.glob("/tmp/*")
assert.eq(True, len(matches) > 0)

assert.eq(True, len(filepath.abs("foo")) > 3)

assert.eq(filepath.clean("bar///foo.md"), "bar/foo.md")
assert.eq(filepath.base("bar/foo.md"), "foo.md")
assert.eq(filepath.dir("bar/foo.md"), "bar")
assert.eq(filepath.ext("foo.md"), ".md")

assert.eq(filepath.is_abs("foo"), False)
assert.eq(filepath.is_abs("/foo"), True)

assert.eq(filepath.join(["foo", "bar"]), "foo/bar")
assert.eq(filepath.rel("/a", "/a/b/c"), "b/c")