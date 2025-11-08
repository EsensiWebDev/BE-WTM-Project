package utils

import (
	"encoding/json"
	"github.com/jinzhu/copier"
	"gorm.io/datatypes"
)

func CopyPatch[T any](dst any, src T) error {
	return copier.CopyWithOption(dst, &src, copier.Option{IgnoreEmpty: true})
}

func CopyStrict[T any](dst any, src T) error {
	return copier.Copy(dst, &src)
}

func MapToJSON(m map[string]string) (datatypes.JSON, error) {
	b, err := json.Marshal(m)
	return datatypes.JSON(b), err
}
