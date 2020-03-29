---
title: 'Getting Started'
weight: 2
---

This page explains the basics of using AsCode to define your infrastructure in Terraform. It assumes that you have already [installed](/docs/install) AsCode.

```sh
> ascode --help
Usage:
  ascode [OPTIONS] <repl | run | version>

AsCode - Terraform Alternative Syntax.

Help Options:
  -h, --help  Show this help message

Available commands:
  repl     Run as interactive shell.
  run      Run parses, resolves, and executes a Starlark file.
  version  Version prints information about this binary.
```

## The `repl` command

The `repl` command provides a handy `REPL` interface for debugging and tinkering with AsCode.

For example you can explore the API of a resource printing the list of arguments:

```sh
> ascode repl
>>> print(dir(helm.resource))
["__kind__", "__provider__", "release", "repository"]
```

Or to validate how a resource will be rendered:
```sh
> ascode repl
>>> aws = tf.provider("aws")
>>> web = aws.resource.instance("web", instance_type="t2.micro")
>>> print(hcl(web))
resource "aws_instance" "web" {
  provider      = aws.id_01E4JV722PS2WPKK7WQ2NMZY6D
  instance_type = "t2.micro"
}
```

## The `run` command

The `run` command executes a valid Starlack program. Using the `--print-hcl` and `--to-hcl`, an HCL encoded version of the `tf` object will be printed or saved to a given file, respectively.

This is the first step to deploy any infrastructure defined with AsCode, using `run` and generating a valid `.tf` file, we can use the standard Terraform tooling to deploy our infrastructure using `terraform init`, `terraform plan` and `terraform apply`.

To learn about writing Starlark programs, please refer to the [Language definition](/docs/starlark/) and the [API Reference](/docs/reference/) sections of this documentation.


### Basic Example

The goal of the example is to create, in DigitalOcean, one `s-1vcpu-1gb` instance called `web` in the `nyc2` region:

> To run this example, you need `terraform` correctly installed on your system.

```sh
> mkdir example; cd example
> echo 'do = tf.provider("digitalocean")' > main.star
> echo 'web = do.resource.droplet("web", name="web", size="s-1vcpu-1gb")' >> main.star
> echo 'web.region = "nyc2"' >> main.star
> echo 'web.image  = "ubuntu-18-04-x64"' >> main.star
```

Now we are ready to run our Starlark program and generate a valid `HCL` file:

```sh
> ascode run main.star --to-hcl main.tf
> cat main.tf
provider "digitalocean" {
  alias   = "id_01E4JXQD8HKW7XEQ7R5S8SP8AQ"
  version = "1.15.1"
}

resource "digitalocean_droplet" "web" {
  provider = digitalocean.id_01E4JXQD8HKW7XEQ7R5S8SP8AQ
  image    = "ubuntu-18-04-x64"
  name     = "web"
  region   = "nyc2"
  size     = "s-1vcpu-1gb"
}


```

And now as it's usual in terraform we can run `init`, `plan` or/and `apply`

```sh
> terraform init
...
> terraform plan
Terraform will perform the following actions:

  # digitalocean_droplet.web will be created
  + resource "digitalocean_droplet" "web" {
      + backups              = false
      + disk                 = (known after apply)
      + id                   = (known after apply)
      + image                = "ubuntu-18-04-x64"
      + ipv4_address         = (known after apply)
      + ipv4_address_private = (known after apply)
      + ipv6                 = false
      + ipv6_address         = (known after apply)
      + ipv6_address_private = (known after apply)
      + locked               = (known after apply)
      + memory               = (known after apply)
      + monitoring           = false
      + name                 = "web"
      + price_hourly         = (known after apply)
      + price_monthly        = (known after apply)
      + private_networking   = false
      + region               = "nyc2"
      + resize_disk          = true
      + size                 = "s-1vcpu-1gb"
      + status               = (known after apply)
      + urn                  = (known after apply)
      + vcpus                = (known after apply)
      + volume_ids           = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

> terraform apply
...
```

## The `version` command

The `version` command prints a report about the versions of the different
dependencies, and AsCode itself used to compile the binary.

```
> ascode version
Go Version: go1.14.1
AsCode Version: v0.0.1
AsCode Commit: 6a682e4
AsCode Build Date: 2020-03-29T12:43:52+02:00
Terraform Version: v0.12.23
Starlark Version: v0.0.0-20200306205701-8dd3e2ee1dd5
```