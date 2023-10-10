package modelx

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type TimeZone struct {
	*time.Location
}

// Value implements the driver.Valuer interface.
func (tz *TimeZone) Value() (driver.Value, error) {
	return tz.Location.String(), nil
}

// Scan implements the sql.Scanner interface.
// A 16-byte slice will be handled by UnmarshalBinary, while
// a longer byte slice or a string will be handled by UnmarshalText.
func (tz *TimeZone) Scan(src interface{}) error {
	switch src := src.(type) {
	case string:
		// uu, err := FromString(src)
		loc, err := time.LoadLocation(src)
		*tz = TimeZone{loc}
		return err
	}

	return fmt.Errorf("refid: cannot convert %T to RefID", src)
}
