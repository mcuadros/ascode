load("assert.star", "assert")

p = tf.provider("aws", "2.13.0")
assert.eq(p.__kind__, "provider")
assert.eq(p.__type__, "aws")
assert.eq(p.__name__, "id_1")
assert.eq(p.__version__, "2.13.0")

# attr names
assert.eq("__provider__" in dir(p), False)
assert.eq("depends_on" in dir(p), False)
assert.eq("add_provisioner" in dir(p), False)
assert.eq("__version__" in dir(p), True)
assert.eq("data" in dir(p), True)
assert.eq("resource" in dir(p), True)

# attr
assert.eq(len(dir(p.data)), 133)
assert.eq(len(dir(p.resource)), 508)

resources = dir(p.resource)
assert.contains(resources, "instance")

# types
assert.eq(type(p), "Provider")
assert.eq(type(p.resource), "ResourceCollectionGroup")
assert.eq(type(p.resource.instance), "ResourceCollection")
assert.eq(type(p.resource.instance()), "Resource<resource>")
assert.eq(type(p.data.ami().filter()), "Resource<nested>")

# string
assert.eq(str(p), "Provider<aws>")
assert.eq(str(p.resource), "ResourceCollectionGroup<aws.resource>")
assert.eq(str(p.resource.instance), "ResourceCollection<aws.resource.aws_instance>")
assert.eq(str(p.resource.instance()), "Resource<aws.resource.aws_instance>")
assert.eq(str(p.data.ami().filter()), "Resource<aws.data.aws_ami.filter>")


assert.eq(len(p.resource.instance), 2)

p.region = "us-west-2"
assert.eq(p.region, "us-west-2")

alias = tf.provider("aws", "2.13.0", "alias")
assert.eq(alias.__name__, "alias")
assert.eq(alias.__version__, "2.13.0")

kwargs = tf.provider("aws", region="foo")
assert.eq(kwargs.region, "foo")

# ResourceCollectionGroup
assert.eq("__kind__" in dir(p.resource), True)
assert.eq(p.resource.__kind__, "resource")
assert.eq("__provider__" in dir(p.resource), True)
assert.eq(p.resource.__provider__, p)

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

