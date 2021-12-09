// Copyright 2021 Ke Fan <litesky@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package puller

import (
	"github.com/tianrandailove/peitho/internal/puller/config"
	"github.com/tianrandailove/peitho/internal/puller/options"
	"github.com/tianrandailove/peitho/pkg/app"
	"github.com/tianrandailove/peitho/pkg/log"
)

const commandDesc = "puller"

func NewApp(basename string) *app.App {
	opts := options.NewOptions()
	app := app.NewApp("puller",
		basename,
		app.WithOptions(opts),
		app.WithDescription(commandDesc),
		app.WithDefaultValidArgs(),
		app.WithRunFunc(run(opts)),
		app.WithNoConfig(),
	)

	return app
}

func run(opts *options.Options) app.RunFunc {
	return func(basename string) error {
		log.Init(opts.Log)
		defer log.Flush()

		cfg, err := config.CreateConfigFromOptions(opts)
		if err != nil {
			return err
		}

		return Run(cfg)
	}
}
