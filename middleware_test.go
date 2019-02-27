package soajsgo

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegistry_Middleware(t *testing.T) {
	tt := []struct {
		name            string
		headerInfo      string
		reg             Registry
		expectedSoaData ContextData
	}{
		{
			name:            "bad header",
			headerInfo:      "nil",
			reg:             Registry{},
			expectedSoaData: ContextData{},
		},
		{
			name:            "empty header",
			headerInfo:      "",
			reg:             Registry{},
			expectedSoaData: ContextData{},
		},
		{
			name:            "all ok",
			headerInfo:      `{"device":"iPhone"}`,
			reg:             Registry{Name: "ok"},
			expectedSoaData: ContextData{Device: "iPhone", Reg: &Registry{Name: "ok"}},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				soa := r.Context().Value(SoajsKey)
				if soa != nil {
					assert.Equal(t, tc.expectedSoaData, soa.(ContextData))
				} else {
					assert.Nil(t, soa)
				}
				_, _ = w.Write([]byte("ok"))
			})
			req := httptest.NewRequest("", "http://localhost:8080/", nil)
			req.Header.Set(headerDataName, tc.headerInfo)
			rec := httptest.NewRecorder()
			middleware := tc.reg.Middleware(handler)
			middleware.ServeHTTP(rec, req)
		})
	}
}

func TestHeaderData(t *testing.T) {
	tt := []struct {
		name         string
		data         string
		expectedInfo *headerInfo
		expectedErr  error
	}{
		{
			name:         "bad header",
			data:         "nil",
			expectedInfo: nil,
			expectedErr:  errors.New("unable to parse SOAJS header: invalid character 'i' in literal null (expecting 'u')"),
		},
		{
			name:         "all ok",
			data:         `{"device":"iPhone"}`,
			expectedInfo: &headerInfo{Device: "iPhone"},
			expectedErr:  nil,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("", "http://localhost:8080/", nil)
			req.Header.Set(headerDataName, tc.data)

			info, err := headerData(req)
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedInfo, info)
		})
	}
}

func TestHost_Path(t *testing.T) {
	tt := []struct {
		name         string
		host         Host
		args         []string
		expectedPath string
	}{
		{
			name: "1",
			host: Host{
				Host: "localhost",
				Port: 8080,
			},
			args:         []string{"test"},
			expectedPath: "localhost:8080/",
		},
		{
			name: "2",
			host: Host{
				Host: "localhost",
				Port: 8080,
			},
			args:         []string{"CONTROLLER", "v"},
			expectedPath: "localhost:8080/CONTROLLER/",
		},
		{
			name: "3",
			host: Host{
				Host: "localhost",
				Port: 8080,
			},
			args:         []string{"CONTROLLER", "1", "-"},
			expectedPath: "localhost:8080/CONTROLLER/v1/",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			p := tc.host.Path(tc.args...)
			assert.Equal(t, tc.expectedPath, p)
		})
	}
}
