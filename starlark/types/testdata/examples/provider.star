def print_provider_info(p):
    print("Provider %s[%s] (%s)" % (p.__type__, p.__name__, p.__version__))
    print("  Defines Data Sources: %d" % len(dir(p.data)))
    print("  Defines Resources: %d" % len(dir(p.resource)))
    print("  Configuration: %s" % p.__dict__)

provider = tf.provider("google", "3.13.0")
provider.project = "acme-app"
provider.region = "us-central1"

print_provider_info(provider)
# Output:
# Provider google[id_1] (3.13.0)
#   Defines Data Sources: 60
#   Defines Resources: 263
#   Configuration: {"project": "acme-app", "region": "us-central1"}

