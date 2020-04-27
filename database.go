package io

import (
	"fmt"
	"github.com/iesreza/foundation/lib/log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"strings"
	"time"
)

var Database *gorm.DB

func SetupDatabase() {
	Events.Go("database.starts")
	config := config.Database
	var err error
	if config.Enabled == false {
		return
	}
	switch strings.ToLower(config.Type) {
	case "mysql":
		connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?"+config.Params, config.Username, config.Password, config.Server, config.Database)
		Database, err = gorm.Open("mysql", connectionString)
	case "postgres":
		connectionString := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=%s "+config.Params, config.Username, config.Password, config.Server, config.Database, config.SSLMode)
		Database, err = gorm.Open("postgres", connectionString)
	case "mssql":
		connectionString := fmt.Sprintf("user id=%s;password=%s;server=%s;database:%s;"+config.Params, config.Username, config.Password, config.Server, config.Database)
		Database, err = gorm.Open("mssql", connectionString)
	default:
		Database, err = gorm.Open("sqlite3", config.Database+config.Params)
	}
	Database.LogMode(config.Debug == "true")
	if err != nil {
		log.Critical(err)
	}
	Events.Go("database.started")

}

func GetDBO() *gorm.DB {
	if Database == nil {
		SetupDatabase()
	}
	return Database
}

type Model struct {
	ID        uint       `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" sql:"index"`
}

/*type Map map[string]interface{}

func (j Map) Value() (driver.Value, error) {
	if j.IsNull() {
		return nil, nil
	}
	return text.ToJSON(j), nil
}

func (j *Map) Scan(value interface{}) error {
	if lib.IsNil(value) {
		*j = nil
		return nil
	}
	s, ok := value.([]byte)
	if !ok {
		fmt.Errorf("Invalid Scan Source")
	}
	err := json.Unmarshal(s,j)
	if err != nil{
		return err
	}
	return nil
}

func (m Map) MarshalJSON() ([]byte, error) {
	fmt.Println(m)
	return json.Marshal(m)
}

func (m *Map) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data,m)
}

func (j Map) IsNull() bool {
	return lib.IsNil(j)
}
*/
