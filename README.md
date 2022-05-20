# terraform-provider-defectdojo

[DefectDojo API Terraform Provider](https://registry.terraform.io/providers/doximity/defectdojo)

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.17

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

You can configure the provider via environment variables:
```
$ export DEFECTDOJO_BASEURL="https://demo.defectdojo.org"
$ export DEFECTDOJO_APIKEY="my-api-key"
```

Or with a username/password:

```
$ export DEFECTDOJO_BASEURL="https://demo.defectdojo.org"
$ export DEFECTDOJO_USERNAME="admin"
$ export DEFECTDOJO_PASSWORD="ebgngrguvegrra"
```

Or in the terraform configuration:

```hcl
provider "defectdojo" {
  base_url = "https://defectdojo.my-company.com"
  api_key = var.dd_api_key # don't put your key in the code!
}
```

```hcl
provider "defectdojo" {
  base_url = "https://defectdojo.my-company.com"
  username = "admin"
  password = var.dd_password # don't put your password in the code!
}
```

Start using resources:

```
data "defectdojo_product_type" "this" {
  name     = "My Product Type"
}

resource "defectdojo_product" "this" {
  name            = var.product_name
  description     = "This product represents is named `${var.product_name}`"
  product_type_id = data.defectdojo_product_type.this.id
}
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```

Run one test at a time:

```shell
TESTARGS="-run TestFunctionName" make testacc
```

## Releasing a new version

Merge your changes to `master` and then push a version tag to master, like:

```
$ git checkout master
$ git pull
$ git tag v0.0.1
$ git push --tags
```

## Contributing

Pull requests are welcome. By contributing to this repository you are agreeing to the [Contributor License Agreement (CONTRIBUTING.md)](./CONTRIBUTING.md)

## Licencse

Licensed under the Apache v2 license. See [LICENSE.txt](./LICENSE.txt)
