package id

import (
	"github.com/rs/xid"
)

// New Generates a new unique id
func New() string {
	id := xid.New()
	return id.String()
}
