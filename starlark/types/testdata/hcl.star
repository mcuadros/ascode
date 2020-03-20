load("assert.star", "assert")

helm = tf.provider("helm", "1.0.0", "default")
helm.kubernetes.token = "foo"

# hcl
assert.eq(hcl(helm), "" +
'provider "helm" {\n' + \
'  alias   = "default"\n' + \
'  version = "1.0.0"\n\n' + \
'  kubernetes {\n' + \
'    token = "foo"\n' + \
'  }\n' + \
'}\n\n')