package model

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type TimeZone struct {
	*time.Location
}

// Value implements the driver.Valuer interface.
func (tz TimeZone) Value() (driver.Value, error) {
	return tz.Location.String(), nil
}

func (tz TimeZone) Equal(other TimeZone) bool {
	return tz.String() == other.String()
}

// Scan implements the sql.Scanner interface.
// A 16-byte slice will be handled by UnmarshalBinary, while
// a longer byte slice or a string will be handled by UnmarshalText.
func (tz *TimeZone) Scan(src interface{}) error {
	switch src := src.(type) {
	case string:
		// uu, err := FromString(src)
		loc, err := time.LoadLocation(src)
		if err != nil {
			return fmt.Errorf("failed to convert TimeZone: %w", err)
		}
		tz.Location = loc
		return nil
	}
	return fmt.Errorf("TimeZone: cannot convert %T to TimeZone", src)
}
