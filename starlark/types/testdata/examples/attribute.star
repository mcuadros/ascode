# When a Resource has an Attribute means that the value it's only available
# during the `apply` phase of Terraform. So in, AsCode an attribute behaves
# like a poor-man pointer.

aws = tf.provider("aws")

ami = aws.resource.ami("ubuntu")

instance = aws.resource.instance("foo")
instance.ami = ami.id

print(instance.ami)
# Output:
# ${aws_ami.ubuntu.id}