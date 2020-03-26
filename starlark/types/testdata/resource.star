load("assert.star", "assert")

ignition = tf.provider("ignition", "1.1.0")



# attr
qux = ignition.data.user()
qux.uid = 42
assert.eq(qux.uid, 42)

qux.uid *= 2
assert.eq(qux.uid, 84)

# attr names 
assert.eq("uid" in dir(qux), True)

# attr names in data sources
assert.eq("depends_on" in dir(qux), False)
assert.eq("add_provisioner" in dir(qux), False)
assert.eq("__provider__" in dir(qux), True)
assert.eq("__type__" in dir(qux), True)
assert.eq("__name__" in dir(qux), True)
assert.eq("__kind__" in dir(qux), True)
assert.eq("__dict__" in dir(qux), True)

# attr not-set
assert.eq(qux.name, None)

# attr not-exists
assert.fails(lambda: qux.foo, "Resource<data> has no .foo field or method")

# attr id
assert.eq(type(qux.id), "Attribute<string>")
assert.eq(str(qux.id), '"${data.ignition_user.id_2.id}"')
aws = tf.provider("aws", "2.13.0")

# attr output assignation
def invalidOutput(): aws.data.instance().public_dns = "foo"
assert.fails(invalidOutput, "Resource<aws.data.aws_instance>: can't set computed public_dns attribute")


# attr output in asignation
web = aws.resource.instance()
web.ami = web.id
def invalidType(): web.get_password_data = web.id
assert.fails(invalidType, "expected bool, got string")

group = aws.resource.autoscaling_group()

# attr names in resources
assert.eq("depends_on" in dir(web), True)
assert.eq("add_provisioner" in dir(web), True)
assert.eq("__provider__" in dir(web), True)
assert.eq("__type__" in dir(web), True)
assert.eq("__name__" in dir(web), True)
assert.eq("__kind__" in dir(web), True)
assert.eq("__dict__" in dir(web), True)

# attr optional computed
assert.eq(str(group.name), '"${aws_autoscaling_group.id_6.name}"')

group.name = "foo"
assert.eq(group.name, "foo")

# attr resource
group.mixed_instances_policy = {
    "launch_template": {
        "launch_template_specification": {
            "launch_template_id": "bar",
        },
    },
}

assert.eq(group.mixed_instances_policy.launch_template.launch_template_specification.launch_template_id, "bar")

# attr collections
assert.eq("__provider__" in dir(web.network_interface), True)
assert.eq("__type__" in dir(web.network_interface), True)
assert.eq("__kind__" in dir(web.network_interface), True)

web.network_interface = [
    {"network_interface_id": "foo"},
    {"network_interface_id": "bar"},
]

assert.eq(len(web.network_interface), 2)
assert.eq(web.network_interface[0].network_interface_id, "foo")
assert.eq(web.network_interface[1].network_interface_id, "bar")

# attr collections clears list
web.network_interface = [
    {"network_interface_id": "qux"},
]

assert.eq(len(web.network_interface), 1)
assert.eq(web.network_interface[0].network_interface_id, "qux")

# attr collection non list
def attrCollectionNonList(): web.network_interface = {}
assert.fails(attrCollectionNonList, "expected list, got dict")

# attr collection non dict elements
def attrCollectionNonDictElement(): web.network_interface = [{}, 42]
assert.fails(attrCollectionNonDictElement, "1: expected dict, got int")

# comparasion simple values
assert.eq(ignition.data.disk(), ignition.data.disk())
assert.ne(ignition.data.disk(device="foo"), ignition.data.disk())

# comparasion with nested
y = ignition.data.disk()
x = ignition.data.disk()

y.partition(start=42)
assert.ne(x, y)

x.partition(start=42)
assert.eq(x, y)

# comparasion with list
assert.ne(ignition.data.user(groups=["foo"]), ignition.data.user())
assert.eq(ignition.data.user(groups=["foo"]), ignition.data.user(groups=["foo"]))

# constructor with name
quux = ignition.data.user("quux")
assert.eq(str(quux.id), '"${data.ignition_user.quux.id}"')

# constructor from kwargs
bar = ignition.data.user(uid=42, system=True)
assert.eq(bar.uid, 42)
assert.eq(bar.system, True)

# constructor from kwargs with name
fred = ignition.data.user("fred", uid=42, system=True)
assert.eq(fred.uid, 42)
assert.eq(fred.system, True)
assert.eq(str(fred.id), '"${data.ignition_user.fred.id}"')

