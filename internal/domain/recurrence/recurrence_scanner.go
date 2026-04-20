package recurrence

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

func (r *RecurrenceRule) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("Failed to scan RecurrenceRule: unsupported type: %T", value)
	}

	if len(bytes) == 0 {
		return nil
	}

	return json.Unmarshal(bytes, r)
}

func (r RecurrenceRule) Value() (driver.Value, error) {
	if r.RecurrenceType == "" &&
		len(r.RecurrenceModifiers) == 0 &&
		len(r.Days) == 0 {
		return nil, nil
	}

	return json.Marshal(r)
}
