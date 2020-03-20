module "moduleA" {
  source = "./moduleA"
}

resource "null_resource" "foo" {
  triggers = {
    foo = "foo-value"
    bar = module.moduleA.bar
    qux = module.moduleA.qux
  }
}