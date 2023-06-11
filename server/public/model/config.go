package model

import (
	"encoding/json"
	"io"
)

const (
	ConnectionSecurityNone     = ""
	ConnectionSecurityPlain    = "PLAIN"
	ConnectionSecurityTLS      = "TLS"
	ConnectionSecuritySTARTTLS = "STARTTLS"

	ServiceSettingsDefaultSiteURL       = "http://localhost:8065"
	ServiceSettingsDefaultTLSCertFile   = ""
	ServiceSettingsDefaultTLSKeyFile    = ""
	ServiceSettingsDefaultListenAddress = ":8065"
	ServiceSettingsDefaultAllowCorsFrom = ""
	ServiceSettingsDefaultReadTimeout   = 300
	ServiceSettingsDefaultWriteTimeout  = 300
	ServiceSettingsDefaultIdleTimeout   = 60
)

type ServiceSettings struct {
	SiteURL                  *string `access:"environment_web_server,authentication,write_restrictable"`
	ListenAddress            *string `access:"environment_web_server,write_restrictable"`
	ConnectionSecurity       *string `access:"environment_web_server,write_restrictable"`
	TLSCertFile              *string `access:"environment_web_server,write_restrictable"`
	TLSKeyFile               *string `access:"enviroment_web_server,write_restrictable"`
	TLSMinVer                *string `access:"write_restrictable"`
	TLSStrictTransport       *bool   `access:"write_restrictable"`
	TLSStrictTransportMaxAge *int64  `access:"write_restrictable"`
	UseLetsEncrypt           *bool   `access:"environment_web_server,write_restrictable"`
	LetsEncryptCertCacheFile *string `access:"environment_web_server,write_restrictable"`
	ReadTimeout              *int    `access:"environment_web_server,write_restrictable"`
	WriteTimeout             *int    `access:"environment_web_server,write_restrictable"`
	IdleTimeout              *int    `access:"environment_web_server,write_restrictable"`
	ForwardHTTPToHTTPS       *bool   `access:"environment_web_server,write_restrictable"`
	AllowCorsFrom            *string `access:"write_restrictable"`
	EnableDeveloper          *bool   `access:"environment_developer,write_restrictable"`
}

func (s *ServiceSettings) SetDefaults() {
	if s.SiteURL == nil {
		if s.EnableDeveloper != nil && *s.EnableDeveloper {
			s.SiteURL = NewString(ServiceSettingsDefaultSiteURL)
		} else {
			s.SiteURL = NewString("")
		}
	}

	if s.ListenAddress == nil {
		s.ListenAddress = NewString(ServiceSettingsDefaultListenAddress)
	}

	if s.EnableDeveloper == nil {
		s.EnableDeveloper = NewBool(false)
	}

	if s.ConnectionSecurity == nil {
		s.ConnectionSecurity = NewString("")
	}

	if s.TLSKeyFile == nil {
		s.TLSKeyFile = NewString(ServiceSettingsDefaultTLSKeyFile)
	}

	if s.TLSCertFile == nil {
		s.TLSCertFile = NewString(ServiceSettingsDefaultTLSCertFile)
	}

	if s.TLSMinVer == nil {
		s.TLSMinVer = NewString("1.2")
	}

	if s.TLSStrictTransport == nil {
		s.TLSStrictTransport = NewBool(false)
	}

	if s.TLSStrictTransportMaxAge == nil {
		s.TLSStrictTransportMaxAge = NewInt64(63072000)
	}

	if s.UseLetsEncrypt == nil {
		s.UseLetsEncrypt = NewBool(false)
	}

	if s.LetsEncryptCertCacheFile == nil {
		s.LetsEncryptCertCacheFile = NewString("./config/letsencrypt.cache")
	}

	if s.ReadTimeout == nil {
		s.ReadTimeout = NewInt(ServiceSettingsDefaultReadTimeout)
	}

	if s.WriteTimeout == nil {
		s.WriteTimeout = NewInt(ServiceSettingsDefaultWriteTimeout)
	}

	if s.IdleTimeout == nil {
		s.IdleTimeout = NewInt(ServiceSettingsDefaultIdleTimeout)
	}

	if s.ForwardHTTPToHTTPS == nil {
		s.ForwardHTTPToHTTPS = NewBool(false)
	}
}

type Config struct {
	ServiceSettings ServiceSettings
}

func (o *Config) Clone() *Config {
	buf, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	var ret Config
	err = json.Unmarshal(buf, &ret)
	if err != nil {
		panic(err)
	}
	return &ret
}

func ConfigFromJSON(data io.Reader) *Config {
	var o *Config
	json.NewDecoder(data).Decode(&o)
	return o
}

func (o *Config) SetDefaults() {
	o.ServiceSettings.SetDefaults()
}
