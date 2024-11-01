package core

import (
	"encoding/json"
	"strconv"
)

type Model struct {
	CreatedAt int64 `gorm:"autoUpdateTime" json:"created_at,omitempty"`
	UpdatedAt int64 `gorm:"autoCreateTime" json:"updated_at,omitempty"`
}

type ModelId struct {
	ID uint `gorm:"primarykey"  json:"id"`
	Model
}

type ModelDelete struct {
	ModelId
	DeletedAt int64 `gorm:"index" json:"deleted_at,omitempty"`
}

type Int64 int64

func (i Int64) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.FormatInt(int64(i), 10))
}
func (i *Int64) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	parseInt, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}
	*i = Int64(parseInt)
	return nil
}
