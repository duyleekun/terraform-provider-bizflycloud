terraform {
  required_providers {
    bizflycloud = {
      source  = "bizflycloud/bizflycloud"
    }
  }
}

provider "bizflycloud" {
    auth_method = "password"
    region_name = "HN"
    email = "username"
    password = ""
}

resource "bizflycloud_vpc_network" "test_vpc_network" {
    name = "test_create_vpc_network"
    description = "test vpc network"
    is_default = false
}