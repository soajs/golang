package soajsgo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type (
	key int
)

const (
	// headerDataName is the SOAJS Gateway injected object attached to the header of each request
	// between the gateway and tech service.
	headerDataName = "soajsinjectobj"
	// SoajsKey use this key to init soajs data from context.
	SoajsKey = key(1)
)

// Middleware is http middleware that gets triggered per request.
func (reg *Registry) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		d, err := headerData(r)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		if d == nil {
			next.ServeHTTP(w, r)
			return
		}
		out := ContextData{
			Tenant:         d.Tenant,
			Urac:           d.Urac,
			ServicesConfig: d.Key.Config,
			Device:         d.Device,
			Geo:            d.Geo,
			Awareness:      d.Awareness,
			Reg:            reg,
		}
		out.Tenant.Key.IKey = d.Key.IKey
		out.Tenant.Key.EKey = d.Key.EKey

		out.Tenant.Application = d.Application
		out.Tenant.Application.PackageACL = d.Package.ACL
		out.Tenant.Application.PackageACLAllEnv = d.Package.ACLAllEnv

		soajs := context.WithValue(r.Context(), SoajsKey, out)
		next.ServeHTTP(w, r.WithContext(soajs))
	})
}

func headerData(r *http.Request) (*headerInfo, error) {
	headerData := r.Header.Get(headerDataName)
	if headerData == "" {
		return nil, nil
	}
	d := new(headerInfo)
	if err := json.Unmarshal([]byte(headerData), d); err != nil {
		return nil, errors.New("unable to parse SOAJS header")
	}
	return d, nil
}

// Path returns compiled service path.
func (a Host) Path(args ...string) string {
	var serviceName, version string
	switch len(args) {
	// controller
	case 1:
		serviceName = args[0]
		// controller, 1
	case 2:
		serviceName = args[0]
		version = args[1]
		// controller, 1, dash [dash is ignored]
	case 3:
		serviceName = args[0]
		version = args[1]
	}
	host := fmt.Sprintf("%s:%d/", a.Host, a.Port)
	if strings.EqualFold(serviceName, "controller") {
		host = fmt.Sprintf("%s%s/", host, serviceName)
		if _, err := strconv.Atoi(version); err == nil {
			host = fmt.Sprintf("%sv%s/", host, version)
		}
	}
	return host
}
