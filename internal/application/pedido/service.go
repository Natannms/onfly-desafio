package pedido_app

import (
	"errors"
	"onfly-api/internal/domain/pedido"
	"onfly-api/internal/infrasctructure/persistence"

	"github.com/google/uuid"
)

type PedidoService struct {
	repo persistence.Repository
}

func NewPedidoService(repo persistence.Repository) *PedidoService {
	return &PedidoService{repo: repo}
}

func (s *PedidoService) CriarPedido(
	solicitanteID uuid.UUID,
	empresaID uuid.UUID,
	destino pedido.Destino,
	periodo pedido.PeriodoViagem,
) (*pedido.PedidoDeViagem, error) {
	p, err := pedido.NovoPedidoDeViagem(solicitanteID, empresaID, destino, periodo)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Salvar(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *PedidoService) AprovarPedido(pedidoID, aprovadorID uuid.UUID) error {
	p, err := s.repo.BuscarPorID(pedidoID)
	if err != nil {
		return err
	}
	if p == nil {
		return errors.New("pedido não encontrado")
	}

	if err := p.AprovarValido(aprovadorID); err != nil {
		return err
	}

	if err := p.AprovarPor(aprovadorID); err != nil {
		return err
	}

	return s.repo.Atualizar(p)
}

func (s *PedidoService) CancelarPedido(pedidoID, usuarioID uuid.UUID) error {
	p, err := s.repo.BuscarPorID(pedidoID)
	if err != nil {
		return err
	}

	if p == nil {
		return errors.New("pedido não encontrado")
	}

	if err := p.CancelarPor(usuarioID); err != nil {
		return err
	}

	return s.repo.Atualizar(p)
}

func (s *PedidoService) BuscarPorID(id uuid.UUID) (*pedido.PedidoDeViagem, error) {
	return s.repo.BuscarPorID(id)
}

func (s *PedidoService) ListarPedidos(filtro pedido.FiltroPedido) ([]pedido.PedidoDeViagem, error) {
	return s.repo.ListarPorFiltro(filtro)
}
