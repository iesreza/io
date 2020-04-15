package jwt

import (
	"encoding/json"
	"github.com/gbrlsnchs/jwt"
	"github.com/lithammer/shortuuid"
	"time"
)


type Config struct{
	Secret     string        `yaml:"secret"`
	Issuer     string        `yaml:"issuer"`
	Audience   []string      `yaml:"audience"`
	Age        time.Duration `yaml:"age"`
	Subject    string        `yaml:"subject"`
}

type Payload struct {
	jwt.Payload
	Data     map[string]interface{}  `json:"data,omitempty"`
	Empty    bool                    `json:"empty,omitempty"`
}

var Hash  *jwt.HMACSHA
var config Config
func Register(c string) {
    json.Unmarshal([]byte(c),&config)
	Hash = jwt.NewHS256([]byte(config.Secret))
}

func Generate(data map[string]interface{},extend... time.Duration) (string,error) {
	now := time.Now()
	pl := Payload{
		Payload: jwt.Payload{
			Issuer:         config.Issuer,
			Subject:        config.Subject,
			Audience:       jwt.Audience(config.Audience),
			ExpirationTime: jwt.NumericDate(now.Add(config.Age)),
			NotBefore:      jwt.NumericDate(now),
			IssuedAt:       jwt.NumericDate(now),
			JWTID:          shortuuid.New(),
		},
		Empty:false,
		Data: data,
	}
	if len(extend) > 0{
		pl.Payload.ExpirationTime = jwt.NumericDate(now.Add(extend[0]))
	} else if d,exist := pl.Get("_extend_duration"); exist{
		duration := d.(time.Duration)
		pl.Payload.ExpirationTime = jwt.NumericDate(now.Add(duration))
	}
	token, err := jwt.Sign(pl, Hash)
	if err != nil {
		return "",err
	}
	return string(token),nil
}

func Verify(token string) (Payload,error) {

	var pl Payload
	_, err := jwt.Verify([]byte(token), Hash, &pl)
	if err != nil {
		return pl,err
	}
	return pl,nil
}

func (p *Payload) Set(key string,value interface{}) {
	p.Data[key] = value
	p.Empty = false
}

func (p *Payload) Get(key string) (interface{},bool) {
	if val, ok := p.Data[key]; ok {
		return val,true
	}
	return nil,false
}

func (p *Payload) Remove(key string){
	delete(p.Data,key)
}

func (p *Payload) ExtendPeriod(d time.Duration){
	p.Set("_extend_duration",d)
}