package objects

import (

)

type ClientInfo struct {
	ApiKey string		`json:"apikey"`
	ClientId string		`json:"clientid"`
	Name string			`json:"name"`
	Manufacturer string	`json:"manufacturer"`
	Model string		`json:"model"`
	DeviceId string		`json:"deviceid"`
}