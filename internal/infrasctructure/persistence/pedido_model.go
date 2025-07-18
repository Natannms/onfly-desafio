package persistence

import (
	"time"

	"github.com/google/uuid"
)

type Pedido struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	SolicitanteID uuid.UUID `gorm:"type:uuid;not null"`
	EmpresaID     uuid.UUID `gorm:"type:uuid;not null"`

	DestinoCidade string
	DestinoEstado string
	DestinoPais   string

	DataIda   time.Time
	DataVolta time.Time

	Status   string
	CriadoEm time.Time
}
