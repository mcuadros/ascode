provider "aws" {
  alias   = "id_1"
  version = "2.13.0"
  region  = "us-west-2"
}

data "aws_ami" "id_5" {
  provider    = aws.id_1
  most_recent = true
  owners      = ["099720109477"]

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-trusty-14.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

resource "aws_autoscaling_group" "id_8" {
  provider           = aws.id_1
  availability_zones = ["us-east-1a"]
  desired_capacity   = 1
  max_size           = 1
  min_size           = 1

  mixed_instances_policy {
    launch_template {
      launch_template_specification {
        launch_template_id = "${aws_launch_template.id_7.id}"
      }

      override {
        instance_type = "c4.large"
      }

      override {
        instance_type = "c3.large"
      }
    }
  }
}

resource "aws_instance" "id_6" {
  provider      = aws.id_1
  ami           = "${data.aws_ami.id_5.id}"
  instance_type = "t2.micro"

  credit_specification {
    cpu_credits = "unlimited"
  }

  network_interface {
    device_index         = 0
    network_interface_id = "${aws_network_interface.id_4.id}"
  }
}

resource "aws_instance" "id_10" {
  provider      = aws.id_1
  ami           = "ami-2757f631"
  instance_type = "t2.micro"
  depends_on    = [aws_s3_bucket.id_9]
}

resource "aws_launch_template" "id_7" {
  provider      = aws.id_1
  image_id      = "${data.aws_ami.id_5.id}"
  instance_type = "c5.large"
  name_prefix   = "example"
}

resource "aws_network_interface" "id_4" {
  provider    = aws.id_1
  private_ips = ["172.16.10.100"]
  subnet_id   = "${aws_subnet.id_3.id}"
  tags        = { Name = "primary_network_iterface" }
}

resource "aws_s3_bucket" "id_9" {
  provider = aws.id_1
  acl      = "private"
  bucket   = "terraform-getting-started-guide"
}

resource "aws_subnet" "id_3" {
  provider          = aws.id_1
  availability_zone = "us-west-2a"
  cidr_block        = "172.16.0.0/24"
  tags              = { Name = "tf-example" }
  vpc_id            = "${aws_vpc.id_2.id}"
}

resource "aws_vpc" "id_2" {
  provider   = aws.id_1
  cidr_block = "172.16.0.0/16"
  tags       = { Name = "tf-example" }
}

