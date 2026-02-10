package main

import (
	"github.com/alecthomas/kong"
	"github.com/bakedSpaceTime/binip/libip"
	"github.com/bakedSpaceTime/binip/libip/config"
)

type AppCmd struct {
}

func (a *AppCmd) Run(c *config.Config) error {
	return libip.App(c)
}

type Info struct {
}

func (i *Info) Run(c *config.Config) error {
	return libip.Info(c)
}

type Test struct {
}

func (t *Test) Run(c *config.Config) error {
	return libip.Test(c)
}

type Reset struct {
}

func (r *Reset) Run(c *config.Config) error {
	return libip.Reset(c)
}

var cli struct {
	App   AppCmd `cmd:"" default:"withargs" help:"Main App."`
	Info  Info   `cmd:"" help:"Show system info"`
	Test  Test   `cmd:"" help:"Run in test mode"`
	Reset Reset  `cmd:"" help:"Reset app db"`
	Debug bool   `help:"Enable debug mode."`
}

func main() {
	ctx := kong.Parse(&cli)

	c := config.NewConfig()
	c.Debug = cli.Debug

	err := ctx.Run(c)
	ctx.FatalIfErrorf(err)
}
