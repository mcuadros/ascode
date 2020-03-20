resource "null_resource" "bar" {
  triggers = {
    bar = "bar-value"
  }
}

module "moduleB" {
  source = "../moduleB"
}


output "qux" {
  value = module.moduleB.qux
}

output "bar" {
  value = null_resource.bar.triggers.bar
}