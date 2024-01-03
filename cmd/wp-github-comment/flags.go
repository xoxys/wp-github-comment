package main

import (
	"github.com/thegeeklab/wp-github-comment/plugin"
	"github.com/urfave/cli/v2"
)

// settingsFlags has the cli.Flags for the plugin.Settings.
//
//go:generate go run docs.go flags.go
func settingsFlags(settings *plugin.Settings, category string) []cli.Flag {
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
