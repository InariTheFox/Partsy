package utils

import (
	"net/url"
	"path"

	"github.com/inarithefox/partsy/server/public/model"
	"github.com/pkg/errors"
)

func GetSubpathFromConfig(config *model.Config) (string, error) {
	if config == nil {
		return "", errors.New("no configuration provided")
	} else if config.ServiceSettings.SiteURL == nil {
		return "/", nil
	}

	u, err := url.Parse(*config.ServiceSettings.SiteURL)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse SiteURL from configuration")
	}

	if u.Path == "" {
		return "/", nil
	}

	return path.Clean(u.Path), nil
}
