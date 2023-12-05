terraform {
  required_providers {
    auth0 = {
      source = "github.com/scastria/auth0"
    }
  }
}

provider "auth0" {
}

data "auth0_email_users" "EmailUsers" {
  email = "ssoni@greenstreet.com"
}