package usuario

type Repository interface {
	Criar(u *Usuario) error
	BuscarPorEmail(email string) (*Usuario, error)
	BuscarPorID(id string) (*Usuario, error)
	ListarTodos() ([]*Usuario, error)
	Atualizar(u *Usuario) error
	Deletar(id string) error
}
