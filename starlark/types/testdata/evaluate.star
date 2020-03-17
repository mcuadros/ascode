load("assert.star", "assert")

bar = "bar"
module = evaluate("evaluate/test.star", bar=bar)
assert.eq(str(module), '<module "test">')
assert.eq(module.foo, bar)