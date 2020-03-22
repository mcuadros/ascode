load("assert.star", "assert")

b = backend("gcs")

# resource
assert.eq(b.__kind__, "backend")
assert.eq(b.__type__, "gcs")
assert.eq(b.__name__, "id_1")
assert.eq(type(b), "Backend<gcs>")

# attr
b.bucket = "tf-state-prod"
b.prefix = "terraform/state"

# attr names
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
