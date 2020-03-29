# Using functions
# This example illustrates how with the usage through the usage of functions, 
# we can simplify and improve the readability of our infrastructure declaration.
# <!--more-->

# Instantiates a new AWS provider, `aws` will be available in the context of 
# the functions. 
aws = tf.provider("aws", region="us-west-2")

# Every instance requires an `ami` data source; this data source contains a
# very specif configuration like the ID of the owner or a name pattern. So
# we define a dictionary with the different values we want to use.
ami_names_owners = {
    "ubuntu": ["ubuntu/images/*/ubuntu-xenial-16.04-amd64-server-*", "099720109477"],
    "ecs": ["*amazon-ecs-optimized", "591542846629"],
}

# `get_ami` returns the ami for the given `distro`. It searches in the
# `ResouceCollection` of the `ami` data source, if finds the `ami` it simply
# returns it; if not creates a new one using the data from the `ami_names_owners` 
# dictionary. 
def get_ami(distro):
    amis = aws.data.ami.search(distro)
    if len(amis) != 0:
        return amis[0]

    data = ami_names_owners[distro]

    ami = aws.data.ami(distro)
    ami.most_recent = True
    ami.filter(name="name", values=[data[0]])
    ami.filter(name="virtualization-type", values=["hvm"])
    ami.owners = [data[1]]

    return ami

# `new_instance` instantiates a new `instance` for the given name, distro and
# type, the type has a default value `t2.micro`. The `distro` value is resolved
# to an `ami` resource using the previously defined function `get_ami`.
def new_instance(name, distro, type="t2.micro"):
    instance = aws.resource.instance(name)
    instance.instance_type = type
    instance.ami = get_ami(distro).id

# Now using a basic `for` loop we can instantiate 5 different web servers, 
# where the even machines are using ubuntu and the odd ones ecs.
for i in range(5):
    distro = "ubuntu"
    if i % 2: 
        distro = "ecs"

    new_instance("web_%d" % i, distro)

# ### Output
# If we execute this script with the flag `--print-hcl` the result shuld be 
# something like this:

"""hcl
provider "aws" {
  alias   = "id_01E4KEA5ZAA1PYERQ8KM5D04GC"
  version = "2.13.0"
  region  = "us-west-2"
}

data "aws_ami" "ubuntu" {
  provider    = aws.id_01E4KEA5ZAA1PYERQ8KM5D04GC
  most_recent = true
  owners      = ["099720109477"]

  filter {
    name   = "name"
    values = ["ubuntu/images/*/ubuntu-xenial-16.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

data "aws_ami" "ecs" {
  provider    = aws.id_01E4KEA5ZAA1PYERQ8KM5D04GC
  most_recent = true
  owners      = ["591542846629"]

  filter {
    name   = "name"
    values = ["*amazon-ecs-optimized"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

resource "aws_instance" "web_0" {
  provider      = aws.id_01E4KEA5ZAA1PYERQ8KM5D04GC
  ami           = "${data.aws_ami.ubuntu.id}"
  instance_type = "t2.micro"
}

resource "aws_instance" "web_1" {
  provider      = aws.id_01E4KEA5ZAA1PYERQ8KM5D04GC
  ami           = "${data.aws_ami.ecs.id}"
  instance_type = "t2.micro"
}

resource "aws_instance" "web_2" {
  provider      = aws.id_01E4KEA5ZAA1PYERQ8KM5D04GC
  ami           = "${data.aws_ami.ubuntu.id}"
  instance_type = "t2.micro"
}

resource "aws_instance" "web_3" {
  provider      = aws.id_01E4KEA5ZAA1PYERQ8KM5D04GC
  ami           = "${data.aws_ami.ecs.id}"
  instance_type = "t2.micro"
}

resource "aws_instance" "web_4" {
  provider      = aws.id_01E4KEA5ZAA1PYERQ8KM5D04GC
  ami           = "${data.aws_ami.ubuntu.id}"
  instance_type = "t2.micro"
}
"""
