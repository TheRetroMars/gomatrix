package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	configDirName   = ".config"
	configAppName   = "gomatrix"
	configFileName  = "config.json"
	localConfigName = "gomatrix.json"
)

// Config holds user-configurable options.
type Config struct {
	Speed        int    `json:"speed"`
	ColorTheme   string `json:"color_theme"`
	ASCIIOnly    bool   `json:"ascii_only"`
	StarMode     bool   `json:"star_mode"`
	StarCount    int    `json:"star_count"`
	CRTMode      bool   `json:"crt_mode"`
	T1Percent    int    `json:"t1_percent"`
	T2Percent    int    `json:"t2_percent"`
	TrueColor    bool   `json:"true_color"`
	GradientMode int    `json:"gradient_mode"`
}

// loadConfig reads the config.json file and parses it.
func loadConfig(cliPath string) (*Config, error) {
	cfg := &Config{
		Speed:        DefaultSpeed,
		ColorTheme:   "Classic",
		ASCIIOnly:    false,
		StarMode:     false,
		StarCount:    DefaultStarCnt,
		CRTMode:      false,
		T1Percent:    DefaultT1Percent,
		T2Percent:    DefaultT2Percent,
		TrueColor:    true,
		GradientMode: ModeSmooth,
	}
	paths := []string{cliPath}
	home, err := os.UserHomeDir()
	if err == nil {
		globalCfg := filepath.Join(
			home, configDirName, configAppName, configFileName)
		paths = append(paths, globalCfg)
	}
	paths = append(paths, localConfigName)

	var data []byte
	for _, p := range paths {
		if p == "" {
			continue
		}
		if b, err := os.ReadFile(p); err == nil {
			data = b
			break
		}
	}
	if len(data) > 0 {
		if err := json.Unmarshal(data, cfg); err != nil {
			return nil, err
		}
	}

	cfg.Speed = clampSpeed(cfg.Speed)
	sanitizeConfig(cfg)

	return cfg, nil
}

// sanitizeConfig validates configuration parameters and applies defaults.
func sanitizeConfig(cfg *Config) {
	if cfg.T1Percent < 0 || cfg.T1Percent > int(MaxPercent) {
		cfg.T1Percent = DefaultT1Percent
	}
	if cfg.T2Percent < 0 || cfg.T2Percent > int(MaxPercent) {
		cfg.T2Percent = DefaultT2Percent
	}
	if cfg.T1Percent >= cfg.T2Percent {
		cfg.T1Percent = DefaultT1Percent
		cfg.T2Percent = DefaultT2Percent
	}
	if cfg.GradientMode < ModeClassic || cfg.GradientMode > ModeVerySmooth {
		cfg.GradientMode = ModeSmooth
	}
}

// clampSpeed constrains speed to [MinSpeed, len(SpeedDelays)].
func clampSpeed(s int) int {
	maxSpd := len(SpeedDelays)
	if s < MinSpeed {
		return MinSpeed
	}
	if s > maxSpd {
		return maxSpd
	}
	return s
}
