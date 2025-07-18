package pedido

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type PedidoDeViagem struct {
	ID            uuid.UUID
	SolicitanteID uuid.UUID
	EmpresaID     uuid.UUID
	Destino       Destino
	Periodo       PeriodoViagem
	Status        Status
	CriadoEm      time.Time
}

type FiltroPedido struct {
	Status        *Status
	Destino       *string
	Inicio        *time.Time
	Fim           *time.Time
	SolicitanteID *uuid.UUID
	Limit         int
	Offset        int
	Cidade        *string
}

func NovoPedidoDeViagem(solicitanteID uuid.UUID, empresaID uuid.UUID, destino Destino, periodo PeriodoViagem) (*PedidoDeViagem, error) {
	if !destino.Valido() {
		return nil, errors.New("destino de viagem inválido")
	}
	if !periodo.Valido() {
		return nil, errors.New("período de viagem inválido")
	}

	return &PedidoDeViagem{
		ID:            uuid.New(),
		SolicitanteID: solicitanteID,
		EmpresaID:     empresaID,
		Destino:       destino,
		Periodo:       periodo,
		Status:        StatusSolicitado,
		CriadoEm:      time.Now(),
	}, nil
}

func (p *PedidoDeViagem) AprovarValido(usuarioID uuid.UUID) error {
	if usuarioID == p.SolicitanteID {
		return errors.New("usuario não é um aprovador valido de solicitação")
	}

	return nil
}
func (p *PedidoDeViagem) AprovarPor(usuarioID uuid.UUID) error {

	if p.Status != StatusSolicitado {
		return errors.New("somente pedidos solicitados podem ser aprovados")
	}
	p.Status = StatusAprovado
	return nil
}

func (p *PedidoDeViagem) CancelarPor(usuarioID uuid.UUID) error {
	if p.Status == StatusCancelado {
		return errors.New("pedido já está cancelado")
	}
	p.Status = StatusCancelado
	return nil
}
