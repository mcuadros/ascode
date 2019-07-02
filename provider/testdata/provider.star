load("assert.star", "assert")

p = provider("aws", "2.13.0")
assert.eq(p.version, "2.13.0")

assert.eq(len(dir(p.data)), 131)
assert.eq(len(dir(p.resource)), 506)

resources = dir(p.resource)
assert.contains(resources, "instance")
assert.eq(type(p.resource.instance), "aws_instance_collection")

p.resource.instance("foo")
p.resource.instance("bar")
assert.eq(len(p.resource.instance), 2)