package migration

import (
	"github.com/google/uuid"
	"time"
)

type Transaction struct {
	ID          int
	Date        time.Time
	Transaction string
}

type User struct {
	ID    uuid.UUID
	Name  string
	Email string
}
