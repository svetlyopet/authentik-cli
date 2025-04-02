package logger

import (
	"fmt"
)

func WriteStdio(targetObjectType, action, targetObjectName string) {
	fmt.Printf("%s/%s %s\n", targetObjectType, targetObjectName, action)
}
