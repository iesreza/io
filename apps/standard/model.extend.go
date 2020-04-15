package __PACKAGE__

import (
	"encoding/json"
	"gopkg.in/dealancer/validate.v2"
)

func (m __MODEL__)All() []__MODEL__  {
	var res []__MODEL__
	db.Find(&res)
	return res
}

func (m __MODEL__)FindAll(query... interface{}) []__MODEL__  {
	var res []__MODEL__
	if len(query) > 1{
		db.Where(query[0],query[1:]).Find(&res)
	}else {
		db.Find(&res)
	}
	return res
}

func (m __MODEL__)FindFirst(query... interface{}) bool  {
	var res __MODEL__
	if len(query) > 1 {
		return db.Where(query[0],query[1:]).First(&res).RecordNotFound()
	}else{
		return db.First(&res).RecordNotFound()
	}
}

func (m *__MODEL__)FindLast(query... interface{}) bool  {
	var res __MODEL__
	if len(query) > 1 {
		return  db.Where(query[0],query[1:]).Last(&res).RecordNotFound()
	}else{
		return  db.Last(&res).RecordNotFound()
	}
}

func (m *__MODEL__)Save() (*__MODEL__,error){
	var err error
	if m.ID > 0 {
		err = db.Save(m).Error
	}else{
		err = db.Create(m).Error
	}
	return m,err
}

func (m *__MODEL__)Delete(){
	if m.ID > 0 {
		db.Delete(m)
	}
}


func (m *__MODEL__)Validate() error{
	return validate.Validate(m)
}

func (m __MODEL__)Json() string {
	b,_ := json.Marshal(m)
	return string(b)
}

func (m *__MODEL__)FromJson(v string) error {
	return json.Unmarshal([]byte(v),m)
}

