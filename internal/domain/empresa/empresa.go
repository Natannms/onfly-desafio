package empresa

import "github.com/google/uuid"

type Empresa struct {
	ID   uuid.UUID
	Nome string
}
