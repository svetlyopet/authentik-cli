package logger

import (
	"fmt"
)

func WriteStdout(targetObjectType, action, targetObjectName string) {
	fmt.Printf("%s/%s %s\n", targetObjectType, targetObjectName, action)
}
