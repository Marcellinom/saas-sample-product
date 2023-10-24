package contracts

import "its.ac.id/base-go/pkg/auth/contracts"

type UserSessionData struct {
	Id                string           `json:"id"`
	Name              string           `json:"name"`
	PreferredUsername string           `json:"preferred_username"`
	Email             string           `json:"email"`
	Picture           string           `json:"picture"`
	ActiveRole        string           `json:"active_role"`
	Roles             []contracts.Role `json:"roles"`
}
