//go:generate mockgen -destination=mock_contract_test.go -package=${GOPACKAGE} -source=contract.go
package service

type DbRepo interface {
	SetFriend(peer_1, peer_2 string) (bool, error)
	GetPeerFollows(initiator string) ([]string, error)
}
