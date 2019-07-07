load("assert.star", "assert")

aws = provider("aws", "2.13.0")
aws.region = "us-west-2"

ubuntu = aws.data.ami()
ubuntu.most_recent = True
ubuntu.filter(name = "name", values = ["ubuntu/images/hvm-ssd/ubuntu-trusty-14.04-amd64-server-*"])
ubuntu.filter(name = "virtualization-type", values = ["hvm"])
ubuntu.owners = ["099720109477"]

web = aws.resource.instance(instance_type = "t2.micro")

#web.instance_type = "t2.micro"
#web.ami = ami.id

template = aws.resource.launch_template()
template.name_prefix = "example"
template.instance_type = "c5.larger"

group = aws.resource.autoscaling_group()
group.availability_zones = ["us-east-1a"]
group.desired_capacity = 1
group.max_size = 1
group.min_size = 1
group.mixed_instances_policy = {
    "launch_template": {
        "launch_template_specification": {
            "launch_template_id": "bar",
        },
    },
}

ami2 = aws.data.ami()
ami2.most_recent = True

print(hcl(aws))
