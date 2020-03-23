def print_provider_stats(p):
    print("Provider %s[%s] (%s)" % (p.__type__, p.__name__, p.__version__))
    print("  Defines Data Sources: %d" % len(dir(p.data)))
    print("  Defines Resources: %d" % len(dir(p.resource)))

provider = tf.provider("aws", "2.13.0")
print_provider_stats(provider)

# Output:
# Provider aws[id_1] (2.13.0)
#   Defines Data Sources: 131
#   Defines Resources: 506
