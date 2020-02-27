package nseconfig

import (
	"bytes"
	"fmt"
	"net"
	"strings"
)

type config interface {
	validate() error
}

type InvalidConfigErrors []error

func (v InvalidConfigErrors) Error() string {
	b := bytes.NewBufferString("validation failed with errors: \n")
	for _, err := range v {
		fmt.Fprintf(b, "\t%s\n", err)
	}
	return b.String()
}

type Config struct {
	Endpoints []*Endpoint
}

func (c Config) validate() error {
	if len(c.Endpoints) == 0 {
		return fmt.Errorf("no endpoints provided")
	}

	var errs InvalidConfigErrors

	for _, endp := range c.Endpoints {
		if err := endp.validate(); err != nil {
			if verr, ok := err.(InvalidConfigErrors); ok {
				errs = append(errs, verr...)
			} else {
				errs = append(errs, err)
			}
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

type Endpoint struct {
	CNNS CNNS

	VL3 VL3
}

func (c Endpoint) validate() error {
	var errs InvalidConfigErrors
	if len(errs) > 0 {
		return errs
	}
	for _, err := range []error{
		c.CNNS.validate(),
		c.VL3.validate(),
	} {
		if err != nil {
			if verr, ok := err.(InvalidConfigErrors); ok {
				errs = append(errs, verr...)
			} else {
				errs = append(errs, err)
			}
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

type CNNS struct {
	Name        string
	Address     string
	AccessToken string
}

func (c CNNS) validate() error {
	var errs InvalidConfigErrors
	if empty(c.Address) {
		errs = append(errs, fmt.Errorf("cnns addreses is not set"))
	}
	if empty(c.Name) {
		errs = append(errs, fmt.Errorf("cnns name is not set"))
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

type VL3 struct {
	IMAP        IMAP
	Ifname      string
	NameServers []string
}

func (v VL3) validate() error {
	var errs InvalidConfigErrors

	if _, _, err := net.ParseCIDR(v.IMAP.PrefixPool); err != nil {
		errs = append(errs, fmt.Errorf("prefix pool is not a valid subnet: %s", err))
	}
	for i, r := range v.IMAP.Routes {
		if _, _, err := net.ParseCIDR(r); err != nil {
			errs = append(errs, fmt.Errorf("route nr %d with value %s is not a valid subnet: %s", i, r, err))
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

type IMAP struct {
	PrefixPool string
	Routes     []string
}

type decoder interface {
	Decode(v interface{}) error
}

type DecoderFn func(v interface{}) error

func (d DecoderFn) Decode(v interface{}) error { return d(v) }

func NewConfig(decoder decoder, cfg config) error {
	if err := decoder.Decode(cfg); err != nil {
		return err
	}

	return cfg.validate()
}

func empty(s string) bool {
	return len(strings.Trim(s, " ")) == 0
}
