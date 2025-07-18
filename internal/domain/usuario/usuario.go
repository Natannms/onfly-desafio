package usuario

import "github.com/google/uuid"

type TipoConta string

const (
	TipoAdministrador TipoConta = "admin"
	TipoEmpresa       TipoConta = "empresa"
)

type Usuario struct {
	ID    uuid.UUID
	Nome  string
	Email string
	Senha string
	Tipo  TipoConta
	Ativo bool
}
