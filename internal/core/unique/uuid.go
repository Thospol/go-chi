package unique

import (
	"strings"

	"github.com/google/uuid"
)

// UUID UUID
func UUID() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}
