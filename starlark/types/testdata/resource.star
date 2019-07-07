load("assert.star", "assert")

p = provider("ignition", "1.1.0")

# attr
qux = p.data.user()
qux.uid = 42
assert.eq(qux.uid, 42)

# attr not-set
assert.eq(qux.name, None)

# attr not-exists
assert.fails(lambda: qux.foo, "data has no .foo field or method")

# attr id
assert.eq(type(qux.id), "computed")
assert.eq(str(qux.id), '"$${data.ignition_user.id_5.id}"')

# attr output assignation
aws = provider("aws", "2.13.0")
def invalidOutput(): aws.data.instance().public_dns = "foo"
assert.fails(invalidOutput, "aws_instance: can't set computed public_dns attribute") 

# attr output in asignation
web = aws.resource.instance()
web.ami = web.id
def invalidType(): web.get_password_data = web.id
assert.fails(invalidType, "expected bool, got string") 

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

# constructor from kwargs
bar = p.data.user(uid=42, system=True)
assert.eq(bar.uid, 42)
assert.eq(bar.system, True)

# constructor from dict
foo = p.data.user({"uid": 42, "system": True})
assert.eq(foo.uid, 42)
assert.eq(foo.system, True)

assert.eq(bar, foo)
assert.eq(foo, p.data.user(foo.__dict__))

# full coverage
user = p.data.user()
user.name = "foo"
user.uid = 42
user.groups = ["foo", "bar"]
user.system = True

assert.eq(type(user), "data")
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