package user_new_peer //nolint:revive,stylecheck

type ProdRepo interface {
	Process(email string, uuid string, msgs []string) error
}

type DBRepo interface {
	GetUUIDForEmail(email string) ([]string, error)
}
