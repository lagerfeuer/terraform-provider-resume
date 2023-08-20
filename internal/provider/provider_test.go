// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"os"
	"testing"
)

const providerConfigTemplate = `
provider "resume" {
  endpoint = "%s"
  token = "%s"
}
`

var providerConfig = fmt.Sprintf(
	providerConfigTemplate,
	os.Getenv("RESUME_API_ENDPOINT"),
	os.Getenv("RESUME_API_TOKEN"),
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"resume": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("RESUME_API_ENDPOINT"); v == "" {
		t.Fatal("RESUME_API_ENDPOINT must be set for acceptance tests")
	}
	if v := os.Getenv("RESUME_API_TOKEN"); v == "" {
		t.Fatal("RESUME_API_TOKEN must be set for acceptance tests")
	}
}
