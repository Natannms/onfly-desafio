package persistence

import (
	"onfly-api/internal/domain/usuario"

	"gorm.io/gorm"
)

type UsuarioRepository struct {
	db *gorm.DB
}

func NewUsuarioRepository(db *gorm.DB) *UsuarioRepository {
	return &UsuarioRepository{db}
}

func (r *UsuarioRepository) Criar(usuario *usuario.Usuario) error {
	model := ToUsuarioModel(usuario)

	if err := r.db.Create(&model).Error; err != nil {
		return err
	}
	*usuario = *model.ToEntity()
	return nil
}

func (r *UsuarioRepository) BuscarPorID(id string) (*usuario.Usuario, error) {
	var model UsuarioModel
	if err := r.db.First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return model.ToEntity(), nil
}

func (r *UsuarioRepository) BuscarPorEmail(email string) (*usuario.Usuario, error) {
	var model UsuarioModel
	if err := r.db.First(&model, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return model.ToEntity(), nil
}

func (r *UsuarioRepository) ListarTodos() ([]*usuario.Usuario, error) {
	var models []UsuarioModel
	if err := r.db.Find(&models).Error; err != nil {
		return nil, err
	}

	var usuarios []*usuario.Usuario
	for _, m := range models {
		usuarios = append(usuarios, m.ToEntity())
	}
	return usuarios, nil
}

func (r *UsuarioRepository) Atualizar(usuario *usuario.Usuario) error {
	model := ToUsuarioModel(usuario)
	return r.db.Save(&model).Error
}

func (r *UsuarioRepository) Deletar(id string) error {
	return r.db.Delete(&UsuarioModel{}, "id = ?", id).Error
}
