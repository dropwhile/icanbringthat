package model

type ModelType interface {
	User | Event | EventItem | Earmark | UserPWReset
}
