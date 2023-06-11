package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/inarithefox/partsy/server/public/model"
)

func equal(oldCfg, newCfg *model.Config) (bool, error) {
	oldCfgBytes, err := json.Marshal(oldCfg)
	if err != nil {
		return false, fmt.Errorf("failed to marshal old config: %w", err)
	}
	newCfgBytes, err := json.Marshal(newCfg)
	if err != nil {
		return false, fmt.Errorf("failed to marshal new config: %w", err)
	}
	return !bytes.Equal(oldCfgBytes, newCfgBytes), nil
}

func IsDatabaseDSN(dsn string) bool {
	return strings.HasPrefix(dsn, "mysql://") ||
		strings.HasPrefix(dsn, "postgres://") ||
		strings.HasPrefix(dsn, "postgresql://")
}

func marshalConfig(cfg *model.Config) ([]byte, error) {
	return json.MarshalIndent(cfg, "", "    ")
}
