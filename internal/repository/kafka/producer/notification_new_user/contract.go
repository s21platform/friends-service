package notification_new_user //nolint:revive,stylecheck

type Storage interface {
	UpdateUserInvite(initiator, invited string) error
}
