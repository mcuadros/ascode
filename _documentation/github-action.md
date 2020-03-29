---
title: 'GitHub Action'
weight: 40
---

AsCode Github Action allows to execute AsCode `run` command in response to a GitHub event such as updating a pull request or pushing a new commit to a specific branch.

This used in combination with the [Terraform GitHub Actions](https://www.terraform.io/docs/github-actions/getting-started.html) allows to execute the different terraform commands `init`, `plan` and `apply` inside of a [GitHub Workflow](https://help.github.com/en/actions/configuring-and-managing-workflows).

## Parameters

| Parameter | **Mandatory**/**Optional** | Description |
| --------- | -------- | ----------- |
| file | **Mandatory** | Starlark file to execute. Default value: `main.star` |
| hcl | **Mandatory** | HCL output file. Default value: `generated.tf` |

## Recommended Workflow


```yaml
name: 'Terraform & AsCode'
on:
  push:
    branches:
      - master
  pull_request:

jobs:
  terraform:
    name: 'Deploy'
    runs-on: ubuntu-latest
    env:
        TF_VERSION: latest
        TF_WORKING_DIR: .
    steps:
      - name: 'Checkout'
        uses: actions/checkout@master

      - name: 'AsCode Run'
        uses: mcuadros/ascode@gh-action

      - name: 'Terraform Init'
        uses: hashicorp/terraform-github-actions@master
        with:
          tf_actions_version: ${{ env.TF_VERSION }}
          tf_actions_subcommand: 'init'
          tf_actions_working_dir: ${{ env.TF_WORKING_DIR }}
          tf_actions_comment: true

      - name: 'Terraform Validate'
        uses: hashicorp/terraform-github-actions@master
        with:
          tf_actions_version: ${{ env.TF_VERSION }}
          tf_actions_subcommand: 'validate'
          tf_actions_working_dir: ${{ env.TF_WORKING_DIR }}
          tf_actions_comment: true

      - name: 'Terraform Plan'
        uses: hashicorp/terraform-github-actions@master
        with:
          tf_actions_version: ${{ env.TF_VERSION }}
          tf_actions_subcommand: 'plan'
          tf_actions_working_dir: ${{ env.TF_WORKING_DIR }}
          tf_actions_comment: true

      - name: 'Terraform Apply'
        uses: hashicorp/terraform-github-actions@master
        if: github.event_name == 'push'
        with:
          tf_actions_version: ${{ env.TF_VERSION }}
          tf_actions_subcommand: 'apply'
          tf_actions_working_dir: ${{ env.TF_WORKING_DIR }}
```