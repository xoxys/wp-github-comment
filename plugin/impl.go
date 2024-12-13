package plugin

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/rs/zerolog/log"
	gh "github.com/thegeeklab/wp-github-comment/github"
	plugin_file "github.com/thegeeklab/wp-plugin-go/v4/file"
)

var ErrPluginEventNotSupported = errors.New("event not supported")

//nolint:revive
func (p *Plugin) run(ctx context.Context) error {
	if err := p.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if err := p.Execute(); err != nil {
		return fmt.Errorf("execution failed: %w", err)
	}

	return nil
}

// Validate handles the settings validation of the plugin.
func (p *Plugin) Validate() error {
	var err error

	if p.Metadata.Pipeline.Event != "pull_request" {
		return fmt.Errorf("%w: %s", ErrPluginEventNotSupported, p.Metadata.Pipeline.Event)
	}

	if p.Settings.Message != "" {
		if p.Settings.Message, p.Settings.IsFile, err = plugin_file.ReadStringOrFile(p.Settings.Message); err != nil {
			return fmt.Errorf("error while reading %s: %w", p.Settings.Message, err)
		}
	}

	if !strings.HasSuffix(p.Settings.BaseURL, "/") {
		p.Settings.BaseURL += "/"
	}

	p.Settings.baseURL, err = url.Parse(p.Settings.BaseURL)
	if err != nil {
		return fmt.Errorf("failed to parse base url: %w", err)
	}

	if p.Settings.Key == "" {
		key := fmt.Sprintf("%s/%s/%d", p.Metadata.Repository.Owner, p.Metadata.Repository.Name, p.Settings.IssueNum)
		hash := sha256.Sum256([]byte(key))
		p.Settings.Key = fmt.Sprintf("%x", hash)
	}

	if p.Settings.Key, _, err = plugin_file.ReadStringOrFile(p.Settings.Key); err != nil {
		return fmt.Errorf("error while reading %s: %w", p.Settings.Key, err)
	}

	return nil
}

// Execute provides the implementation of the plugin.
func (p *Plugin) Execute() error {
	client := gh.NewClient(p.Network.Context, p.Settings.baseURL, p.Settings.APIKey, p.Network.Client)
	client.Issue.Opt = gh.IssueOptions{
		Repo:    p.Metadata.Repository.Name,
		Owner:   p.Metadata.Repository.Owner,
		Message: p.Settings.Message,
		Update:  p.Settings.Update,
		Key:     p.Settings.Key,
		Number:  p.Metadata.Curr.PullRequest,
	}

	if p.Settings.SkipMissing && !p.Settings.IsFile {
		log.Info().
			Msg("comment skipped: 'message' is not a valid path or file does not exist while 'skip-missing' is enabled")

		return nil
	}

	_, err := client.Issue.AddComment(p.Network.Context)
	if err != nil {
		return fmt.Errorf("failed to create or update comment: %w", err)
	}

	return nil
}
