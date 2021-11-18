package options

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/pflag"
)

// SweeperOption define opotion for chaincode sweeper
type SweeperOption struct {
	Enable   bool `json:"enable" mapstructure:"enable"`
	Interval int  `json:"interval" mapstructure:"interval"`
}

// NewSweeperOption create a zero value instance
func NewSweeperOption() *SweeperOption {
	return &SweeperOption{
		Enable:   true,
		Interval: 60,
	}
}

// Validate validate option value.
func (o *SweeperOption) Validate() []error {
	errs := []error{}

	if o.Interval <= 0 {
		errs = append(errs, fmt.Errorf("Interval cannot be zero"))
	}

	return errs
}

// AddFlags bind command flag.
func (o *SweeperOption) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVar(
		&(o.Enable),
		"sweeper.enable",
		o.Enable,
		"enable to sweep that restarting chaincode deployment",
	)
	fs.IntVar(&(o.Interval), "sweeper.interval", o.Interval, "interval for sweeping")
}

// String to json string.
func (o *SweeperOption) String() string {
	data, _ := json.Marshal(o)

	return string(data)
}
