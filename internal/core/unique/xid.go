package unique

import "github.com/rs/xid"

func NewXid() string {
	return xid.New().String()
}
