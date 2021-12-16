// Copyright 2021 Ke Fan <litesky@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package options

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/pflag"
)

const (
	IMAGE_MODE_REGISTRY = "registry"
	IMAGE_MODE_DELIVERY = "delivery"
)

// PeithoOption defines options for peitho.
type PeithoOption struct {
	ImageMode           string `json:"imageMode"           mapstructure:"imageMode"`
	PullerAccessAddress string `json:"pullerAccessAddress" mapstructure:"pullerAccessAddress"`
	PullerImage         string `json:"pullerImage"         mapstructure:"pullerImage"`
}

// NewPeithoOption create a `zero` value instance.
func NewPeithoOption() *PeithoOption {
	return &PeithoOption{
		ImageMode:           IMAGE_MODE_REGISTRY,
		PullerAccessAddress: "",
	}
}

// Validate validate option value.
func (o *PeithoOption) Validate() []error {
	errs := []error{}

	if IMAGE_MODE_DELIVERY != o.ImageMode && IMAGE_MODE_REGISTRY != o.ImageMode {
		errs = append(errs, fmt.Errorf("imageMode must be registry or deivery"))
	}

	if o.ImageMode == IMAGE_MODE_DELIVERY && o.PullerAccessAddress == "" {
		errs = append(errs, fmt.Errorf("pullerAccessAddress must not be empty"))
	}

	if o.ImageMode == IMAGE_MODE_DELIVERY && o.PullerImage == "" {
		errs = append(errs, fmt.Errorf("pullerImage must not be empty"))
	}

	return errs
}

// AddFlags bind command flag.
func (o *PeithoOption) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(
		&(o.ImageMode),
		"imageMode",
		o.ImageMode,
		"how to delivery chaincode image, if choose registry mode, then pull image from registry, else from peitho",
	)
	fs.StringVar(
		&(o.PullerAccessAddress),
		"pullerAccessAddress",
		o.PullerAccessAddress,
		"puller access the url for pulling image",
	)
	fs.StringVar(&(o.PullerImage), "pullerImage", o.PullerImage, "the pullerImage for chancode initcontainer")
}

// String to json string.
func (o *PeithoOption) String() string {
	data, _ := json.Marshal(o)

	return string(data)
}
