package types

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"strings"
)

// UUIDArray adalah custom type untuk field Postgres uuid[].
type UUIDArray []uuid.UUID

// Scan membaca data dari database dan mengubahnya menjadi UUIDArray.
func (a *UUIDArray) Scan(value interface{}) error {
	var str string

	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	default:
		return errors.New("failed to parse UUIDArray: unsupported data type")
	}

	str = strings.TrimPrefix(str, "{")
	str = strings.TrimSuffix(str, "}")
	parts := strings.Split(str, ",")

	*a = make(UUIDArray, 0, len(parts))
	for _, s := range parts {
		s = strings.TrimSpace(strings.Trim(s, `"`)) // hapus quotes jika ada
		if s == "" {
			continue
		}
		u, err := uuid.Parse(s)
		if err != nil {
			return fmt.Errorf("invalid UUID in array: %v", err)
		}
		*a = append(*a, u)
	}

	return nil
}

// Value mengubah UUIDArray menjadi format Postgres {uuid1,uuid2,...}.
func (a UUIDArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "{}", nil
	}

	strs := make([]string, 0, len(a))
	for _, u := range a {
		strs = append(strs, fmt.Sprintf(`"%s"`, u.String()))
	}
	return "{" + strings.Join(strs, ",") + "}", nil
}

// GormDataType memberi tahu GORM untuk mapping ke uuid[] di Postgres.
func (UUIDArray) GormDataType() string {
	return "uuid[]"
}