# constructor from dict
foo = ignition.data.user({"uid": 42, "system": True})
assert.eq(foo.uid, 42)
assert.eq(foo.system, True)

# constructor from dict with name
baz = ignition.data.user("baz", {"uid": 42, "system": True})
assert.eq(baz.uid, 42)
assert.eq(baz.system, True)
assert.eq(str(baz.id), '"${data.ignition_user.baz.id}"')

# constructor from dict with name and kwargs
baz = ignition.data.user("baz", {"uid": 42, "system": True}, uid=84)
assert.eq(baz.uid, 84)
assert.eq(baz.system, True)
assert.eq(str(baz.id), '"${data.ignition_user.baz.id}"')


assert.eq(bar, foo)
assert.eq(foo, ignition.data.user(foo.__dict__))

# constructor errors
def consNonDict(): ignition.data.user(1)
assert.fails(consNonDict, "resource: expected string or dict, got int")

def consNonNameDict(): ignition.data.user(1, 1)
assert.fails(consNonNameDict, "resource: expected string, got int")

def consNameDict(): ignition.data.user("foo", 1)
assert.fails(consNameDict, "resource: expected dict, got int")

def consKwargsNonName(): ignition.data.user(1, uid=42)
assert.fails(consKwargsNonName, "resource: expected string or dict, got int")

# full coverage
user = ignition.data.user()
user.name = "foo"
user.uid = 42
user.groups = ["foo", "bar"]
user.system = True

assert.eq(str(user), "Resource<ignition.data.ignition_user>")
assert.eq(user.__dict__, {
    "name": "foo",
    "uid": 42,
    "groups": ["foo", "bar"],
    "system": True,
})

disk = ignition.data.disk()

root = disk.partition()
root.label = "root"
root.start = 2048
root.size = 4 * 1024 * 1024

home = disk.partition()
home.label = "home"
home.start = root.size + root.start
home.size = 4 * 1024 * 1024

assert.eq(disk.__dict__, {
    "partition": [{
        "label": "root",
        "start": 2048,
        "size": 4194304
    }, {
        "start": 4196352,
        "size": 4194304,
        "label": "home"
    }]
})


# depends_on
instanceA = aws.resource.instance()
instanceB = aws.resource.instance()
instanceA.depends_on(instanceB)

def dependsOnNonResource(): instanceA.depends_on(42)
assert.fails(dependsOnNonResource, "expected Resource<\\[data|resource\\]>, got int")

def dependsOnNestedResource(): instanceA.depends_on(disk.partition())
assert.fails(dependsOnNestedResource, "expected Resource<\\[data|resource\\]>, got Resource<nested.partition>")

def dependsOnItself(): instanceA.depends_on(instanceA)
assert.fails(dependsOnItself, "can't depend on itself")

# __provider__
assert.eq(web.__provider__, aws)
assert.eq(baz.__provider__, ignition)
assert.eq(instanceA.__provider__, aws)
assert.eq(home.__provider__, ignition)
assert.eq(aws.resource.instance.__provider__, aws)

# __kind__
assert.eq(ignition.data.user().__kind__, "data")
assert.eq(aws.resource.instance.__kind__, "resource")
assert.eq(aws.resource.instance().__kind__, "resource")
assert.eq(aws.resource.autoscaling_group().mixed_instances_policy.__kind__, "nested")
assert.eq(web.network_interface.__kind__, "nested")

# __type__
assert.eq(ignition.data.user().__type__, "ignition_user")
assert.eq(aws.resource.instance.__type__, "aws_instance")
assert.eq(aws.resource.instance().__type__, "aws_instance")
assert.eq(aws.resource.autoscaling_group().mixed_instances_policy.__type__, "mixed_instances_policy")
assert.eq(web.network_interface.__type__, "network_interface")

# __name__
assert.eq(ignition.data.user().__name__, "id_30")
assert.eq(aws.resource.instance().__name__, "id_31")
assert.eq(ignition.data.user("given").__name__, "given")

# __call__
assert.eq(ignition.data.user().__name__, "id_32")
assert.eq(ignition.data.user("foo").__name__, "foo")
assert.eq(ignition.data.user(uid=42).uid, 42)
assert.eq(ignition.data.user({"uid": 42}).uid, 42)

foo = ignition.data.user("foo", {"uid": 42})
assert.eq(foo.__name__, "foo")
assert.eq(foo.uid, 42)

foo = ignition.data.user("foo", uid=42)
assert.eq(foo.__name__, "foo")
assert.eq(foo.uid, 42)