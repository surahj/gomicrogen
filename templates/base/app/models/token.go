package models

import tokenutils "github.com/mudphilo/gwt"

type TokenData struct {
	UserID   int64           `json:"user_id"`
	UserName string          `json:"user_name"`
	Expiry   int64           `json:"expiry"`
	Role     tokenutils.Role `json:"role"`
}
