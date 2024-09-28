package models

// User represents a user in your system.
type User struct {
	ID        string `firestore:"id"`
	FirstName string `firestore:"firstName"`
	LastName  string `firestore:"lastName"`
	Email     string `firestore:"email"`
	Picture   string `firestore:"picture"`
}
