aws = tf.provider("aws")
aws.resource.instance("foo", instance_type="t2.micro")
aws.resource.instance("bar", instance_type="a1.medium")
aws.resource.instance("qux", instance_type="t2.micro")

r = aws.resource.instance.search("bar") 
print("Instance type of `bar`: %s" % r[0].instance_type)

r = aws.resource.instance.search("instance_type", "t2.micro")
print("Instances with 't2.micro`: %d" % len(r))

# Output:
# Instance type of `bar`: a1.medium
# Instances with 't2.micro`: 2
