load("assert.star", "assert")

aws = provider("aws", "2.13.0")

# compute of scalar
web = aws.resource.instance()
web.ami = aws.data.ami().id
assert.eq(type(web.ami), "computed")
assert.eq(str(web.ami), '"${data.aws_ami.8731.id}"')

# compute of set
table = aws.data.dynamodb_table()
assert.eq(str(table.ttl), '"${data.aws_dynamodb_table.8731.ttl}"')
assert.eq(str(table.ttl[0]), '"${data.aws_dynamodb_table.8731.ttl.0}"')
assert.eq(str(table.ttl[0].attribute_name), '"${data.aws_dynamodb_table.8731.ttl.0.attribute_name}"')

# compute of list
instance = aws.data.instance()
assert.eq(str(instance.credit_specification), '"${data.aws_instance.8731.credit_specification}"')
assert.eq(str(instance.credit_specification[0]), '"${data.aws_instance.8731.credit_specification.0}"')
assert.eq(str(instance.credit_specification[0].cpu_credits), '"${data.aws_instance.8731.credit_specification.0.cpu_credits}"')

# compute of map
assert.eq(str(aws.resource.instance().root_block_device.volume_size), '"${resource.aws_instance.8731.root_block_device.volume_size}"')