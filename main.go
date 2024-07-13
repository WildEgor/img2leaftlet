// Package main ...
package main

import (
	"github.com/WildEgor/img2leaftlet/cmd"
	"github.com/WildEgor/img2leaftlet/internal/logger"
)

func main() {
	logger.Init()
	cmd.Execute()
}
