package usuario_app_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	usuario_app "onfly-api/internal/application/usuario"
	"onfly-api/internal/domain/usuario"
	infra "onfly-api/internal/infrasctructure/persistence"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UsuarioServiceSuite struct {
	suite.Suite
	db      *gorm.DB
	service *usuario_app.Service
	uBase   *usuario.Usuario
	email   string
	userID  uuid.UUID
}

func TestUsuarioServiceSuite(t *testing.T) {
	suite.Run(t, new(UsuarioServiceSuite))
}

func (s *UsuarioServiceSuite) SetupSuite() {
	err := godotenv.Load("../../../.env.test")
	s.Require().NoError(err)

	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	s.Require().NoError(err)

	err = db.AutoMigrate(&infra.UsuarioModel{})
	s.Require().NoError(err)

	s.db = db
	repo := infra.NewUsuarioRepository(s.db)
	s.service = usuario_app.NewUsuarioService(repo)

	s.email = fmt.Sprintf("usuario_%s@example.com", uuid.NewString())
	s.userID = uuid.New()

	s.uBase, err = s.service.Criar("Usuário Base", s.email, "senha123", usuario.TipoAdministrador)
	s.Require().NoError(err)
}

func (s *UsuarioServiceSuite) TestBuscarPorID() {
	u, err := s.service.BuscarPorID(s.uBase.ID.String())
	s.NoError(err)
	s.Equal(s.uBase.Email, u.Email)
}

func (s *UsuarioServiceSuite) TestBuscarPorEmail() {
	u, err := s.service.BuscarPorEmail(s.uBase.Email)
	s.NoError(err)
	s.Equal(s.uBase.ID, u.ID)
}

func (s *UsuarioServiceSuite) TestListarTodos() {
	lista, err := s.service.ListarTodos()
	s.NoError(err)
	s.GreaterOrEqual(len(lista), 1)
}

func (s *UsuarioServiceSuite) TestAtualizar() {
	s.uBase.Nome = "Nome Atualizado"
	err := s.service.Atualizar(s.uBase)
	s.NoError(err)

	u, err := s.service.BuscarPorID(s.uBase.ID.String())
	s.NoError(err)
	s.Equal("Nome Atualizado", u.Nome)
}

func (s *UsuarioServiceSuite) TestDeletar() {
	err := s.service.Deletar(s.uBase.ID.String())
	s.NoError(err)

	u, err := s.service.BuscarPorID(s.uBase.ID.String())
	s.Error(err)
	s.Nil(u)
}

func (s *UsuarioServiceSuite) TestCriarUsuario() {
	email := fmt.Sprintf("usuario_%d@example.com", time.Now().UnixNano())
	u, err := s.service.Criar("Natan", email, "123456", usuario.TipoAdministrador)

	s.NoError(err)
	s.NotNil(u)
	s.Equal("Natan", u.Nome)
	s.Equal(email, u.Email)
	s.Equal(usuario.TipoAdministrador, u.Tipo)

	s.NotEqual("123456", u.Senha)
}
