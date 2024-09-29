package models

type Flashcard struct {
	Word        string `firestore:"word"`
	Proficiency *int   `firestore:"proficiency"`
}
