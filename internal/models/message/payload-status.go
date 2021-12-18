package message

import (
	"encoding/json"
	"strings"
	"time"

	l "github.com/redhatinsights/payload-tracker-go/internal/logging"
)

var (
	dateFormat = time.RFC3339
)

// PayloadStatusMessage it the definition of the Payload Message Status kafka message
type PayloadStatusMessage struct {
	Service     string       `json:"service"`
	Source      string       `json:"source,omitempty"`
	Account     string       `json:"account,omitempty"`
	RequestID   string       `json:"request_id"`
	InventoryID string       `json:"inventory_id,omitempty"`
	SystemID    string       `json:"system_id,omitempty"`
	Status      string       `json:"status"`
	StatusMSG   string       `json:"status_msg,omitempty"`
	PayloadID   uint         `json:"payload_id,omitempty"`
	Date        FormatedTime `json:"date"`
}

type FormatedTime struct {
	time.Time
}

func (t *FormatedTime) UnmarshalJSON(b []byte) error {
	var date string
	err := json.Unmarshal(b, &date)
	if err != nil {
		l.Log.Error("ERROR: Unmarshaling time: ", err)
		return err
	}

	date = strings.Join(strings.Fields(date), "T")
	
	// Add a Z to the end of the timestamp if it doesn't exist
	if !strings.HasSuffix(date, "Z") {
		date = date + "Z"
	}

	t.Time, err = time.Parse(dateFormat, date)
	if err != nil {
		l.Log.Error("ERROR: Parsing date into new format: ", err)
		return err
	}

	return nil
}
