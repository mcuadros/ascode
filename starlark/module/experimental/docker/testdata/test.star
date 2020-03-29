load('experimental/docker', 'docker')
load('assert.star', 'assert')

attr = docker.image("mcuadros/ascode", "latest")
assert.eq(attr.name, "docker.io/mcuadros/ascode")
assert.eq(attr.domain, "docker.io")
assert.eq(attr.path, "mcuadros/ascode")
assert.eq(dir(attr), ["domain", "name", "path", "tags", "version"])

image = docker.image("fedora", "latest")
assert.eq(image.name, "docker.io/library/fedora")
assert.eq(image.domain, "docker.io")
assert.eq(image.path, "library/fedora")
assert.eq(image.version(), "latest")

semver = docker.image("fedora", ">=22 <30")
assert.eq(semver.name, "docker.io/library/fedora")
assert.eq(semver.version(), "29")
assert.eq(semver.version(True), "docker.io/library/fedora:29")

prometheus = docker.image("quay.io/prometheus/prometheus", "1.8.x")
assert.eq(prometheus.name, "quay.io/prometheus/prometheus")
assert.eq(prometheus.version(), "v1.8.2")

tagNotFound = docker.image("fedora", "not-found")
assert.fails(lambda: tagNotFound.version(), 'tag "not-found" not found in repository')

