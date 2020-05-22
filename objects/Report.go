package objects

import (
	
)

type Report struct {
	ApiKey string	`json:"apikey"`
	Id int			`json:"id"`
	ClientId string	`json:"clientid"`
	Time int64		`json:"time"`
	Sequence int	`json:"sequence"`
	Message string	`json:"message"`
	FilePath string	`json:"filepath"`
	TimeString string
}