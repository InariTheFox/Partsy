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

type Config struct {
	ServiceSettings   ServiceSettings
	SqlSettings       SqlSettings
	FileSettings      FileSettings
	EmailSettings     EmailSettings
	RateLimitSettings RateLimitSettings
	OpenIdSettings    OpenIdSettings
}

type ServiceSettings struct {
	SiteURL                  *string `access:"environment_web_server,authentication,write_restrictable"`
	ListenAddress            *string `access:"environment_web_server,write_restrictable"`
	ConnectionSecurity       *string `access:"environment_web_server,write_restrictable"`
	TLSCertFile              *string `access:"environment_web_server,write_restrictable"`
	TLSKeyFile               *string `access:"enviroment_web_server,write_restrictable"`
	TLSMinVer                *string `access:"environment_web_server,write_restrictable"`
	TLSStrictTransport       *bool   `access:"environment_web_server,write_restrictable"`
	TLSStrictTransportMaxAge *int64  `access:"environment_web_server,write_restrictable"`
	UseLetsEncrypt           *bool   `access:"environment_web_server,write_restrictable"`
	LetsEncryptCertCacheFile *string `access:"environment_web_server,write_restrictable"`
	ReadTimeout              *int    `access:"environment_web_server,write_restrictable"`
	WriteTimeout             *int    `access:"environment_web_server,write_restrictable"`
	IdleTimeout              *int    `access:"environment_web_server,write_restrictable"`
	ForwardHTTPToHTTPS       *bool   `access:"environment_web_server,write_restrictable"`
	AllowCorsFrom            *string `access:"environment_web_server,write_restrictable"`
	EnableDeveloper          *bool   `access:"environment_developer,write_restrictable"`
}

type SqlSettings struct {
	DriverName            *string  `access:"environment_database,write_restrictable"`
	DataSource            *string  `access:"environment_database,write_restrictable"`
	DataSourceReplicas    []string `access:"environment_database,write_resitrctable"`
	MaxOpenConnections    *int     `access:"environment_database,write_restrictable"`
	ConnectionMaxIdleTime *int     `access:"environment_database,write_restrictable"`
	ConnectionMaxLifetime *int     `access:"environment_database,write_restrictable"`
	QueryTimeout          *int     `access:"environment_database,write_restrictable"`
}

type FileSettings struct {
	EnableFileAttachments *bool   `access:"environment_file_storage,write_restrictable"`
	MaxFileSize           *int64  `access:"environment_file_storage,write_restrictable"`
	Directory             *string `access:"environment_file_storage,write_restrictable"`
}

type EmailSettings struct {
	Enable                            *bool   `access:"site_notification"`
	ReplyToAddress                    *string `access:"site_notification,write_restrictable"`
	SMTPUsername                      *string `access:"environment_smtp,write_restrictable"`
	SMTPPassword                      *string `access:"environment_smtp,write_restrictable"`
	SMTPPort                          *int    `access:"environment_smtp,write_restrictable"`
	SMTPServerTimeout                 *int    `access:"environment_smtp,write_restrictable"`
	ConnectionSecurity                *string `access:"environment_smtp,write_restrictable"`
	SkipServerCertificateVerification *bool   `access:"environment_smtp,write_restrictable"`
}

type RateLimitSettings struct {
	Enable *bool `access:"environment_rate_limit,write_restrictable"`
	Rate   *int  `access:"environment_rate_limit,write_restrictable"`
}

type OpenIdSettings struct {
	Enable            *bool   `access:"authentication,environment_openid,write_restrictable"`
	Secret            *string `access:"authentication,environment_openid,write_restrictable"`
	Id                *string `access:"authentication,environment_openid,write_restrictable"`
	Scope             *string `access:"authentication,environment_openid,write_restrictable"`
	AuthEndpoint      *string `access:"authentication,environment_openid,write_restrictable"`
	TokenEndpoint     *string `access:"authentication,environment_openid,write_restrictable"`
	DiscoveryEndpoint *string `access:"authentication,environment_openid,write_restrictable"`
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

func (s *SqlSettings) SetDefaults() {

}

func (s *FileSettings) SetDefaults() {
	if s.EnableFileAttachments == nil {
		s.EnableFileAttachments = NewBool(false)
	}
}

func (s *EmailSettings) SetDefaults() {
	if s.Enable == nil {
		s.Enable = NewBool(false)
	}

	if s.SkipServerCertificateVerification == nil {
		s.SkipServerCertificateVerification = NewBool(false)
	}
}

func (s *RateLimitSettings) SetDefaults() {
	if s.Enable == nil {
		s.Enable = NewBool(false)
	}
}

func (s *OpenIdSettings) SetDefaults() {
	if s.Enable == nil {
		s.Enable = NewBool(false)
	}

	if s.Secret == nil {
		s.Secret = NewString("")
	}

	if s.Id == nil {
		s.Id = NewString("")
	}

	if s.AuthEndpoint == nil {
		s.AuthEndpoint = NewString("")
	}

	if s.TokenEndpoint == nil {
		s.TokenEndpoint = NewString("")
	}

	if s.DiscoveryEndpoint == nil {
		s.DiscoveryEndpoint = NewString("")
	}
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
	o.SqlSettings.SetDefaults()
	o.FileSettings.SetDefaults()
	o.EmailSettings.SetDefaults()
	o.RateLimitSettings.SetDefaults()
	o.OpenIdSettings.SetDefaults()
}
