// Copyright 2021 Ke Fan <litesky@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package options

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/pflag"
)

// PullerOption defines options for docker cluster.
type PullerOption struct {
	DockerEndpoint string `json:"dockerEndpoint" mapstructure:"dockerEndpoint"`
	Image          string `json:"image"          mapstructure:"image"`
	PullAddress    string `json:"pullAddress"    mapstructure:"pullAddress"`
}

// NewPullerOption create a `zero` value instance.
func NewPullerOption() *PullerOption {
	return &PullerOption{
		DockerEndpoint: "",
		Image:          "",
		PullAddress:    "",
	}
}

// Validate validate option value.
func (o *PullerOption) Validate() []error {
	errs := []error{}

	if o.DockerEndpoint == "" {
		errs = append(errs, fmt.Errorf("docker endpoint can not be empty"))
	}

	if o.Image == "" {
		errs = append(errs, fmt.Errorf("image cannot be empty"))
	}

	if o.PullAddress == "" {
		errs = append(errs, fmt.Errorf("pullAddress cannot be empty"))
	}

	return errs
}

// AddFlags bind command flag.
func (o *PullerOption) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&(o.DockerEndpoint), "docker.endpoint", o.DockerEndpoint, "The endpoint for accessing docker server.")
	fs.StringVar(&(o.Image), "image", o.Image, "the image name")
	fs.StringVar(&(o.PullAddress), "pullAddress", o.PullAddress, "the url to download the image")
}

// String to json string.
func (o *PullerOption) String() string {
	data, _ := json.Marshal(o)

	return string(data)
}
