package models

// User represents a user in your system.
type User struct {
	Email    string `firestore:"email"`
	Picture  string `firestore:"picture"`
	Username string `firestore:"username"`
}
