package pedido

import "time"

type PeriodoViagem struct {
	Ida   time.Time
	Volta time.Time
}

func (p PeriodoViagem) Valido() bool {
	return p.Ida.Before(p.Volta)
}
