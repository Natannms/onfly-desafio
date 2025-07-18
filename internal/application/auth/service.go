package auth

import (
	"errors"

	"onfly-api/cmd/fiber_http/jwt"
	"onfly-api/internal/domain/usuario"
	"onfly-api/internal/infrasctructure/persistence"

	"github.com/google/uuid"
)

type Service struct {
	repo usuario.Repository
}

func NewAuthService() *Service {
	db := persistence.GetDB() // você pode adaptar para injetar via DI/container se estiver usando isso
	repo := persistence.NewUsuarioRepository(db)
	return &Service{repo: repo}
}

func (s *Service) Register(nome, email, senha string) (*usuario.Usuario, error) {
	existente, _ := s.repo.BuscarPorEmail(email)
	if existente != nil {
		return nil, errors.New("usuário já existe")
	}

	hashed, err := usuario.HashSenha(senha)
	if err != nil {
		return nil, err
	}

	u := &usuario.Usuario{
		ID:    uuid.New(),
		Nome:  nome,
		Email: email,
		Senha: hashed,
		Tipo:  usuario.TipoEmpresa,
		Ativo: true,
	}

	err = s.repo.Criar(u)
	return u, err
}

func (s *Service) Login(email, senha string) (string, error) {
	u, err := s.repo.BuscarPorEmail(email)
	if err != nil || u == nil {
		return "", errors.New("credenciais inválidas")
	}

	if err := usuario.VerificarSenha(senha, u.Senha); err != nil {
		return "", errors.New("credenciais inválidas")
	}

	token, err := jwt.GerarToken(u.ID.String())
	if err != nil {
		return "", errors.New("erro ao gerar token")
	}

	return token, nil
}

func (s *Service) ResetPassword(email, novaSenha string) error {
	u, err := s.repo.BuscarPorEmail(email)
	if err != nil || u == nil {
		return errors.New("usuário não encontrado")
	}

	hashed, err := usuario.HashSenha(novaSenha)
	if err != nil {
		return err
	}

	u.Senha = hashed
	return s.repo.Atualizar(u)
}
