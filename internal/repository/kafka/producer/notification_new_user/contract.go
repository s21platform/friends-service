package notification_new_user //nolint:revive,stylecheck

type DBRepo interface {
	UpdateUserInvite(initiator, invited string) error
}
