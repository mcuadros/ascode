load("assert.star", "assert")

aws = provider("aws", "2.13.0")
assert.eq(aws.version, "2.13.0")

example = aws.resource.instance("example")
example.instance_type = "t2.micro"
example.ami = "ami-abc123"

print(aws.resource.instance.__json__)