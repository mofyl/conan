package log

import (
	"fmt"
	"go.uber.org/zap"
	"testing"
)

func TestGetConfig(t *testing.T) {
	cfg := getDefaultLog()

	fmt.Println(cfg.OutFile)
}

func TestLog(t *testing.T) {
	Debug("qqq", zap.String("asda", "asda"))
}
