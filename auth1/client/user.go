package client

const (
	EmailUsersPath = "users-by-email"
	Email          = "email"
)

type EmailUser struct {
	Email  string `json:"email,omitempty"`
	UserId string `json:"user_id,omitempty"`
}
