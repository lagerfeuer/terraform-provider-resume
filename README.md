# Terraform Provider Resume

_This template repository is built on the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework). The template repository built on the [Terraform Plugin SDK](https://github.com/hashicorp/terraform-plugin-sdk) can be found at [terraform-provider-scaffolding](https://github.com/hashicorp/terraform-provider-scaffolding). See [Which SDK Should I Use?](https://www.terraform.io/docs/plugin/which-sdk.html) in the Terraform documentation for additional information._

This repository is a basic [Terraform](https://www.terraform.io) provider, containing:

- A resource and a data source (`internal/provider/`),
- Examples (`examples/`) and generated documentation (`docs/`),
- Miscellaneous meta files.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.20

## Building The Provider

```shell
make
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

_TBD_

## Developing the Provider

### Compiling
To compile the provider, run
```shell
make install
```
This will build the provider and put the provider binary in the `$GOPATH/bin` directory.
Terraform does not yet know about our local provider, so run the following in the root of the project to create and source a [Terraform config file](https://developer.hashicorp.com/terraform/cli/config/config-file):
```shell
eval $(./setup.sh)
```
After that you may try one of the examples:
```shell
cd examples/provider
terraform plan
```

### Documentation
To generate or update documentation, run `go generate`.

### Testing
In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```
