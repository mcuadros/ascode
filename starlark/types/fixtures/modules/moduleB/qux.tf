resource "null_resource" "qux" {
  triggers = {
    qux = "qux-value"
  }
}

output "qux" {
  value = null_resource.qux.triggers.qux
}