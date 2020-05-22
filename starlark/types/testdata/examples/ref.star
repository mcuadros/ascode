# Non-computed arguments are passed by value to one argument to another 
# in some cases may be required to make a referece to and relay in the HCL
# interpolation. 

aws = tf.provider("aws")

instance = aws.resource.instance("foo")
instance.ami = "foo"

print(ref(instance, "ami"))
# Output:
# ${aws_instance.foo.ami}