package soajsgo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type registryPath string

func registryAddress() (registryPath, error) {
	registryAPI := os.Getenv(EnvRegistryAPIAddress)
	if registryAPI == "" {
		return "", fmt.Errorf("could not find environment variable %s", EnvRegistryAPIAddress)
	}
	if index := strings.Index(registryAPI, ":"); index == -1 {
		return "", fmt.Errorf("invalid format for %s. Got [%s], expected [hostname:port]", EnvRegistryAPIAddress, registryAPI)
	}
	port := strings.Split(registryAPI, ":")[1]
	if _, err := strconv.Atoi(port); err != nil {
		return "", fmt.Errorf("port must be an integer, got %s", port)
	}
	return registryPath(registryAPI), nil
}

func (r registryPath) register() string {
	return fmt.Sprintf("http://%s/register", r)
}

func (r registryPath) getRegistry(serviceName, envCode string) string {
	return fmt.Sprintf("http://%s/getRegistry?env=%s&serviceName=%s", r, envCode, serviceName)
}

func registryResponse(res *http.Response) (*Registry, error) {
	if res.StatusCode < 200 || res.StatusCode > 299 {
		b, _ := ioutil.ReadAll(res.Body)
		return nil, fmt.Errorf("non 2xx status code: %d %s", res.StatusCode, b)
	}
	var regRes registryAPIResponse
	err := json.NewDecoder(res.Body).Decode(&regRes)
	if err != nil {
		return nil, fmt.Errorf("could not decode registry response: %v", err)
	}
	if len(regRes.Errors.Details) > 0 {
		return nil, fmt.Errorf("unable to register service at gateway: [%d] [%s]",
			regRes.Errors.Details[0].Code,
			regRes.Errors.Details[0].Message)
	}
	if !regRes.Result {
		return nil, errors.New("negative result by registry")
	}
	return &regRes.Registry, nil
}
