terraform {
  required_providers {
    example = {
      source  = "example.com/tfp-example/example"
      version = "0.0.1"
    }
  }
}

provider "example" {
  host = "http://localhost:8080"
}

resource "example_shop_article" "example" {
  name        = "Princess Rosalea"
  description = "Child Shampoo & Conditioner"
}
