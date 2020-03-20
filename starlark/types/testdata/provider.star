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

# compare
assert.ne(p, kwargs)
assert.ne(p, kwargs)

foo = p.resource.instance("foo", ami="valueA")
bar = p.resource.instance("bar", ami="valueA", disable_api_termination=False)
qux = p.resource.instance("qux", ami="valueB", disable_api_termination=True)

result = p.resource.instance.search("id", "foo")
assert.eq(len(result), 1)
assert.eq(result[0], foo)

assert.eq(len(p.resource.instance.search("ami", "valueA")), 2)
assert.eq(len(p.resource.instance.search("disable_api_termination", True)), 1)
assert.eq(len(p.resource.instance.search("disable_api_termination", False)), 1)

assert.eq(p.resource.instance.search("foo")[0], foo)