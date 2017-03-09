package couchbase

import (
	"bytes"
	"fmt"
)

type User struct {
	Name  string
	Id    string
	Type  string
	Roles []Role
}

type Role struct {
	Role       string
	BucketName string `json:"bucket_name"`
}

// Return user-role data, as parsed JSON.
// Sample:
//   [{"id":"ivanivanov","name":"Ivan Ivanov","roles":[{"role":"cluster_admin"},{"bucket_name":"default","role":"bucket_admin"}]},
//    {"id":"petrpetrov","name":"Petr Petrov","roles":[{"role":"replication_admin"}]}]
func (c *Client) GetUserRoles() ([]interface{}, error) {
	ret := make([]interface{}, 0, 1)
	err := c.parseURLResponse("/settings/rbac/users", &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (c *Client) GetUserInfoAll() ([]User, error) {
	ret := make([]User, 0, 16)
	err := c.parseURLResponse("/settings/rbac/users", &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func rolesToParamFormat(roles []Role) string {
	var buffer bytes.Buffer
	for i, role := range roles {
		if i > 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString(role.Role)
		if role.BucketName != "" {
			buffer.WriteString("[")
			buffer.WriteString(role.BucketName)
			buffer.WriteString("]")
		}
	}
	return buffer.String()
}

func (c *Client) PutUserInfo(u *User) error {
	params := map[string]interface{}{
		"name":  u.Name,
		"roles": rolesToParamFormat(u.Roles),
	}
	var target string
	switch u.Type {
	case "saslauthd":
		target = "/settings/rbac/users/" + u.Id
	case "builtin":
		target = "/settings/rbac/users/builtin/" + u.Id
	default:
		return fmt.Errorf("Unknown user type: %s", u.Type)
	}
	var ret string // PUT returns an empty string. We ignore it.
	err := c.parsePutURLResponse(target, params, &ret)
	return err
}
