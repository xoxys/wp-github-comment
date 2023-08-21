package plugin

import (
	"net/url"

	wp "github.com/thegeeklab/wp-plugin-go/plugin"
)

// Plugin implements provide the plugin.
type Plugin struct {
	*wp.Plugin
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

func New(options wp.Options, settings *Settings) *Plugin {
	p := &Plugin{}

	options.Execute = p.run

	p.Plugin = wp.New(options)
	p.Settings = settings

	return p
}
