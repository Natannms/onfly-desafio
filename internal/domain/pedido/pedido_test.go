package pedido

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNovoPedidoDeViagem(t *testing.T) {
	solicitanteID := uuid.New()
	empresaID := uuid.New()

	dest := Destino{
		Cidade: "Belo Horizonte",
		Estado: "Minas Gerais",
		Pais:   "Brasil",
	}

	periodo := PeriodoViagem{
		Ida:   time.Now().Add(24 * time.Hour),
		Volta: time.Now().Add(48 * time.Hour),
	}

	p, err := NovoPedidoDeViagem(solicitanteID, empresaID, dest, periodo)

	assert.NoError(t, err)
	assert.NotNil(t, p)
	assert.Equal(t, StatusSolicitado, p.Status)
}

func TestPedidoNaoPermiteDestinoInvalido(t *testing.T) {
	solicitanteID := uuid.New()
	empresaID := uuid.New()

	dest := Destino{
		Cidade: "",
		Estado: "",
		Pais:   "",
	}

	periodo := PeriodoViagem{
		Ida:   time.Now().Add(24 * time.Hour),
		Volta: time.Now().Add(48 * time.Hour),
	}

	p, err := NovoPedidoDeViagem(solicitanteID, empresaID, dest, periodo)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "destino de viagem inválido")
	assert.Nil(t, p)
}
func TestPedidoNaoPermitePeriodoInvalido(t *testing.T) {
	solicitanteID := uuid.New()
	empresaID := uuid.New()

	dest := Destino{
		Cidade: "Belzonte",
		Estado: "Minas Gerais",
		Pais:   "Brasil",
	}

	periodo := PeriodoViagem{
		Ida:   time.Now().Add(48 * time.Hour),
		Volta: time.Now().Add(24 * time.Hour),
	}

	p, err := NovoPedidoDeViagem(solicitanteID, empresaID, dest, periodo)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "período de viagem inválido")
	assert.Nil(t, p)
}

func TestSolicitanteNaoPodeAprovarPedido(t *testing.T) {
	solicitanteID := uuid.New()
	empresaID := uuid.New()

	dest := Destino{
		Cidade: "Belo Horizonte",
		Estado: "Minas Gerais",
		Pais:   "Brasil",
	}

	periodo := PeriodoViagem{
		Ida:   time.Now().Add(24 * time.Hour),
		Volta: time.Now().Add(48 * time.Hour),
	}

	p, _ := NovoPedidoDeViagem(solicitanteID, empresaID, dest, periodo)
	isvalidApprv := p.AprovarValido(solicitanteID)
	assert.EqualError(t, isvalidApprv, "usuario não é um aprovador valido de solicitação")
	assert.Equal(t, StatusSolicitado, p.Status)
}

func TestSomentePedidosSolicitadosPodemSerAprovados(t *testing.T) {
	solicitanteID := uuid.New()
	empresaID := uuid.New()
	aprovador := uuid.New()
	dest := Destino{
		Cidade: "Belo Horizonte",
		Estado: "Minas Gerais",
		Pais:   "Brasil",
	}

	periodo := PeriodoViagem{
		Ida:   time.Now().Add(24 * time.Hour),
		Volta: time.Now().Add(48 * time.Hour),
	}

	p, _ := NovoPedidoDeViagem(solicitanteID, empresaID, dest, periodo)

	p.Status = StatusAprovado // alterando status da solicitação
	errAprv := p.AprovarPor(aprovador)

	assert.Error(t, errAprv)
	assert.Equal(t, errAprv.Error(), "somente pedidos solicitados podem ser aprovados")
}

func TestValidaApovador(t *testing.T) {
	solicitanteID := uuid.New()
	empresaID := uuid.New()
	aprovador := uuid.New()
	dest := Destino{
		Cidade: "Belo Horizonte",
		Estado: "Minas Gerais",
		Pais:   "Brasil",
	}

	periodo := PeriodoViagem{
		Ida:   time.Now().Add(24 * time.Hour),
		Volta: time.Now().Add(48 * time.Hour),
	}

	p, _ := NovoPedidoDeViagem(solicitanteID, empresaID, dest, periodo)
	isvalidApprv := p.AprovarValido(aprovador)

	assert.NoError(t, isvalidApprv)
}

func TestAprovadorInvalido(t *testing.T) {
	solicitanteID := uuid.New()
	empresaID := uuid.New()
	aprovador := solicitanteID
	dest := Destino{
		Cidade: "Belo Horizonte",
		Estado: "Minas Gerais",
		Pais:   "Brasil",
	}

	periodo := PeriodoViagem{
		Ida:   time.Now().Add(24 * time.Hour),
		Volta: time.Now().Add(48 * time.Hour),
	}

	p, _ := NovoPedidoDeViagem(solicitanteID, empresaID, dest, periodo)
	isvalidApprv := p.AprovarValido(aprovador)

	assert.Error(t, isvalidApprv)
	assert.Equal(t, isvalidApprv.Error(), "usuario não é um aprovador valido de solicitação")
}
func TestSolicitacaoAprovada(t *testing.T) {
	solicitanteID := uuid.New()
	empresaID := uuid.New()
	aprovador := uuid.New()
	dest := Destino{
		Cidade: "Belo Horizonte",
		Estado: "Minas Gerais",
		Pais:   "Brasil",
	}

	periodo := PeriodoViagem{
		Ida:   time.Now().Add(24 * time.Hour),
		Volta: time.Now().Add(48 * time.Hour),
	}

	p, _ := NovoPedidoDeViagem(solicitanteID, empresaID, dest, periodo)
	errAprv := p.AprovarPor(aprovador)
	assert.NoError(t, errAprv)
}
func TestSolicitacaoCancelada(t *testing.T) {
	solicitanteID := uuid.New()
	empresaID := uuid.New()
	aprovador := uuid.New()
	dest := Destino{
		Cidade: "Belo Horizonte",
		Estado: "Minas Gerais",
		Pais:   "Brasil",
	}

	periodo := PeriodoViagem{
		Ida:   time.Now().Add(24 * time.Hour),
		Volta: time.Now().Add(48 * time.Hour),
	}

	p, _ := NovoPedidoDeViagem(solicitanteID, empresaID, dest, periodo)

	p.AprovarPor(aprovador)
	errCancel := p.CancelarPor(aprovador)
	assert.NoError(t, errCancel)
}

func TestNaoPermiteCancelarSolicitacoesComStatusCancelado(t *testing.T) {
	solicitanteID := uuid.New()
	empresaID := uuid.New()
	aprovador := uuid.New()
	dest := Destino{
		Cidade: "Belo Horizonte",
		Estado: "Minas Gerais",
		Pais:   "Brasil",
	}

	periodo := PeriodoViagem{
		Ida:   time.Now().Add(24 * time.Hour),
		Volta: time.Now().Add(48 * time.Hour),
	}

	p, _ := NovoPedidoDeViagem(solicitanteID, empresaID, dest, periodo)

	p.AprovarPor(aprovador)
	p.Status = StatusCancelado
	errCancel := p.CancelarPor(aprovador)
	assert.Error(t, errCancel)
	assert.Equal(t, errCancel.Error(), "pedido já está cancelado")

}
