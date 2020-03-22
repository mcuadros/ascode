load("assert.star", "assert")

b = backend("local")
b.path = "fixtures/modules/terraform.tfstate"

s = b.state()
assert.ne(s, None)
assert.ne(s["null"], None)
assert.eq(s["null"]["resource"]["resource"]["foo"].triggers["foo"], "foo-value")
assert.eq(b.state("module.moduleA")["null"]["resource"]["resource"]["bar"].triggers["bar"], "bar-value")
assert.eq(b.state("module.moduleA.module.moduleB")["null"]["resource"]["resource"]["qux"].triggers["qux"], "qux-value")

c = backend("local")
c.path = "fixtures/state/terraform.tfstate"

s = c.state()
assert.ne(s["google"]["data"]["client_config"], None)
assert.eq(s["google"]["data"]["client_config"]["default"].id, "2020-03-19 15:06:27.25614138 +0000 UTC")
assert.eq(s["google"]["data"]["client_config"]["default"].project, "project-foo")

cluster = s["google"]["resource"]["container_cluster"]["primary"]
assert.eq(cluster.addons_config.network_policy_config.disabled, True)

release = s["helm"]["resource"]["release"]["nats-operator"]
assert.eq(release.set[0].name, "cluster.auth.enabled")
assert.eq(release.set[1].name, "image.tag")
