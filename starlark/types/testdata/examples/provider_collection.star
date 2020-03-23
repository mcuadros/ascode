tf.provider("aws", "2.13.0", "qux")
tf.provider("aws", "2.13.0", "bar")
tf.provider("google")

# providers can be access by indexing
aws_names = tf.provider["aws"].keys()
print("aws providers:", sorted(aws_names))

# or by the get method
google_names = tf.provider.get("google").keys()
print("google providers:", google_names)

# Output:
# aws providers: ["bar", "qux"]
# google providers: ["id_1"]