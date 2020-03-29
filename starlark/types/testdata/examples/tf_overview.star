tf.backend = backend("gcs", bucket="tf-state")

aws = tf.provider("aws", "2.54.0", region="us-west-2")
aws.resource.instance("foo", instance_type="t2.micro")

print(hcl(tf))

# Output:
# terraform {
#   backend "gcs" {
#     bucket = "tf-state"
#   }
# }
#
# provider "aws" {
#   alias   = "id_1"
#   version = "2.54.0"
#   region  = "us-west-2"
# }
# 
# resource "aws_instance" "foo" {
#   provider      = aws.id_1
#   instance_type = "t2.micro"
# }
