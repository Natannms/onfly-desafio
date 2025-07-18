package pedido_app_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	pedido_app "onfly-api/internal/application/pedido"
	dominio "onfly-api/internal/domain/pedido"
	infra "onfly-api/internal/infrasctructure/persistence"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	err := godotenv.Load("../../../.env.test")
	if err != nil {
		panic("Erro ao carregar .env: " + err.Error())
	}

	dsn := os.Getenv("DATABASE_URL")
	fmt.Println("DATA BASE URL", dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Erro ao conectar no banco de testes: " + err.Error())
	}

	err = db.AutoMigrate(
		&infra.Pedido{},
		&infra.UsuarioModel{},
	)
	if err != nil {
		panic("Erro ao migrar schema de teste: " + err.Error())
	}

	testDB = db

	code := m.Run()
	os.Exit(code)
}

func TestCriarPedidoDeViagem(t *testing.T) {
	repo := infra.NewPedidoRepository(testDB)
	service := pedido_app.NewPedidoService(repo)

	solicitanteID := uuid.New()
	empresaID := uuid.New()

	destino := dominio.Destino{
		Cidade: "Recife",
		Estado: "PE",
		Pais:   "Brasil",
	}

	periodo := dominio.PeriodoViagem{
		Ida:   time.Now().Add(24 * time.Hour),
		Volta: time.Now().Add(72 * time.Hour),
	}

	novoPedido, err := service.CriarPedido(solicitanteID, empresaID, destino, periodo)

	assert.NoError(t, err)
	assert.NotNil(t, novoPedido)
	assert.Equal(t, destino, novoPedido.Destino)
	assert.Equal(t, dominio.StatusSolicitado, novoPedido.Status)

	encontrado, err := repo.BuscarPorID(novoPedido.ID)
	assert.NoError(t, err)
	assert.NotNil(t, encontrado)
	assert.Equal(t, novoPedido.ID, encontrado.ID)
	assert.Equal(t, destino.Cidade, encontrado.Destino.Cidade)
	assert.Equal(t, dominio.StatusSolicitado, encontrado.Status)
}

func TestAprovarPedidoComSucesso(t *testing.T) {
	repo := infra.NewPedidoRepository(testDB)
	service := pedido_app.NewPedidoService(repo)

	solicitanteID := uuid.New()
	aprovadorID := uuid.New()
	empresaID := uuid.New()

	p, err := service.CriarPedido(solicitanteID, empresaID, dominio.Destino{
		Cidade: "Florianópolis", Estado: "SC", Pais: "Brasil",
	}, dominio.PeriodoViagem{
		Ida: time.Now().Add(48 * time.Hour), Volta: time.Now().Add(96 * time.Hour),
	})
	assert.NoError(t, err)

	err = service.AprovarPedido(p.ID, aprovadorID)
	assert.NoError(t, err)

	encontrado, _ := repo.BuscarPorID(p.ID)
	assert.Equal(t, dominio.StatusAprovado, encontrado.Status)
}

func TestNaoPermiteAprovadorSerSolicitante(t *testing.T) {
	repo := infra.NewPedidoRepository(testDB)
	service := pedido_app.NewPedidoService(repo)

	userID := uuid.New()
	empresaID := uuid.New()

	p, _ := service.CriarPedido(userID, empresaID, dominio.Destino{
		Cidade: "BH", Estado: "MG", Pais: "Brasil",
	}, dominio.PeriodoViagem{
		Ida: time.Now().Add(24 * time.Hour), Volta: time.Now().Add(48 * time.Hour),
	})

	err := service.AprovarPedido(p.ID, userID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "usuario não é um aprovador valido")
}

func TestNaoPermiteAprovarPedidoComStatusInvalido(t *testing.T) {
	repo := infra.NewPedidoRepository(testDB)
	service := pedido_app.NewPedidoService(repo)

	solicitanteID := uuid.New()
	aprovadorID := uuid.New()
	empresaID := uuid.New()

	p, _ := service.CriarPedido(solicitanteID, empresaID, dominio.Destino{
		Cidade: "Curitiba", Estado: "PR", Pais: "Brasil",
	}, dominio.PeriodoViagem{
		Ida: time.Now().Add(24 * time.Hour), Volta: time.Now().Add(48 * time.Hour),
	})

	p.Status = dominio.StatusAprovado
	_ = repo.Atualizar(p)

	err := service.AprovarPedido(p.ID, aprovadorID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "somente pedidos solicitados podem ser aprovados")
}

