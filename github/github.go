package github

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-github/v68/github"
	"golang.org/x/oauth2"
)

var ErrCommentNotFound = errors.New("comment not found")

type Client struct {
	client *github.Client
	Issue  *Issue
}

type Issue struct {
	client IssueService
	Opt    IssueOptions
}

type IssueOptions struct {
	Number  int
	Message string
	Key     string
	Repo    string
	Owner   string
	Update  bool
}

// NewGitHubClient creates a new GitHubClient instance that wraps the provided GitHub API client.
// The GitHubClient provides a higher-level interface for interacting with the GitHub API,
// including methods for managing GitHub issues.
func NewClient(ctx context.Context, url *url.URL, token string, client *http.Client) *Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(
		context.WithValue(ctx, oauth2.HTTPClient, client),
		ts,
	)

	c := github.NewClient(tc)
	c.BaseURL = url

	return &Client{
		client: c,
		Issue: &Issue{
			client: &IssueServiceImpl{client: c},
			Opt:    IssueOptions{},
		},
	}
}

// AddComment adds a new comment or updates an existing comment on a GitHub issue.
// If the Update field is true, it will append a unique identifier to the comment
// body and attempt to find and update the existing comment with that identifier.
// Otherwise, it will create a new comment on the issue.
func (i *Issue) AddComment(ctx context.Context) (*github.IssueComment, error) {
	issueComment := &github.IssueComment{
		Body: &i.Opt.Message,
	}

	if i.Opt.Update {
		// Append plugin comment ID to comment message so we can search for it later
		*issueComment.Body = fmt.Sprintf("%s\n<!-- id: %s -->\n", i.Opt.Message, i.Opt.Key)

		comment, err := i.FindComment(ctx)
		if err != nil && !errors.Is(err, ErrCommentNotFound) {
			return nil, err
		}

		if comment != nil {
			comment, _, err = i.client.EditComment(ctx, i.Opt.Owner, i.Opt.Repo, *comment.ID, issueComment)

			return comment, err
		}
	}

	comment, _, err := i.client.CreateComment(ctx, i.Opt.Owner, i.Opt.Repo, i.Opt.Number, issueComment)

	return comment, err
}

// FindComment returns the GitHub issue comment that contains the specified key, or nil if no such comment exists.
// It retrieves all comments on the issue and searches for one that contains the specified key in the comment body.
func (i *Issue) FindComment(ctx context.Context) (*github.IssueComment, error) {
	var allComments []*github.IssueComment

	opts := &github.IssueListCommentsOptions{}

	for {
		comments, resp, err := i.client.ListComments(ctx, i.Opt.Owner, i.Opt.Repo, i.Opt.Number, opts)
		if err != nil {
			return nil, err
		}

		allComments = append(allComments, comments...)

		if resp == nil || resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	for _, comment := range allComments {
		if strings.Contains(*comment.Body, fmt.Sprintf("<!-- id: %s -->", i.Opt.Key)) {
			return comment, nil
		}
	}

	return nil, fmt.Errorf("%w: failed to find comment with key %s", ErrCommentNotFound, i.Opt.Key)
}
