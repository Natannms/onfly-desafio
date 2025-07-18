package persistence

import (
	"errors"
	"fmt"
	"onfly-api/internal/domain/pedido"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Salvar(p *pedido.PedidoDeViagem) error
	BuscarPorID(id uuid.UUID) (*pedido.PedidoDeViagem, error)
	ListarPorFiltro(filtro pedido.FiltroPedido) ([]pedido.PedidoDeViagem, error)
	Atualizar(p *pedido.PedidoDeViagem) error
}

type PedidoRepository struct {
	db *gorm.DB
}

func toModel(p *pedido.PedidoDeViagem) Pedido {
	return Pedido{
		ID:            p.ID,
		SolicitanteID: p.SolicitanteID,
		EmpresaID:     p.EmpresaID,
		DestinoCidade: p.Destino.Cidade,
		DestinoEstado: p.Destino.Estado,
		DestinoPais:   p.Destino.Pais,
		DataIda:       p.Periodo.Ida,
		DataVolta:     p.Periodo.Volta,
		Status:        string(p.Status),
		CriadoEm:      p.CriadoEm,
	}
}

func toEntity(m *Pedido) *pedido.PedidoDeViagem {
	return &pedido.PedidoDeViagem{
		ID:            m.ID,
		SolicitanteID: m.SolicitanteID,
		EmpresaID:     m.EmpresaID,
		Destino: pedido.Destino{
			Cidade: m.DestinoCidade,
			Estado: m.DestinoEstado,
			Pais:   m.DestinoPais,
		},
		Periodo: pedido.PeriodoViagem{
			Ida:   m.DataIda,
			Volta: m.DataVolta,
		},
		Status:   pedido.Status(m.Status),
		CriadoEm: m.CriadoEm,
	}
}

func NewPedidoRepository(db *gorm.DB) *PedidoRepository {
	return &PedidoRepository{db}
}

func (r *PedidoRepository) Salvar(p *pedido.PedidoDeViagem) error {
	model := toModel(p)
	return r.db.Create(&model).Error
}

func (r *PedidoRepository) Atualizar(p *pedido.PedidoDeViagem) error {
	model := toModel(p)
	fmt.Println("atualizando", p)
	fmt.Println("atualizando", model)
	return r.db.Save(&model).Error
}

func (r *PedidoRepository) BuscarPorID(id uuid.UUID) (*pedido.PedidoDeViagem, error) {
	var model Pedido
	err := r.db.First(&model, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toEntity(&model), nil
}

func (r *PedidoRepository) ListarPorFiltro(f pedido.FiltroPedido) ([]pedido.PedidoDeViagem, error) {
	var models []Pedido
	tx := r.db

	if f.Status != nil {
		tx = tx.Where("status = ?", string(*f.Status))
	}
	if f.Cidade != nil {
		tx = tx.Where("destino_cidade ILIKE ?", "%"+*f.Cidade+"%")
	}
	if f.Cidade != nil {
		tx = tx.Where("destino_cidade ILIKE ?", "%"+*f.Cidade+"%")
	}
	if f.SolicitanteID != nil {
		tx = tx.Where("solicitante_id = ?", *f.SolicitanteID)
	}
	if f.Inicio != nil {
		tx = tx.Where("data_ida >= ?", *f.Inicio)
	}
	if f.Fim != nil {
		tx = tx.Where("data_volta <= ?", *f.Fim)
	}

	if f.Limit > 0 {
		tx = tx.Limit(f.Limit)
	} else {
		tx = tx.Limit(20)
	}

	tx = tx.Offset(f.Offset)

	if err := tx.Find(&models).Error; err != nil {
		return nil, err
	}

	var pedidos []pedido.PedidoDeViagem
	for _, m := range models {
		pedidos = append(pedidos, *toEntity(&m))
	}

	return pedidos, nil
}
