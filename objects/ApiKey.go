package objects

import (
	"strconv"
)

type ApiKey struct {
	UserId int
	ApiKey string
	AppName string
}

func (receiver *ApiKey) String() string {
	var retVal string
	retVal = "{UserId=" + strconv.Itoa(receiver.UserId) + ", ApiKey=" + receiver.ApiKey + ", AppName=" + receiver.AppName + "}"
	return retVal
} 