package objects

import (

)

type ReportMessage struct {
	Sequence int 	`json:"sequence"`
	Time int		`json:"time"`
	Message string	`json:"message"`
	ClientId string `json:"clientid"`
}