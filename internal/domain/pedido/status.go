package pedido

type Status string

const (
	StatusSolicitado Status = "solicitado"
	StatusAprovado   Status = "aprovado"
	StatusCancelado  Status = "cancelado"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusSolicitado, StatusAprovado, StatusCancelado:
		return true
	}
	return false
}
