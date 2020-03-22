load("assert.star", "assert")

bar = "bar"

# context by kwargs
module = evaluate("evaluate/test.star", bar=bar)
assert.eq(str(module), '<module "test">')
assert.eq(module.foo, bar)

# context by dict
module = evaluate("evaluate/test.star", {"bar": bar})
assert.eq(str(module), '<module "test">')
assert.eq(module.foo, bar)

# context dict overrided by kwargs
module = evaluate("evaluate/test.star", {"bar": bar}, bar="foo")
assert.eq(str(module), '<module "test">')
assert.eq(module.foo, "foo")

# context dict with non strings
def contextNonString(): evaluate("evaluate/test.star", {1: bar})
assert.fails(contextNonString, "expected string keys in dict, got int at index 0")
