package service

type DbRepo interface {
	SetFriend(peer_1, peer_2 string) (bool, error)
	GetPeerFollows(initiator string) ([]string, error)
}
