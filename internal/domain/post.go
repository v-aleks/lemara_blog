package domain

import (
	"time"

	"github.com/google/uuid"
)

// Посты

type Post struct {
	ID        	uuid.UUID 		`json:"id"`
	Title     	string 			`json:"title"`
	Content   	string 			`json:"content"`
	Author    	string 			`json:"author"`
	CreatedAt 	time.Time 		`json:"created_at"`
	UpdatedAt 	time.Time 		`json:"updated_at"`
	Tags      	[]Tag 			`json:"tags"`
}

type PostCreateRequest struct {
	Title     	string 			`json:"title"`
	Content   	string 			`json:"content"`
	Author    	string 			`json:"author"`
	Tags      	[]string 		`json:"tags"`
}


type Tag struct {
	ID   	uuid.UUID 		`json:"id"`
	Name 	string 		`json:"name"`
}
