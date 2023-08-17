package main

import (
	"fmt"
	"os"

	"github.com/tenax66/shura"
	"go.uber.org/zap"
)

func main() {
	os.Setenv("HTTP_PROXY", "socks5://localhost:9050")

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	result, err := shura.Run()

	if err != nil {
		sugar.Warn(err)
		return
	}

	fmt.Print(result)
}
