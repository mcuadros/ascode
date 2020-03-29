# Runtime Modules
# AsCode comes with a variety of modules available like `http`, `math`, 
# `encoding/json`, etc. All this [modules](/docs/reference/) are available 
# runtime through the [`load`](/docs/starlark/statements/#load-statements) 
# statement. This example shows the usage of this modules and some others.
# <!--more-->

# ## Basic Module
# The load statement expects at least two arguments; the first is the name of
# the module, and the same the symbol to extract to it. The runtime modules 
# always define a symbol called equals to the last part of the module name.
load("encoding/base64", "base64")
load("http", "http")

# This modules are very usuful to do basic operations such as encoding of
# strings, like in this case to `base64` or to make HTTP requests.
dec = base64.encode("ascode is amazing")

msg = http.get("https://httpbin.org/base64/%s" % dec)
print(msg.body())

# ### Output
"""sh
ascode is amazing
"""

# ## Advanced Modules
# Also, AsCode has some more specif modules, like the `docker` module. The 
# docker modules allow you to manipulate docker image names.
load("experimental/docker", "docker")

# A docker image tag can be defined using semver, instead of using the infamous
# 'latest' tag, or fixing a particular version. This allows us to be up-to-date
# without breaking our deployment.
golang = docker.image("golang", "1.13.x")

# We can use this in the definition of resources, allowing use to upgrade
# the version of our containers in every `terraform apply`
p = tf.provider("docker", "2.7.0")
container = p.resource.container("golang", image=golang.version(full=True))

# version queries the docker repository and returns the correct tag.
print(hcl(container))

# ### Output
"""hcl
resource "docker_container" "foo" {
  provider = docker.id_01E4KHW2RSW0FQM93KN5W70Y42
  image    = "docker.io/library/golang:1.13.9"
}
"""