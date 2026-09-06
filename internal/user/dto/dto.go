package dto

// Response is the public representation of a user. Sensitive fields such as
// the password hash are never exposed here.
type Response struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
