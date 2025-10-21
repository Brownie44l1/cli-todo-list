package internals

import "time"

type Model struct {
	Id        int
	Title     string
	Status    bool
	CreatedAt time.Time
}
