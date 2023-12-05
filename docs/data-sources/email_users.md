# Data Source: auth0_email_users
Represents a collection of users that all share the same email address
## Example usage
```hcl
data "auth0_email_users" "example" {
  email = "hello@example.com"
}
```
## Argument Reference
* `email` - **(Required, String)** The email of the users.
## Attribute Reference
* `id` - **(String)** Same as `email`
* `user_ids` - **(List of String)** The list of user ids that share the specified email address.
