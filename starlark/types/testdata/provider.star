load("assert.star", "assert")

p = provider("aws", "2.13.0")
assert.eq(p.version, "2.13.0")

assert.eq(len(dir(p.data)), 131)
assert.eq(len(dir(p.resource)), 506)

resources = dir(p.resource)
assert.contains(resources, "instance")
assert.eq(type(p), "Provider<aws>")
assert.eq(type(p.resource.instance), "ResourceCollection<resource.aws_instance>")
assert.eq(type(p.resource.instance()), "Resource<resource.aws_instance>")

p.resource.instance()
assert.eq(len(p.resource.instance), 2)

p.region = "us-west-2"
assert.eq(p.region, "us-west-2")

alias = provider("aws", "2.13.0", "alias")
assert.eq(alias.alias, "alias")
assert.eq(alias.version, "2.13.0")

kwargs = provider("aws", region="foo")
assert.eq(kwargs.region, "foo")