func TestCancelarPedidoComSucesso(t *testing.T) {
	repo := infra.NewPedidoRepository(testDB)
	service := pedido_app.NewPedidoService(repo)

	solicitanteID := uuid.New()
	aprovadorID := uuid.New()
	empresaID := uuid.New()

	p, _ := service.CriarPedido(solicitanteID, empresaID, dominio.Destino{
		Cidade: "Natal", Estado: "RN", Pais: "Brasil",
	}, dominio.PeriodoViagem{
		Ida: time.Now().Add(24 * time.Hour), Volta: time.Now().Add(48 * time.Hour),
	})

	_ = service.AprovarPedido(p.ID, aprovadorID)

	err := service.CancelarPedido(p.ID, aprovadorID)
	assert.NoError(t, err)

	encontrado, _ := repo.BuscarPorID(p.ID)
	assert.Equal(t, dominio.StatusCancelado, encontrado.Status)
}

func TestBuscarPedidoPorID(t *testing.T) {
	repo := infra.NewPedidoRepository(testDB)
	service := pedido_app.NewPedidoService(repo)

	solicitanteID := uuid.New()
	empresaID := uuid.New()

	p, _ := service.CriarPedido(solicitanteID, empresaID, dominio.Destino{
		Cidade: "Rio de Janeiro", Estado: "RJ", Pais: "Brasil",
	}, dominio.PeriodoViagem{
		Ida: time.Now().Add(48 * time.Hour), Volta: time.Now().Add(96 * time.Hour),
	})

	encontrado, err := service.BuscarPorID(p.ID)
	assert.NoError(t, err)
	assert.Equal(t, p.ID, encontrado.ID)
	assert.Equal(t, dominio.StatusSolicitado, encontrado.Status)
}

func TestListarPedidosPorStatus(t *testing.T) {
	repo := infra.NewPedidoRepository(testDB)
	service := pedido_app.NewPedidoService(repo)

	solicitanteID := uuid.New()
	empresaID := uuid.New()

	service.CriarPedido(solicitanteID, empresaID, dominio.Destino{Cidade: "Fortaleza", Estado: "CE", Pais: "Brasil"}, dominio.PeriodoViagem{Ida: time.Now().Add(24 * time.Hour), Volta: time.Now().Add(48 * time.Hour)})
	service.CriarPedido(solicitanteID, empresaID, dominio.Destino{Cidade: "Fortaleza", Estado: "CE", Pais: "Brasil"}, dominio.PeriodoViagem{Ida: time.Now().Add(24 * time.Hour), Volta: time.Now().Add(48 * time.Hour)})

	status := dominio.StatusSolicitado
	filtro := dominio.FiltroPedido{Status: &status}

	lista, err := service.ListarPedidos(filtro)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(lista), 2)
}

func TestListarPedidosPaginados(t *testing.T) {
	repo := infra.NewPedidoRepository(testDB)
	service := pedido_app.NewPedidoService(repo)

	solicitanteID := uuid.New()
	empresaID := uuid.New()

	uniquePrefix := uuid.New().String()

	for i := 0; i < 7; i++ {
		cidade := fmt.Sprintf("TestePaginado_%s_%d", uniquePrefix, i)

		_, err := service.CriarPedido(solicitanteID, empresaID, dominio.Destino{
			Cidade: cidade,
			Estado: "GO",
			Pais:   "Brasil",
		}, dominio.PeriodoViagem{
			Ida:   time.Now().Add(24 * time.Hour),
			Volta: time.Now().Add(48 * time.Hour),
		})
		assert.NoError(t, err)
	}

	status := dominio.StatusSolicitado
	cidadePrefix := fmt.Sprintf("TestePaginado_%s", uniquePrefix)
	baseFiltro := dominio.FiltroPedido{
		Status: &status,
		Cidade: &cidadePrefix,
	}

	pagina1, err := service.ListarPedidos(dominio.FiltroPedido{
		Status: baseFiltro.Status,
		Cidade: baseFiltro.Cidade,
		Limit:  3,
		Offset: 0,
	})
	assert.NoError(t, err)
	assert.Len(t, pagina1, 3)

	pagina2, err := service.ListarPedidos(dominio.FiltroPedido{
		Status: baseFiltro.Status,
		Cidade: baseFiltro.Cidade,
		Limit:  3,
		Offset: 3,
	})
	assert.NoError(t, err)
	assert.Len(t, pagina2, 3)

	pagina3, err := service.ListarPedidos(dominio.FiltroPedido{
		Status: baseFiltro.Status,
		Cidade: baseFiltro.Cidade,
		Limit:  3,
		Offset: 6,
	})
	assert.NoError(t, err)
	assert.Len(t, pagina3, 1)
}
