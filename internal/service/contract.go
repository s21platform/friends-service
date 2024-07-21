package service

type dbRepo interface {
	SetFriend(peer_1, peer_2 string) (bool, error)
	CheckFriend(peer_1, peer_2 string) (bool, error)
}
