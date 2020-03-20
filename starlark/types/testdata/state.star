
load("assert.star", "assert")

b = backend("local")
b.path = "fixtures/modules/terraform.tfstate"

assert.eq(b.state().null.resource.resource.foo.triggers["foo"], "foo-value")
assert.eq(b.state("module.moduleA").null.resource.resource.bar.triggers["bar"], "bar-value")
assert.eq(b.state("module.moduleA.module.moduleB").null.resource.resource.qux.triggers["qux"], "qux-value")

c = backend("local")
c.path = "fixtures/state/terraform.tfstate"

cs = c.state()
assert.eq(cs.google.data.client_config.default.id, "2020-03-19 15:06:27.25614138 +0000 UTC")
assert.eq(cs.google.data.client_config.default.project, "project-foo")

cluster = cs.google.resource.container_cluster.primary
assert.eq(cluster.addons_config.network_policy_config.disabled, True)

release = cs.helm.resource.release["nats-operator"]
assert.eq(release.set[0].name, "cluster.auth.enabled")
assert.eq(release.set[1].name, "image.tag")
