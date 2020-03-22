load("assert.star", "assert")

# constructor
foo = provisioner("file", source="conf/myapp.conf", destination="/etc/myapp.conf")
assert.eq(foo.source, "conf/myapp.conf")
assert.eq(foo.destination, "/etc/myapp.conf")


file = provisioner("file")
assert.eq(file.__kind__, "provisioner")
assert.eq(file.__type__, "file")
assert.eq(file.__name__, "id_2")

# attr
file.content = "foo"
assert.eq(file.content, "foo")

# attr names
assert.eq("content" in dir(file), True)

# hcl
assert.eq(hcl(file), "" +
'provisioner "file" {\n' + \
'  content = "foo"\n' + \
'}\n')


# type
assert.eq(type(file), "Provisioner<file>")
