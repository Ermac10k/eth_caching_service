package logger

import (
	"fmt"
	"time"
)

// Logger is custom logger
type Logger struct {}

func (writer Logger) Write(bytes []byte) (int, error) {
    return fmt.Print(time.Now().Format("2021-12-01 15:04:05.999 : ") + string(bytes))
}