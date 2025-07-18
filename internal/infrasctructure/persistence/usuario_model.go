package persistence

import (
	"onfly-api/internal/domain/usuario"
	"time"

	"github.com/google/uuid"
)

type UsuarioModel struct {
	ID           string `gorm:"primaryKey"`
	Nome         string
	Email        string `gorm:"uniqueIndex"`
	Senha        string
	Role         string
	CriadoEm     time.Time
	AtualizadoEm time.Time
}

func NewUsuarioModel() UsuarioModel {
	return UsuarioModel{
		ID:    "",
		Nome:  "",
		Email: "",
		Senha: "",
		Role:  "",
	}
}
func ToUsuarioModel(u *usuario.Usuario) *UsuarioModel {
	return &UsuarioModel{
		ID:    u.ID.String(),
		Nome:  u.Nome,
		Email: u.Email,
		Senha: u.Senha,
		Role:  string(u.Tipo),
	}
}

func (u *UsuarioModel) ToEntity() *usuario.Usuario {
	id, _ := uuid.Parse(u.ID)
	return &usuario.Usuario{
		ID:    id,
		Nome:  u.Nome,
		Email: u.Email,
		Senha: u.Senha,
		Tipo:  usuario.TipoConta(u.Role),
	}
}
