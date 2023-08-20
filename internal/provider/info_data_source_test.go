package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"regexp"
	"testing"
)

func TestAccInfoDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "resume_info" "this" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.resume_info.this", "name", "Resume API",
					),
					resource.TestMatchResourceAttr(
						"data.resume_info.this", "version", regexp.MustCompile("^v?([0-9]+.){2}[0-9]$"),
					),
					resource.TestMatchResourceAttr(
						"data.resume_info.this",
						"environment",
						regexp.MustCompile("(development|production)"),
					),
				),
			},
		},
	})
}
