terraform {
  required_providers {
    random = {
      source  = "hashicorp/random"
      version = "3.6.3"
    }
  }
}

provider "random" {
  alias = "random-alias"
}

resource "random_string" "random" {

  length   = 16
  provider = random.random-alias
  lifecycle {
    create_before_destroy = false
  }

  special  = true
  
}
