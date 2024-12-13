package plugin

import (
	"fmt"
	"net/url"

	plugin_base "github.com/thegeeklab/wp-plugin-go/v4/plugin"
	"github.com/urfave/cli/v2"
)

//go:generate go run ../internal/doc/main.go -output=../docs/data/data-raw.yaml

// Plugin implements provide the plugin.
type Plugin struct {
	*plugin_base.Plugin
	Settings *Settings
}

// Settings for the Plugin.
type Settings struct {
	BaseURL     string
	IssueNum    int
	Key         string
	Message     string
	Update      bool
	APIKey      string
	SkipMissing bool
	IsFile      bool

	baseURL *url.URL
}

func New(e plugin_base.ExecuteFunc, build ...string) *Plugin {
	p := &Plugin{
		Settings: &Settings{},
	}

	options := plugin_base.Options{
		Name:                "wp-github-comment",
		Description:         "Add comments to GitHub Issues and Pull Requests",
		Flags:               Flags(p.Settings, plugin_base.FlagsPluginCategory),
		Execute:             p.run,
		HideWoodpeckerFlags: true,
	}

	if len(build) > 0 {
		options.Version = build[0]
	}

	if len(build) > 1 {
		options.VersionMetadata = fmt.Sprintf("date=%s", build[1])
	}

	if e != nil {
		options.Execute = e
	}

	p.Plugin = plugin_base.New(options)

	return p
}

// Flags returns a slice of CLI flags for the plugin.
func Flags(settings *Settings, category string) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "api-key",
			EnvVars:     []string{"PLUGIN_API_KEY", "GITHUB_COMMENT_API_KEY"},
			Usage:       "personal access token to access the GitHub API",
			Destination: &settings.APIKey,
			Category:    category,
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "base-url",
			EnvVars:     []string{"PLUGIN_BASE_URL", "GITHUB_COMMENT_BASE_URL"},
			Usage:       "API URL",
			Value:       "https://api.github.com/",
			Destination: &settings.BaseURL,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "key",
			EnvVars:     []string{"PLUGIN_KEY", "GITHUB_COMMENT_KEY"},
			Usage:       "unique identifier to assign to a comment",
			Destination: &settings.Key,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "message",
			EnvVars:     []string{"PLUGIN_MESSAGE", "GITHUB_COMMENT_MESSAGE"},
			Usage:       "path to file or string that contains the comment text",
			Destination: &settings.Message,
			Category:    category,
			Required:    true,
		},
		&cli.BoolFlag{
			Name:        "update",
			EnvVars:     []string{"PLUGIN_UPDATE", "GITHUB_COMMENT_UPDATE"},
			Usage:       "enable update of an existing comment that matches the key",
			Value:       false,
			Destination: &settings.Update,
			Category:    category,
		},
		&cli.BoolFlag{
			Name:        "skip-missing",
			EnvVars:     []string{"PLUGIN_SKIP_MISSING", "GITHUB_COMMENT_SKIP_MISSING"},
			Usage:       "skip comment creation if the given message file does not exist",
			Value:       false,
			Destination: &settings.SkipMissing,
			Category:    category,
		},
	}
}
