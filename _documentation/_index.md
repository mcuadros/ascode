---
title: 'Documentation'
weight: 1
---

## AsCode - Terraform Alternative Syntax 

**AsCode** is a tool to define infrastructure as code using the [Starlark](https://github.com/google/starlark-go/blob/master/doc/spec.md) language on top of Ter[Terraform](https://github.com/hashicorp/terraform)raform. It allows to describe your infrastructure using an expressive language in Terraform without writing a single line of HCL, meanwhile, you have the complete ecosystem of [providers](https://www.terraform.io/docs/providers/index.html) 

### Why?

Terraform is a great tool, with support for almost everything you can imagine, making it the industry leader. Terraform is based on HCL, a JSON-alike declarative language, with minimal control flow functionality. IMHO, to unleash the power of the IaC, a powerful, expressive language should be used, where basic elements like loops or functions are first-class citizens.

### What is Starlark?

> Starlark is a dialect of Python intended for use as a configuration language. A Starlark interpreter is typically embedded within a larger application, and this application may define additional domain-specific functions and data types beyond those provided by the core language. For example, Starlark is embedded within (and was originally developed for) the Bazel build tool, and Bazel's build language is based on Starlark.