// Copyright 2021 Ke Fan <litesky@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package config

import "github.com/tianrandailove/peitho/internal/puller/options"

// Config is the running configuration structure of the IAM pump service.
type Config struct {
	*options.Options
}

// CreateConfigFromOptions creates a running configuration instance based.
func CreateConfigFromOptions(opts *options.Options) (*Config, error) {
	return &Config{opts}, nil
}
