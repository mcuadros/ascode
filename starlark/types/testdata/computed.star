load("assert.star", "assert")

aws = provider("aws", "2.13.0")

# compute of scalar
web = aws.resource.instance()
web.ami = aws.data.ami().id
assert.eq(type(web.ami), "Computed")
assert.eq(str(web.ami), '"${data.aws_ami.id_3.id}"')

# compute of set
table = aws.data.dynamodb_table()
assert.eq(str(table.ttl), '"${data.aws_dynamodb_table.id_4.ttl}"')
assert.eq(str(table.ttl[0]), '"${data.aws_dynamodb_table.id_4.ttl.0}"')
assert.eq(str(table.ttl[0].attribute_name), '"${data.aws_dynamodb_table.id_4.ttl.0.attribute_name}"')

# compute of list
instance = aws.data.instance()
assert.eq(str(instance.credit_specification), '"${data.aws_instance.id_5.credit_specification}"')
assert.eq(str(instance.credit_specification[0]), '"${data.aws_instance.id_5.credit_specification.0}"')
assert.eq(str(instance.credit_specification[0].cpu_credits), '"${data.aws_instance.id_5.credit_specification.0.cpu_credits}"')

# compute of map
computed = str(aws.resource.instance().root_block_device.volume_size)
assert.eq(computed, '"${resource.aws_instance.id_6.root_block_device.volume_size}"')