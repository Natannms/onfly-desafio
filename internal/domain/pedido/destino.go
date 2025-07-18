package pedido

type Destino struct {
	Cidade string
	Estado string
	Pais   string
}

func (d Destino) Valido() bool {
	return d.Cidade != "" && d.Pais != "" && d.Estado != ""
}
