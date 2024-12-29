package todosRepo

import "time"

type Todo struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	IsCompleted bool      `json:"isCompleted"`
	CreatedAt   time.Time `json:"createdAt"`
}
