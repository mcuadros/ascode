load("os", "os")
load("assert.star", "assert")

aws = provider("aws", "2.13.0")
aws.region = "us-west-2"

# Based on:
# https://www.terraform.io/docs/providers/aws/r/instance.html#example
vpc = aws.resource.vpc()
vpc.cidr_block = "172.16.0.0/16"
vpc.tags = {"Name": "tf-example"}

subnet = aws.resource.subnet()
subnet.vpc_id = vpc.id
subnet.cidr_block = "172.16.0.0/24"
subnet.availability_zone = "us-west-2a"
subnet.tags = {"Name": "tf-example"}

iface = aws.resource.network_interface()
iface.subnet_id = subnet.id
iface.private_ips = ["172.16.10.100"]
iface.tags = {"Name": "primary_network_iterface"}

ubuntu = aws.data.ami()
ubuntu.most_recent = True
ubuntu.filter(name = "name", values = ["ubuntu/images/hvm-ssd/ubuntu-trusty-14.04-amd64-server-*"])
ubuntu.filter(name = "virtualization-type", values = ["hvm"])
ubuntu.owners = ["099720109477"]

instance = aws.resource.instance()
instance.ami = ubuntu.id
instance.instance_type = "t2.micro"
instance.credit_specification.cpu_credits = "unlimited"
instance.network_interface = [{
   "network_interface_id": iface.id,
   "device_index": 0
}]

# Based on:
# https://www.terraform.io/docs/providers/aws/r/autoscaling_group.html#mixed-instances-policy
template = aws.resource.launch_template()
template.name_prefix = "example"
template.image_id = ubuntu.id
template.instance_type = "c5.large"

group = aws.resource.autoscaling_group()
group.availability_zones = ["us-east-1a"]
group.min_size = 1
group.max_size = 1
group.desired_capacity = 1
group.mixed_instances_policy = {
   "launch_template": {
      "launch_template_specification": {
         "launch_template_id": template.id,
      },
      "override": [
        {"instance_type": "c4.large"},
        {"instance_type": "c3.large"}
      ],
   },
}

# Based on:
# https://learn.hashicorp.com/terraform/getting-started/dependencies.html#implicit-and-explicit-dependencies
bucket = aws.resource.s3_bucket()
bucket.bucket = "terraform-getting-started-guide"
bucket.acl = "private"

example = aws.resource.instance()
example.ami = "ami-2757f631"
example.instance_type = "t2.micro"
example.depends_on(bucket)

assert.eq(hcl(aws), os.read_file("fixtures/aws.tf"))
