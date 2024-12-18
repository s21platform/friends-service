//go:generate mockgen -destination=mock_contract_test.go -package=${GOPACKAGE} -source=contract.go
package service

type DBRepo interface {
	SetFriend(peer1, peer2 string) (bool, error)
	GetPeerFollows(initiator string) ([]string, error)
	GetWhoFollowsPeer(initiator string) ([]string, error)
	RemoveSubscribe(peer1, peer2 string) error
	SetInvitePeer(initiator, email string) error
	RemoveFriends(peer1, peer2 string) (bool, error)
	GetCountFriends(uuid string) (int64, int64, error)
	IsRowFriendExist(peer1, peer2 string) (bool, error)
}
