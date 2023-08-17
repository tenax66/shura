package main

import (
	"fmt"

	"github.com/tenax66/shura"
	"go.uber.org/zap"
)

func main() {
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
