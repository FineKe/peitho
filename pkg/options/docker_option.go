// Copyright 2021 Ke Fan <litesky@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package options

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/pflag"
)

// DockerOption defines options for docker cluster.
type DockerOption struct {
	Endpoint string   `json:"endpoint" mapstructure:"endpoint"`
	Registry Registry `json:"registry" mapstructure:"registry"`
}

// Registry defines options for docker registry.
type Registry struct {
	Username      string `json:"username"       mapstructure:"username"`
	Password      string `json:"password"       mapstructure:"password"`
	Email         string `json:"email"          mapstructure:"email"`
	Serveraddress string `json:"server-address" mapstructure:"server-address"`
	Project       string `json:"project"        mapstructure:"project"`
}

// NewDockerOption create a `zero` value instance.
func NewDockerOption() *DockerOption {
	return &DockerOption{
		Endpoint: "",
		Registry: Registry{},
	}
}

// Validate validate option value.
func (o *DockerOption) Validate() []error {
	errs := []error{}

	if o.Endpoint == "" {
		errs = append(errs, fmt.Errorf("docker endpoint can not be empty"))
	}

	if o.Registry.Project == "" {
		errs = append(errs, fmt.Errorf("registry project cannot be empty"))
	}

	return errs
}

// AddFlags bind command flag.
func (o *DockerOption) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&(o.Endpoint), "docker.endpoint", o.Endpoint, "The endpoint for accessing docker server.")
	fs.StringVar(
		&(o.Registry.Serveraddress),
		"docker.registry.server-address",
		o.Registry.Serveraddress,
		"docker registry server address",
	)
	fs.StringVar(&(o.Registry.Project), "docker.registry.project", o.Registry.Project, "docker registry project name")
	fs.StringVar(&(o.Registry.Email), "docker.registry.email", o.Registry.Email, "docker registry auth email")
	fs.StringVar(
		&(o.Registry.Username),
		"docker.registry.username",
		o.Registry.Username,
		"docker registry auth username",
	)
	fs.StringVar(
		&(o.Registry.Password),
		"docker.registry.password",
		o.Registry.Password,
		"docker registry auth password",
	)
}

// String to json string.
func (o *DockerOption) String() string {
	data, _ := json.Marshal(o)

	return string(data)
}
