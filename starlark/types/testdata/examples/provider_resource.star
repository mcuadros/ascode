helm = tf.provider("helm", "1.0.0")

podinfo = helm.resource.release("podinfo")
podinfo.chart = "podinfo"
podinfo.version = "3.1.8"

print(hcl(podinfo))
# Output:
# resource "helm_release" "podinfo" {
#   provider = helm.id_1
#   chart    = "podinfo"
#   version  = "3.1.8"
# }
