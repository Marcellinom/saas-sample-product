package config

import (
	"github.com/dptsi/its-go/sso"
)

func ssoConfig() sso.Config {
	return sso.Config{
		IsRoleFromSso:               true,
		IsGroupMappedDirectlyToRole: true,
		Roles:                       map[string]sso.Role{},
		RolePermissions:             map[string][]string{},
		GroupRoleMapping:            map[string]string{},
	}
}
