package new_friend

type DBRepo interface {
	GetUUIDForEmail(email string) ([]string, error)
}

type NotificationNewFriend interface {
	ProduceMessage(message interface{}) error
}
