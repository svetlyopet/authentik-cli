package logger

import (
	"fmt"
)

func LogObjectChange(targetObjectType, action, targetObjectName string) {
	fmt.Printf("%s/%s %s\n", targetObjectType, targetObjectName, action)
}

func LogObjectDetails(resourceData []byte) {
	fmt.Printf("%s\n", string(resourceData))
}
