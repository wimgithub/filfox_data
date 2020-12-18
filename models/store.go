package model

type Store interface {
	BeginTx() (Store, error)
	Rollback() error
	CommitTx() error

	AddFilData(data []*Data) error
	GetFilFoxCount() (count int64, err error)
	GetFilFoxData(begin, end, height int64, msg, to, t string) (data []*Data, err error)
}
