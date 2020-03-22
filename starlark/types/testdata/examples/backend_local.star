# new instance of a local backend
# https://www.terraform.io/docs/backends/types/local.html
b = backend("local")
b.path = "terraform.tfstate"

# it reads the state
s = b.state()

for provider in sorted(list(s)):
    print("%s:" % provider)
    for resource in sorted(list(s[provider]["resource"])):
        count = len(s[provider]["resource"][resource])
        print("    %s (%d)" % (resource, count))

# Output:
# google:
#     container_cluster (1)
#     container_node_pool (1)
# helm:
#     release (5)
# kubernetes:
#     cluster_role_binding (1)
#     namespace (1)
#     secret (1)
