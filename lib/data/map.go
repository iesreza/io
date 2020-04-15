package data

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)




type Map map[string]interface{}

func (p Map) Value() (driver.Value, error) {
	j, err := json.Marshal(p)
	return j, err
}

func (p *Map) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("Type assertion .([]byte) failed.")
	}

	var i interface{}
	err := json.Unmarshal(source, &i)
	if err != nil {
		return err
	}

	*p, ok = i.(map[string]interface{})
	if !ok {
		return fmt.Errorf("Type assertion .(map[string]interface{}) failed.")
	}

	return nil
}