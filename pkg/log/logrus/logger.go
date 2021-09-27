// Copyright 2021 Ke Fan <litesky@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package logrus adds a hook to the logrus logger hooks.
package logrus

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

// NewLogger create a logrus logger, add hook to it and return it.
func NewLogger(zapLogger *zap.Logger) *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(ioutil.Discard)
	logger.AddHook(newHook(zapLogger))

	return logger
}
