load("assert.star", "assert")

aws = tf.provider("aws", "2.13.0")

ami = aws.data.ami()

# attribute of scalar
web = aws.resource.instance()
web.ami = ami.id
assert.eq(type(web.ami), "Attribute<string>")
assert.eq(str(web.ami), '"${data.aws_ami.id_2.id}"')
assert.eq(web.ami.__resource__, ami)
assert.eq(web.ami.__type__, "string")

# attr names
assert.eq("__resource__" in dir(web.ami), True)
assert.eq("__type__" in dir(web.ami), True)

# attribute of set
table = aws.data.dynamodb_table()
assert.eq(str(table.ttl), '"${data.aws_dynamodb_table.id_4.ttl}"')
assert.eq(str(table.ttl[0]), '"${data.aws_dynamodb_table.id_4.ttl.0}"')
assert.eq(str(table.ttl[0].attribute_name), '"${data.aws_dynamodb_table.id_4.ttl.0.attribute_name}"')

# attribute of list
instance = aws.data.instance()
assert.eq(str(instance.credit_specification), '"${data.aws_instance.id_5.credit_specification}"')
assert.eq(str(instance.credit_specification[0]), '"${data.aws_instance.id_5.credit_specification.0}"')
assert.eq(str(instance.credit_specification[0].cpu_credits), '"${data.aws_instance.id_5.credit_specification.0.cpu_credits}"')

# attribute of block
attribute = str(aws.resource.instance().root_block_device.volume_size)
assert.eq(attribute, '"${aws_instance.id_6.root_block_device.0.volume_size}"')

# attribute on data source
assert.eq(str(aws.resource.instance().id), '"${aws_instance.id_7.id}"')

# attribute on resource
assert.eq(str(aws.data.ami().id), '"${data.aws_ami.id_8.id}"')

gcp = tf.provider("google", "3.13.0")

# attribute on list with MaxItem:1
cluster = gcp.resource.container_cluster("foo")
assert.eq(str(cluster.master_auth.client_certificate), '"${google_container_cluster.foo.master_auth.0.client_certificate}"')

# attr non-object
assert.fails(lambda: web.ami.foo, "Attribute<string> it's not a object")


# fn wrapping
assert.eq(str(fn("base64encode", web.ami)), '"${base64encode(data.aws_ami.id_2.id)}"')
