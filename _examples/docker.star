load("experimental/docker", "docker")

p = tf.provider("docker", "2.7.0", "foo")

# using docker.image semver can be used to choose the docker image, `
golang = docker.image("golang", "1.13.x")

foo = p.resource.container("foo")
foo.name = "foo"

# version queries the docker repository and returns the correct tag.
foo.image = golang.version(full=True)