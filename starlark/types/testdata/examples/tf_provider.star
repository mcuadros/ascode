aws = tf.provider("aws", "2.13.0", "aliased")
helm = tf.provider("helm", "1.0.0")
ignition = tf.provider("ignition")

print("terraform =", tf.version)
for type in sorted(list(tf.provider)):
    for name in tf.provider[type]:
        p = tf.provider[type][name]
        print("   %s (%s) = %s" % (type, name, p.__version__))

# Output:
# terraform = 0.12.23
#    aws (aliased) = 2.13.0
#    helm (id_1) = 1.0.0
#    ignition (id_2) = 1.2.1