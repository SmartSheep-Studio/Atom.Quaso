package hyperutils

import (
	"fmt"
	"strings"
)

func GenScope(scopes ...string) []string {
	var data []string
	for _, scope := range scopes {
		operator := strings.SplitN(scope, ":", 2)[0]
		object := strings.SplitN(scope, ":", 2)[1]
		data = append(data, fmt.Sprintf("%s:quaso.%s", operator, object))
	}
	return data
}

func GenPerms(perms ...string) []string {
	var data []string
	for _, perm := range perms {
		data = append(data, fmt.Sprintf("quaso.%s", perm))
	}
	return data
}
