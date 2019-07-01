
load("assert.star", "assert")

p = provider("aws", "2.13.0")
assert.eq(p.version, "2.13.0")

assert.eq(len(dir(p.data)), 131)
assert.eq(len(dir(p.resource)), 506)

resources = dir(p.resource)
assert.contains(resources, "instance")
assert.eq(type(p.resource.instance), "builtin_function_or_method")