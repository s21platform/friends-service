package notification_new_user //nolint:revive,stylecheck

type DBRepo interface {
	UpdateUserInvite(initiator, invited string) error
	SetFriend(peer1, peer2 string) (bool, error)
}
