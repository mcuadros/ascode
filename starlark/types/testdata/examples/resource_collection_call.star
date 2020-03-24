aws = tf.provider("aws")

# Resources can be defined in different ways...

# you can set the resource attributes...
v = aws.resource.instance()
v.instance_type = "t2.micro"
v.tags = {"name": "HelloWorld"}

# or using a dict in the constructor...
d = aws.resource.instance({
    "instance_type": "t2.micro",
    "tags": {"name": "HelloWorld"},
})

# or even using kwargs
k = aws.resource.instance(instance_type="t2.micro", tags={"name": "HelloWorld"})

# and all the resources are equivalent:
print(v.__dict__)
print(d.__dict__)
print(k.__dict__)

# Output:
# {"instance_type": "t2.micro", "tags": {"name": "HelloWorld"}}
# {"instance_type": "t2.micro", "tags": {"name": "HelloWorld"}}
# {"instance_type": "t2.micro", "tags": {"name": "HelloWorld"}}
