package usuario_app

import (
	"onfly-api/internal/domain/usuario"

	"github.com/google/uuid"
)

type Service struct {
	repo usuario.Repository
}

func NewUsuarioService(repo usuario.Repository) *Service {
	return &Service{repo}
}

func (s *Service) Criar(nome, email, senha string, tipo usuario.TipoConta) (*usuario.Usuario, error) {
	hashed, err := usuario.HashSenha(senha)
	if err != nil {
		return nil, err
	}

	u := &usuario.Usuario{
		ID:    uuid.New(),
		Nome:  nome,
		Email: email,
		Senha: hashed,
		Tipo:  tipo,
		Ativo: true,
	}
	err = s.repo.Criar(u)
	return u, err
}

func (s *Service) BuscarPorID(id string) (*usuario.Usuario, error) {
	return s.repo.BuscarPorID(id)
}

func (s *Service) BuscarPorEmail(email string) (*usuario.Usuario, error) {
	return s.repo.BuscarPorEmail(email)
}

func (s *Service) ListarTodos() ([]*usuario.Usuario, error) {
	return s.repo.ListarTodos()
}

func (s *Service) Atualizar(u *usuario.Usuario) error {
	return s.repo.Atualizar(u)
}

func (s *Service) Deletar(id string) error {
	return s.repo.Deletar(id)
}
