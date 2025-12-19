package domain

import "time"

type Message struct {
	ID        string    `firestore:"id"`
	Content   string    `firestore:"content"`
	Author    string    `firestore:"author"`
	CreatedAt time.Time `firestore:"createdAt"`
}
