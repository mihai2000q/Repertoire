package internal

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type Date time.Time

// Scan implements the sql.Scanner interface for database deserialization
func (d *Date) Scan(value interface{}) error {
	var t time.Time
	switch v := value.(type) {
	case time.Time:
		t = v
	case []byte:
		var err error
		t, err = time.Parse("2006-01-02", string(v))
		if err != nil {
			return err
		}
	case string:
		var err error
		t, err = time.Parse("2006-01-02", v)
		if err != nil {
			return err
		}
	case nil:
		return nil
	default:
		return fmt.Errorf("unsupported type for Date: %T", value)
	}

	*d = Date(t.UTC().Truncate(24 * time.Hour))
	return nil
}

// Value implements the driver.Valuer interface for database serialization
func (d Date) Value() (driver.Value, error) {
	if time.Time(d).IsZero() {
		return nil, nil
	}
	y, m, day := time.Time(d).Date()
	return time.Date(y, m, day, 0, 0, 0, 0, time.UTC), nil
}

// GormDataType specifies the database type for GORM
func (Date) GormDataType() string {
	return "date"
}

// MarshalJSON implements custom JSON marshaling
func (d Date) MarshalJSON() ([]byte, error) {
	if time.Time(d).IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + time.Time(d).UTC().Format("2006-01-02") + `"`), nil
}

// UnmarshalJSON implements custom JSON unmarshaling
func (d *Date) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return fmt.Errorf("invalid date format: %v", err)
	}

	*d = Date(t.UTC().Truncate(24 * time.Hour))
	return nil
}

// String returns the date in ISO8601 format
func (d Date) String() string {
	return time.Time(d).Format("2006-01-02")
}

// GobEncode implements the gob.GobEncoder interface
func (d Date) GobEncode() ([]byte, error) {
	return time.Time(d).GobEncode()
}

// GobDecode implements the gob.GobDecoder interface
func (d *Date) GobDecode(data []byte) error {
	return (*time.Time)(d).GobDecode(data)
}
