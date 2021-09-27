// Copyright 2021 Ke Fan <litesky@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package options

// Validate checks Options and return a slice of found errs.
func (o *Options) Validate() []error {
	var errs []error

	errs = append(errs, o.K8sOption.Validate()...)
	errs = append(errs, o.DockerOption.Validate()...)
	errs = append(errs, o.Log.Validate()...)

	return errs
}
