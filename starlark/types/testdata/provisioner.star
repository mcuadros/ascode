load("assert.star", "assert")

# constructor
foo = provisioner("file", source="conf/myapp.conf", destination="/etc/myapp.conf")
assert.eq(foo.source, "conf/myapp.conf")
assert.eq(foo.destination, "/etc/myapp.conf")

file = provisioner("file")

# attr
file.content = "foo"
assert.eq(file.content, "foo")
assert.eq(len(dir(file)), 3)

# hcl
assert.eq(hcl(file), "" +
'provisioner "file" {\n' + \
'  content = "foo"\n' + \
'}\n')


# type
assert.eq(type(file), "Provisioner<file>")
