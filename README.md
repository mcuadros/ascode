# AsCode - The Real Infrastructure as Code

**AsCode** is a tool for define infrastructure as code using the [Starlark](https://github.com/google/starlark-go/blob/master/doc/spec.md) language on top of [Terraform](https://github.com/hashicorp/terraform). It allows to describe your infrastructure using an expressive language in Terraform without writing a single line of [HCL](https://www.terraform.io/docs/configuration/syntax.html), meanwhile, you have the complete ecosystem of [providers](https://www.terraform.io/docs/providers/index.html)

### Why?

Terraform is a great tool, with support for almost everything you can imagine, making it the industry leader. Terraform is based on HCL, a JSON-alike declarative language, with a very limited control flow functionalities. IMHO, to really unleash the power of the IaC, a powerful, expressive language should be used, where basic elements like loops or functions are first class citizens.


### What is Starlark?

> Starlark is a dialect of Python intended for use as a configuration language. A Starlark interpreter is typically embedded within a larger application, and this application may define additional domain-specific functions and data types beyond those provided by the core language. For example, Starlark is embedded within (and was originally developed for) the Bazel build tool, and Bazel's build language is based on Starlark.

## Examples

### Simple

Creating am Amazon EC2 Instance is as easy as:

```pyhon
aws = tf.provider("aws", "2.13.0")
aws.region = "us-west-2"

aws.resource.instance(instance_type ="t2.micro", ami="ami-2757f631")
```
### Using functions

In this example we create 40 instances, 20 using ubuntu and 20 using ECS.

```python
aws = tf.provider("aws")
aws.region = "us-west-2"

# It creates a new instance for the given name, distro and type.
def new_instance(name, distro, type="t2.micro"):
    instance = aws.resource.instance(name)
    instance.instance_type = type
    instance.ami = get_ami_id(distro)

    return instance

amis = {}
ami_names_owners = {
    "ubuntu": ["ubuntu/images/*/ubuntu-xenial-16.04-amd64-server-*", "099720109477"],
    "ecs": ["*amazon-ecs-optimized", "591542846629"],
}

# We create the AMI data-source for the given distro.
def get_ami_id(distro):
    if distro in amis:
        return amis[distro]

    data = ami_names_owners[distro]

    ami = aws.data.ami(distro)
    ami.most_recent = True
    ami.filter(name="name", values=[data[0]])
    ami.filter(name="virtualization-type", values=["hvm"])
    ami.owners = [data[1]]

    amis[distro] = ami.id
    return ami.id

# Creates 20 instances of each distro.
for i in range(20):
    new_instance("ubuntu_%d" % i, "ubuntu")
    new_instance("ecs_%d" % i, "ecs")
 ```

### Using the runtime

ascode comes with a built-in runtime with functions to work with `yaml`, `json`, `http`, etc. Take a look to the [documentation](/_documentation/runtime).

```
load("encoding/base64", "base64")
load("http", "http")

dec = base64.encode("ascode is amazing")

msg = http.get("https://httpbin.org/base64/%s" % dec)
print(msg.body())
```


## Installation

The recommended way to install *ascode* it's download the binary from the [releases](https://github.com/mcuadros/ascode/releases) section.


## License

GPL-3.0, see [LICENSE](LICENSE)
