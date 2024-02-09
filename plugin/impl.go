package plugin

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/google/go-github/v59/github"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

var ErrPluginEventNotSupported = errors.New("event not supported")

//nolint:revive
func (p *Plugin) run(ctx context.Context) error {
	if err := p.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	//nolint:contextcheck
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
		if p.Settings.Message, p.Settings.IsFile, err = readStringOrFile(p.Settings.Message); err != nil {
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

	if p.Settings.Key, _, err = readStringOrFile(p.Settings.Key); err != nil {
		return fmt.Errorf("error while reading %s: %w", p.Settings.Key, err)
	}

	return nil
}

// Execute provides the implementation of the plugin.
func (p *Plugin) Execute() error {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: p.Settings.APIKey})
	tc := oauth2.NewClient(
		context.WithValue(p.Network.Context, oauth2.HTTPClient, p.Network.Client),
		ts,
	)

	client := github.NewClient(tc)
	client.BaseURL = p.Settings.baseURL

	commentClient := commentClient{
		Client:   client,
		Repo:     p.Metadata.Repository.Name,
		Owner:    p.Metadata.Repository.Owner,
		Message:  p.Settings.Message,
		Update:   p.Settings.Update,
		Key:      p.Settings.Key,
		IssueNum: p.Metadata.Curr.PullRequest,
	}

	if p.Settings.SkipMissing && !p.Settings.IsFile {
		log.Info().
			Msg("comment skipped: 'message' is not a valid path or file does not exist while 'skip-missing' is enabled")

		return nil
	}

	err := commentClient.issueComment(p.Network.Context)
	if err != nil {
		return fmt.Errorf("failed to create or update comment: %w", err)
	}

	return nil
}
