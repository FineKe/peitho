// Copyright 2021 Ke Fan <litesky@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package options

import (
	"encoding/json"

	cliflag "github.com/marmotedu/component-base/pkg/cli/flag"

	"github.com/tianrandailove/peitho/pkg/log"
	"github.com/tianrandailove/peitho/pkg/options"
)

// Options run a peitho server.
type Options struct {
	K8sOption    *options.K8sOption     `json:"k8s"    mapstructure:"k8s"`
	DockerOption *options.DockerOption  `json:"docker" mapstructure:"docker"`
	Log          *log.Options           `json:"log"    mapstructure:"log"`
	Sweeperption *options.SweeperOption `json:"sweeper" mapstructure:"sweeper"`
	PeithoOption *options.PeithoOption `json:"peitho" mapstructure:"peitho"`

}

// NewOptions creates a new Options object with default parameters.
func NewOptions() *Options {
	option := Options{
		K8sOption:    options.NewK8sOption(),
		DockerOption: options.NewDockerOption(),
		Log:          log.NewOptions(),
		Sweeperption: options.NewSweeperOption(),
		PeithoOption: options.NewPeithoOption(),
	}

	return &option
}

func (o *Options) Flags() (fss cliflag.NamedFlagSets) {
	o.K8sOption.AddFlags(fss.FlagSet("k8s"))
	o.DockerOption.AddFlags(fss.FlagSet("docker"))
	o.Log.AddFlags(fss.FlagSet("log"))
	o.Sweeperption.AddFlags(fss.FlagSet("sweeper"))
	o.PeithoOption.AddFlags(fss.FlagSet("peitho"))
  
	return fss
}

func (o *Options) String() string {
	data, _ := json.Marshal(o)

	return string(data)
}

// Complete set default Options.
func (o *Options) Complete() error {
	return nil
}
