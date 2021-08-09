package main

import (
	"github.com/alecthomas/kong"
	"github.com/mingalevme/avito/cmd"
	"github.com/mingalevme/avito/pkg/env"
)

var CLI struct{
	Check cmd.CheckCmd `cmd help:"Check"`
}

func main() {
	namespace := env.GetEnv("MINGALEVME_AVITO_ENV_NAMESPACE", "APP_")
	e := env.New(namespace, env.GetOSEnvMap())
	ctx := kong.Parse(&CLI)
	_ = ctx.Run(e)
}
