package objects

import (

)

type Report struct {
	ApiKey string	`json:"apikey"`
	Id int			`json:"id"`
	ClientId string	`json:"clientid"`
	Time int		`json:"time"`
	Sequence int	`json:"sequence"`
	Message string	`json:"message"`
	FilePath string	`json:"filepath"`
}