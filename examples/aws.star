aws = provider("aws")
print(dir(aws))

ami = aws.data.ami("ubuntu")
ami.most_recent = True
ami.filter("name", values=["ubuntu/images/hvm-ssd/ubuntu-trusty-14.04-amd64-server-*"])
ami.filter("virtualization-type", values=["hvm"])
ami.owners = ["099720109477"]
print(ami.__dict__)


web = aws.resource.instance("web", instance_type="t2.micro")

#web.instance_type = "t2.micro"
#web.ami = ami.id


template = aws.resource.launch_template("example")
template.name_prefix = "example"
template.instance_type = "c5.larger"

group = aws.resource.autoscaling_group("example")
group.availability_zones = ["us-east-1a"]
group.desired_capacity = 1
group.max_size = 1
group.min_size = 1
#group.mixed_instances_policy.launch_template.launch_template_specification.launch_template_id = "foo"

group.mixed_instances_policy = {
    "launch_template": {
        "launch_template_specification": {
            "launch_template_id": "bar"
        }
    }
}

print(group.__dict__)