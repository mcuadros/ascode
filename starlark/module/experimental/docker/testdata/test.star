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

full = docker.image("fedora", "24")
assert.eq(full.name, "docker.io/library/fedora")
assert.eq(full.version(True), "docker.io/library/fedora:24")

semver = docker.image("fedora", ">=22 <30")
assert.eq(semver.name, "docker.io/library/fedora")
assert.eq(semver.version(), "29")

golang = docker.image("golang", "1.13.x")
assert.eq(golang.name, "docker.io/library/golang")
assert.eq(golang.version(), "1.13.8")

tagNotFound = docker.image("fedora", "not-found")
assert.eq(tagNotFound.name, "docker.io/library/fedora")

def tagNotExistant(): tagNotFound.version()
assert.fails(tagNotExistant,'tag "not-found" not found in repository')

