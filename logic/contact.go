package logic

import "encoding/json"

type Contact struct {
	ID string `json:"contact_id"`
}

func ParseContact(data []byte) (contact Contact, err error) {
	err = json.Unmarshal(data, &contact)
	return
}
