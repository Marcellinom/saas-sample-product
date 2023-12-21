package contracts

import "its.ac.id/base-go/pkg/auth/contracts"

type UserSessionData struct {
	Id                  string           `json:"id"`
	Name                string           `json:"name"`
	Nickname            string           `json:"nickname"`
	PreferredUsername   string           `json:"preferred_username"`
	Email               string           `json:"email"`
	EmailVerified       bool             `json:"email_verified"`
	Picture             string           `json:"picture"`
	Gender              string           `json:"gender"`
	Birthdate           string           `json:"birthdate"`
	Zoneinfo            string           `json:"zoneinfo"`
	Locale              string           `json:"locale"`
	PhoneNumber         string           `json:"phone_number"`
	PhoneNumberVerified bool             `json:"phone_number_verified"`
	ActiveRole          string           `json:"active_role"`
	Roles               []contracts.Role `json:"roles"`
}
