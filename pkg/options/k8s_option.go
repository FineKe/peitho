// Copyright 2021 Ke Fan <litesky@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package options

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

// K8sOption defines options for kubernates cluster.
type K8sOption struct {
	Kubeconfig string   `json:"kubeconfig" mapstructure:"kubeconfig"`
	Namespace  string   `json:"namespace"  mapstructure:"namespace"`
	DNS        []string `json:"dns"        mapstructure:"dns"`
}

// NewK8sOption create a `zero` value instance.
func NewK8sOption() *K8sOption {
	return &K8sOption{
		Kubeconfig: "",
		Namespace:  "",
		DNS:        []string{},
	}
}

// Validate validate option value.
func (o *K8sOption) Validate() []error {
	errs := []error{}

	if o.Kubeconfig == "" {
		errs = append(errs, fmt.Errorf("kubeconfig path can not be empty"))
	} else if _, err := os.Stat(o.Kubeconfig); err != nil {
		errs = append(errs, fmt.Errorf("kubeconfig file not exists"))
	}

	if o.Namespace == "" {
		errs = append(errs, fmt.Errorf("namespace cannot be empty"))
	}

	return errs
}

// AddFlags bind command flag.
func (o *K8sOption) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(
		&(o.Kubeconfig),
		"k8s.kubeconfig",
		o.Kubeconfig,
		"Path to kubeconfig file for accessing kubernates cluster.",
	)
	fs.StringVar(&(o.Namespace), "k8s.namespace", o.Namespace, "Namespace for accessing kubernates cluster.")
	fs.StringSliceVar(&o.DNS, "k8s.dns", o.DNS, "DNS for chaincode to resolve peer node")
}

// String to json string.
func (o *K8sOption) String() string {
	data, _ := json.Marshal(o)

	return string(data)
}
