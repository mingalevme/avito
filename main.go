package main

import (
	"github.com/alecthomas/kong"
	"github.com/mingalevme/avito/cmd"
	"github.com/mingalevme/avito/internal/env"
)

var CLI struct{
	Check cmd.CheckCmd `cmd help:"Check"`
}

func main() {
	e := env.New("APP_", env.GetOSEnvMap())
	ctx := kong.Parse(&CLI)
	_ = ctx.Run(e)
}
