package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

var (
	resourceName = "resume_resume.test"
)

func TestAccResumeResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{

				Config: providerConfig + `
resource "resume_resume" "test" {
	name = "Test McTester"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Test McTester"),
					resource.TestCheckNoResourceAttr(resourceName, "address"),
					resource.TestCheckNoResourceAttr(resourceName, "phone_number"),
					resource.TestCheckNoResourceAttr(resourceName, "website"),
				),
			},
			// Import state
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read
			{
				Config: providerConfig + `
resource "resume_resume" "test" {
	name = "TJ McTester"
	address = "1 Test Lane"
	phone_number = "555-555-5555"
	website = "https://test.com"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "TJ McTester"),
					resource.TestCheckResourceAttr(resourceName, "address", "1 Test Lane"),
					resource.TestCheckResourceAttr(resourceName, "phone_number", "555-555-5555"),
					resource.TestCheckResourceAttr(resourceName, "website", "https://test.com"),
				),
			},
		},
	})
}
