package storelayer

type store struct {
	dbMock bool
}

type Store interface {
}

func New() *store {
	return &store{
		dbMock: true,
	}
}
