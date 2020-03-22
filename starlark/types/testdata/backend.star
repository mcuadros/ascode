load("assert.star", "assert")

b = backend("gcs")

# resource
assert.eq(b.__kind__, "backend")
assert.eq(b.__type__, "gcs")
assert.eq(type(b), "Backend<gcs>")

# attr
b.bucket = "tf-state-prod"
b.prefix = "terraform/state"

# attr names
assert.eq("__provider__" in dir(b), False)
assert.eq("__name__" in dir(b), False)
assert.eq("depends_on" in dir(b), False)
assert.eq("add_provisioner" in dir(b), False)
assert.eq("state" in dir(b), True)
assert.eq("bucket" in dir(b), True)

# hcl
assert.eq(hcl(b), "" +
'terraform {\n' + \
'  backend "gcs" {\n' + \
'    bucket = "tf-state-prod"\n' + \
'    prefix = "terraform/state"\n' + \
'  }\n' + \
'}\n\n')
