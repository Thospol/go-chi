package json

import (
	"saaa-api/internal/core/utils"
)

// Init init read json file
func Init(path string, v interface{}) error {
	err := utils.ReadJSONFile(path, v)
	if err != nil {
		return err
	}

	return nil
}
