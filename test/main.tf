terraform {
  required_providers {
    auth1 = {
      source = "github.com/scastria/auth1"
    }
  }
}

provider "auth1" {
}

data "auth1_email_users" "EmailUsers" {
  email = "ssoni@greenstreet.com"
}