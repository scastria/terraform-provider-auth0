# Auth0 Provider
The Auth0 provider is intended to be used in addition to the official Auth0 provider to fill in any missing gaps in
functionality.  The provider needs to be configured with the proper credentials before it can be used.
## Example Usage
```hcl
terraform {
  required_providers {
    auth0 = {
      source  = "scastria/auth0"
      version = "~> 0.1"
    }
  }
}

# Configure the Auth0 Provider
provider "auth0" {
  domain = "XXXX"
  client_id = "YYYY"
  client_secret = "ZZZZ"
}
```
## Argument Reference
* `domain` - **(Required, String)** The Auth0 domain hostname. Can be specified via env variable `AUTH0_DOMAIN`.
* `client_id` - **(Required, String)** The client id for obtaining a token to Auth0. Can be specified via env variable `AUTH0_CLIENT_ID`.
* `client_secret` - **(Required, String)** The client secret for obtaining a token to Auth0. Can be specified via env variable `AUTH0_CLIENT_SECRET`.
