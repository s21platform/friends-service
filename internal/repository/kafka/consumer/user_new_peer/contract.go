package user_new_peer //nolint:revive,stylecheck

type ProdRepo interface {
	Process(email string, msgs []string) error
}

type DBRepo interface {
	GetUUIDForEmail(email []byte) ([]string, error)
}
