load("assert.star", "assert")

b = backend("gcs")
b.bucket = "tf-state-prod"
b.prefix = "terraform/state"

assert.eq(hcl(b), "" +
'terraform {\n' + \
'  backend "gcs" {\n' + \
'    bucket = "tf-state-prod"\n' + \
'    prefix = "terraform/state"\n' + \
'  }\n' + \
'}\n')
