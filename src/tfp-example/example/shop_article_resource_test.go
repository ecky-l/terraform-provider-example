package example

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func Test_ShopArticle(t *testing.T) {
	name := "Princess Rosalea"
	description := "Child Shampoo & Conditioner"

	resourceTf := `
provider "example" {
  host = "http://localhost:8080"
}

resource "example_shop_article" "shampoo" {
    name = "` + name + `"
    description = "` + description + `"
}
`

	shampooResource := "example_shop_article.shampoo"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceTf,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						shampooResource, "name", name),
					resource.TestCheckResourceAttr(
						shampooResource, "description", description),
				),
			},
		},
	})
}
