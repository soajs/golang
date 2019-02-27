package soajsgo

import (
	"errors"
	"fmt"
	"regexp"
)

type (
	// Config represents service configuration from json file.
	// see: https://soajsorg.atlassian.net/wiki/spaces/SOAJ/pages/61347270/Service
	Config struct {
		Type          string `json:"type"`
		Prerequisites struct {
			CPU    string `json:"cpu"`
			Memory string `json:"memory"`
		} `json:"prerequisites"`
		ServiceName           string      `json:"serviceName"`
		ServiceGroup          string      `json:"serviceGroup"`
		ServiceVersion        string      `json:"serviceVersion"`
		ServiceIP             string      `json:"serviceIP"`
		ServicePort           int         `json:"servicePort"`
		RequestTimeout        int         `json:"requestTimeout"`
		RequestTimeoutRenewal int         `json:"requestTimeoutRenewal"`
		Swagger               bool        `json:"swagger"`
		ExtKeyRequired        bool        `json:"extKeyRequired"`
		Urac                  bool        `json:"urac"`
		UracProfile           bool        `json:"urac_Profile"`
		UracACL               bool        `json:"urac_ACL"`
		ProvisionACL          bool        `json:"provision_ACL"`
		Oauth                 bool        `json:"oauth"`
		Maintenance           maintenance `json:"maintenance"`
	}
)

var versionRegexp = regexp.MustCompile(`[0-9]+(.[0-9]+)?`)

// Validate validates soajs config.
func (c *Config) Validate() error {
	if c.Type == "" {
		return errors.New("could not find [Type] in your config, Type is <required>")
	}
	if c.ServiceName == "" {
		return errors.New("could not find [ServiceName] in your config, ServiceName is <required>")
	}
	if c.ServicePort == 0 {
		return errors.New("could not find [ServicePort] in your config, ServicePort is <required>")
	}
	if c.ServiceVersion == "" {
		return errors.New("could not find [ServiceVersion] in your config, ServiceVersion is <required>")
	}
	if !versionRegexp.MatchString(c.ServiceVersion) {
		return fmt.Errorf("error with [ServiceVersion] in your config, ServiceVersion syntax is [%s]", versionRegexp)
	}
	return nil
}
