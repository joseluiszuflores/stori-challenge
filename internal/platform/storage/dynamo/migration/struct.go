package migration

import (
	"time"

	"github.com/google/uuid"
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
