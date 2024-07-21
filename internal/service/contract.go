package service

type dbRepo interface {
	SetFriend(peer_1, peer_2 string) (bool, error)
	isRowFriendExist(peer_1, peer_2 string) (bool, error)
}
