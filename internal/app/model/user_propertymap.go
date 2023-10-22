package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"

	"golang.org/x/exp/constraints"
)

const (
	DefaultReminderThresholdHours = 24
)

func ValidateReminderThresholdHours[T constraints.Unsigned](v T) (uint8, error) {
	if v > 168 || v < 2 {
		return 0, fmt.Errorf("value outside constraints")
	}
	return uint8(v), nil
}

type UserSettings struct {
	ReminderThresholdHours uint8 `json:"reminder_threshold"`
	// weird negative name here, so zero value defaults
	// to enabling reminders
	EnableReminders bool `json:"enable_reminders"`
}

func (p UserSettings) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *UserSettings) Scan(src interface{}) error {
	var s []byte
	switch src := src.(type) {
	case UserSettings:
		*p = src
		return nil
	case *UserSettings:
		*p = *src
		return nil
	case []byte:
		s = src
	case string:
		s = []byte(src)
	default:
		return fmt.Errorf("cannot convert %T to UserPropertyMap", src)
	}

	err := json.Unmarshal(s, &p)
	if err != nil {
		return err
	}

	// set any defaults for missing values, where the zero
	// value is not what we want/need
	if p.ReminderThresholdHours == 0 {
		p.ReminderThresholdHours = DefaultReminderThresholdHours
	}

	return nil
}

func NewUserPropertyMap() *UserSettings {
	return &UserSettings{
		ReminderThresholdHours: DefaultReminderThresholdHours,
	}
}

type UserSettingsMatcher struct {
	expected UserSettings
}

func NewUserSettingsMatcher(expected UserSettings) UserSettingsMatcher {
	return UserSettingsMatcher{expected}
}

func (m UserSettingsMatcher) Match(v interface{}) bool {
	var settings UserSettings
	var err error
	switch x := v.(type) {
	case UserSettings:
		settings = x
	case *UserSettings:
		settings = *x
	case string:
		err = json.Unmarshal([]byte(x), &settings)
	case []byte:
		err = json.Unmarshal(x, &settings)
	default:
		return false
	}
	if err != nil {
		return false
	}
	return reflect.DeepEqual(m.expected, settings)
}
