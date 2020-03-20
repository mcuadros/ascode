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