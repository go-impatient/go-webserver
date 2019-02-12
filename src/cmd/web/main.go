package main

import (
	"fmt"
	"os"

	"github.com/moocss/go-webserver/src/config"
	"github.com/moocss/go-webserver/src/log"
	"github.com/moocss/go-webserver/src/pkg/version"
	"github.com/moocss/go-webserver/src/server"
	"github.com/moocss/go-webserver/src/storer"
	"github.com/urfave/cli"
	"golang.org/x/sync/errgroup"
)

var usageStr = `
              ___.                                           
__  _  __ ____\_ |__   ______ ______________  __ ___________ 
\ \/ \/ // __ \| __ \ /  ___// __ \_  __ \  \/ // __ \_  __ \
 \     /\  ___/| \_\ \\___ \\  ___/|  | \/\   /\  ___/|  | \/
  \/\_/  \___  >___  /____  >\___  >__|    \_/  \___  >__|   
             \/    \/     \/     \/                 \/      
Usage: webserver [options]
Server Options:
	-c, --config <file>              Configuration file path
	-a, --address <address>          Address to bind (default: any)
	-p, --port <port>                Use port for clients (default: 9090)
Common Options:
	-h, --help                       Show this message
	-v, --version                    Show version
`

var flags = []cli.Flag{
	cli.BoolFlag{
		EnvVar: "WEBSERVER_DEBUG",
		Name:   "debug",
		Usage:  "enable server debug mode",
	},
	cli.StringFlag{
		EnvVar: "WEBSERVER_CONFING",
		Name:   "config, c",
		Usage:  "set config file",
	},
}

func start(c *cli.Context) error {
	var (
		g errgroup.Group
	)

	// 初始化数据
	storer.DB.Init()

	// 设置默认配置
	if err := config.Init(c.String("c")); err != nil {
		log.Infof("Load yaml config file error: '%v'", err)
	}

	// overwrite server port and address
	if c.String("port") != "" {
		config.Bear.C.Core.Port = c.String("port")
	}
	if c.String("host") != "" {
		config.Bear.C.Core.Host = c.String("host")
	}

	g.Go(func() error {
		// 启动服务
		return server.RunHTTPServer()
	})
	g.Go(func() error {
		// 健康检查
		return server.PingServer()
	})

	if err := g.Wait(); err != nil {
		log.Error("接口服务出错了：", err)
	}

	return g.Wait()
}

func run() {
	app := cli.NewApp()
	app.Name = "webserver"
	app.Version = version.Info.String() // version.Version.String()
	app.Usage = "go web server"
	app.UsageText = usageStr
	app.Action = start
	app.Flags = flags
	app.Before = func(c *cli.Context) error {
		fmt.Fprintf(c.App.Writer, "brace for impact\n")
		return nil
	}
	app.After = func(c *cli.Context) error {
		fmt.Fprintf(c.App.Writer, "did we lose anyone?\n")
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}

func main() {
	run()
}