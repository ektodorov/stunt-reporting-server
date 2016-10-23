package objects

import (

)

type Report struct {
	Id int
	ClientId string
	Time int
	Sequence int
	Message string
	FilePath string
}