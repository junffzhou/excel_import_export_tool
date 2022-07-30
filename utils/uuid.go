package utils

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"strings"
)


func GetUUid()string{
	return strings.Replace(fmt.Sprintf("%s", uuid.NewV1()), "-", "", -1)
}

