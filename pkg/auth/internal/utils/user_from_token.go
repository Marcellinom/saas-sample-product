package utils

import (
	"reflect"

	"github.com/lestrrat-go/jwx/jwt"
	"its.ac.id/base-go/pkg/auth/contracts"
)

func UserFromToken(t jwt.Token) *contracts.User {
	u := contracts.NewUser(t.Subject())
	rolesInterface, exist := t.Get("roles")
	if !exist {
		return u
	}

	s := reflect.ValueOf(rolesInterface)

	for i := 0; i < s.Len(); i++ {
		roleMap := s.Index(i).Elem()
		name := roleMap.MapIndex(reflect.ValueOf("name")).Interface().(string)
		permissionsInterface := roleMap.MapIndex(reflect.ValueOf("permissions")).Interface().([]interface{})
		isDefault := roleMap.MapIndex(reflect.ValueOf("is_default")).Interface().(bool)

		permissionsReflect := reflect.ValueOf(permissionsInterface)
		permissions := make([]string, permissionsReflect.Len())
		for j := 0; j < permissionsReflect.Len(); j++ {
			permissions[j] = permissionsReflect.Index(j).Interface().(string)
		}

		u.AddRole(name, permissions, isDefault)
	}
	activeRoleInterface, exist := t.Get("active_role")
	if !exist {
		return u
	}
	activeRole, ok := activeRoleInterface.(string)
	if !ok {
		return u
	}
	u.SetActiveRole(activeRole)

	return u
}
