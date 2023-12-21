package contracts

import (
	"github.com/pkg/errors"
)

var (
	ErrUserDoesNotHaveRole = errors.New("user_does_not_have_role")
)

type Role struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
	IsDefault   bool     `json:"is_default"`
}

type User struct {
	id                  string
	name                string
	preferred_username  string
	email               string
	emailVerified       bool
	picture             string
	gender              string
	birthdate           string
	zoneinfo            string
	locale              string
	phoneNumber         string
	phoneNumberVerified bool
	activeRole          string
	activeRoleName      string
	roles               []Role
	hashedPassword      string
}

func NewUser(id string) *User {
	return &User{
		id:             id,
		activeRole:     "",
		activeRoleName: "",
		roles:          []Role{},
	}
}

func (u *User) Id() string {
	return u.id
}

func (u *User) Name() string {
	return u.name
}

func (u *User) SetName(name string) {
	u.name = name
}

func (u *User) PreferredUsername() string {
	return u.preferred_username
}

func (u *User) SetPreferredUsername(preferredUsername string) {
	u.preferred_username = preferredUsername
}

func (u *User) Email() string {
	return u.email
}

func (u *User) SetEmail(email string) {
	u.email = email
}

func (u *User) EmailVerified() bool {
	return u.emailVerified
}

func (u *User) SetEmailVerified(emailVerified bool) {
	u.emailVerified = emailVerified
}

func (u *User) Picture() string {
	return u.picture
}

func (u *User) SetPicture(picture string) {
	u.picture = picture
}

func (u *User) Gender() string {
	return u.gender
}

func (u *User) SetGender(gender string) {
	u.gender = gender
}

func (u *User) Birthdate() string {
	return u.birthdate
}

func (u *User) SetBirthdate(birthdate string) {
	u.birthdate = birthdate
}

func (u *User) Zoneinfo() string {
	return u.zoneinfo
}

func (u *User) SetZoneinfo(zoneinfo string) {
	u.zoneinfo = zoneinfo
}

func (u *User) Locale() string {
	return u.locale
}

func (u *User) SetLocale(locale string) {
	u.locale = locale
}

func (u *User) PhoneNumber() string {
	return u.phoneNumber
}

func (u *User) SetPhoneNumber(phoneNumber string) {
	u.phoneNumber = phoneNumber
}

func (u *User) PhoneNumberVerified() bool {
	return u.phoneNumberVerified
}

func (u *User) SetPhoneNumberVerified(phoneNumberVerified bool) {
	u.phoneNumberVerified = phoneNumberVerified
}

func (u *User) ActiveRole() string {
	return u.activeRole
}

func (u *User) ActiveRoleName() string {
	return u.activeRoleName
}

func (u *User) Roles() []Role {
	return u.roles
}

func (u *User) AddRole(id string, name string, permissions []string, isDefault bool) {
	u.roles = append(u.roles, Role{
		Id:          id,
		Name:        name,
		Permissions: permissions,
		IsDefault:   isDefault,
	})

	if isDefault || u.activeRole == "" {
		u.SetActiveRole(id)
	}
}

func (u *User) SetActiveRole(id string) error {
	for _, role := range u.roles {
		if role.Id == id {
			u.activeRole = id
			u.activeRoleName = role.Name
			return nil
		}
	}

	return ErrUserDoesNotHaveRole
}

func (u *User) HasPermission(permission string) bool {
	for _, role := range u.roles {
		if role.Id == u.activeRole {
			for _, perm := range role.Permissions {
				if perm == permission {
					return true
				}
			}
		}
	}

	return false
}

func (u *User) HashedPassword() string {
	return u.hashedPassword
}

func (u *User) SetHashedPassword(hashedPassword string) {
	u.hashedPassword = hashedPassword
}
