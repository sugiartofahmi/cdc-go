package enums

type UserStatusEnum int

const (
	UserStatusActive    UserStatusEnum = 1
	UserStatusInactive  UserStatusEnum = 2
	UserStatusSuspended UserStatusEnum = 3
)

func (e UserStatusEnum) GetLabel() string {
	switch e {
	case UserStatusActive:
		return "Active"
	case UserStatusInactive:
		return "Inactive"
	case UserStatusSuspended:
		return "Suspended"
	default:
		return "Unknown"
	}
}
