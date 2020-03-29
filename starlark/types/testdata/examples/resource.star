# Create a new instance of the latest Ubuntu 14.04 on an
# t2.micro node with an AWS Tag naming it "HelloWorld"

aws = tf.provider("aws", "2.54.0")
aws.region = "us-west-2"

ubuntu_filter = "ubuntu/images/*/ubuntu-xenial-16.04-amd64-server-*"
canonical = "099720109477"

ami = aws.data.ami("ubuntu")
ami.most_recent = True
ami.filter(name="name", values=[ubuntu_filter])
ami.filter(name="virtualization-type", values=["hvm"])
ami.owners = [canonical]


instance = aws.resource.instance("web")
instance.instance_type = "t2.micro"
instance.ami = ami.id
instance.tags = {
    "name": "HelloWorld"
}

print(hcl(tf))
# Output:
# provider "aws" {
#   alias   = "id_1"
#   version = "2.54.0"
#   region  = "us-west-2"
# }
#
# data "aws_ami" "ubuntu" {
#   provider    = aws.id_1
#   most_recent = true
#   owners      = ["099720109477"]
# 
#   filter {
#     name   = "name"
#     values = ["ubuntu/images/*/ubuntu-xenial-16.04-amd64-server-*"]
#   }
# 
#   filter {
#     name   = "virtualization-type"
#     values = ["hvm"]
#   }
# }
# 
# resource "aws_instance" "web" {
#   provider      = aws.id_1
#   ami           = "${data.aws_ami.ubuntu.id}"
#   instance_type = "t2.micro"
#   tags          = { name = "HelloWorld" }
# }
