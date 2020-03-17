load("assert.star", "assert")

p = provider("ignition", "1.1.0")

# attr
qux = p.data.user()
qux.uid = 42
assert.eq(qux.uid, 42)

qux.uid *= 2
assert.eq(qux.uid, 84)

# attr not-set
assert.eq(qux.name, None)

# attr not-exists
assert.fails(lambda: qux.foo, "Resource<data.ignition_user> has no .foo field or method")

# attr id
assert.eq(type(qux.id), "Computed")
assert.eq(str(qux.id), '"${data.ignition_user.id_2.id}"')

# attr output assignation
aws = provider("aws", "2.13.0")
def invalidOutput(): aws.data.instance().public_dns = "foo"
assert.fails(invalidOutput, "aws_instance: can't set computed public_dns attribute")

# attr output in asignation
web = aws.resource.instance()
web.ami = web.id
def invalidType(): web.get_password_data = web.id
assert.fails(invalidType, "expected bool, got string")

group = aws.resource.autoscaling_group()

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

# attr resource non dict
def attrResourceNonDict(): group.mixed_instances_policy = []
assert.fails(attrResourceNonDict, "expected dict, got list")

# attr collections
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
assert.eq(p.data.disk(), p.data.disk())
assert.ne(p.data.disk(device="foo"), p.data.disk())

# comparasion with nested
y = p.data.disk()
x = p.data.disk()

y.partition(start=42)
assert.ne(x, y)

x.partition(start=42)
assert.eq(x, y)

# comparasion with list
assert.ne(p.data.user(groups=["foo"]), p.data.user())
assert.eq(p.data.user(groups=["foo"]), p.data.user(groups=["foo"]))

# constructor with name
quux = p.data.user("quux")
assert.eq(str(quux.id), '"${data.ignition_user.quux.id}"')

# constructor from kwargs
bar = p.data.user(uid=42, system=True)
assert.eq(bar.uid, 42)
assert.eq(bar.system, True)

# constructor from kwargs with name
fred = p.data.user("fred", uid=42, system=True)
assert.eq(fred.uid, 42)
assert.eq(fred.system, True)
assert.eq(str(fred.id), '"${data.ignition_user.fred.id}"')

# constructor from dict
foo = p.data.user({"uid": 42, "system": True})
assert.eq(foo.uid, 42)
assert.eq(foo.system, True)

# constructor from dict with name
baz = p.data.user("baz", {"uid": 42, "system": True})
assert.eq(baz.uid, 42)
assert.eq(baz.system, True)
assert.eq(str(baz.id), '"${data.ignition_user.baz.id}"')

assert.eq(bar, foo)
assert.eq(foo, p.data.user(foo.__dict__))

# constructor errors
def consNonDict(): p.data.user(1)
assert.fails(consNonDict, "resource: expected string or dict, got int")

def consNonNameDict(): p.data.user(1, 1)
assert.fails(consNonNameDict, "resource: expected string, got int")

def consNameDict(): p.data.user("foo", 1)
assert.fails(consNameDict, "resource: expected dict, got int")

def consKwargsNonName(): p.data.user(1, uid=42)
assert.fails(consKwargsNonName, "resource: expected string, got int")

# full coverage
user = p.data.user()
user.name = "foo"
user.uid = 42
user.groups = ["foo", "bar"]
user.system = True

assert.eq(type(user), "Resource<data.ignition_user>")
assert.eq(user.__dict__, {
    "name": "foo",
    "uid": 42,
    "groups": ["foo", "bar"],
    "system": True,
})

disk = p.data.disk()

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
userA = p.data.user()
userB = p.data.user()
userA.depends_on(userB)

def dependsOnNonResource(): userA.depends_on(42)
assert.fails(dependsOnNonResource, "expected Resource<\\[data|resource\\].\\*>, got int")

def dependsOnNestedResource(): userA.depends_on(disk.partition())
assert.fails(dependsOnNestedResource, "expected Resource<\\[data|resource\\].\\*>, got Resource<nested.partition>")

def dependsOnItself(): userA.depends_on(userA)
assert.fails(dependsOnItself, "can't depend on itself")

# __provider__
assert.eq(web.__provider__, aws)
assert.eq(baz.__provider__, p)
assert.eq(userA.__provider__, p)
assert.eq(home.__provider__, p)