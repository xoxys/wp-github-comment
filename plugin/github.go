package plugin

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v61/github"
)

type GithubClient struct {
	Client *github.Client
	Issue  *GithubIssue
}

type GithubIssue struct {
	*github.Client
	Number  int
	Message string
	Key     string
	Repo    string
	Owner   string
	Update  bool
}

// Constructor function for Parent.
func NewGithubClient(client *github.Client) *GithubClient {
	return &GithubClient{
		Client: client,
		Issue:  &GithubIssue{Client: client},
	}
}

// AddComment adds a new comment or updates an existing comment on a GitHub issue.
// If the Update field is true, it will append a unique identifier to the comment
// body and attempt to find and update the existing comment with that identifier.
// Otherwise, it will create a new comment on the issue.
func (i *GithubIssue) AddComment(ctx context.Context) error {
	var (
		err     error
		comment *github.IssueComment
		resp    *github.Response
	)

	issueComment := &github.IssueComment{
		Body: &i.Message,
	}

	if i.Update {
		// Append plugin comment ID to comment message so we can search for it later
		message := fmt.Sprintf("%s\n<!-- id: %s -->\n", i.Message, i.Key)
		issueComment.Body = &message

		comment, err = i.FindComment(ctx)

		if err == nil && comment != nil {
			_, resp, err = i.Client.Issues.EditComment(ctx, i.Owner, i.Repo, *comment.ID, issueComment)
		}
	}

	if err == nil && resp == nil {
		_, _, err = i.Client.Issues.CreateComment(ctx, i.Owner, i.Repo, i.Number, issueComment)
	}

	if err != nil {
		return err
	}

	return nil
}

// FindComment returns the GitHub issue comment that contains the specified key, or nil if no such comment exists.
// It retrieves all comments on the issue and searches for one that contains the specified key in the comment body.
func (i *GithubIssue) FindComment(ctx context.Context) (*github.IssueComment, error) {
	var allComments []*github.IssueComment

	opts := &github.IssueListCommentsOptions{}

	for {
		comments, resp, err := i.Client.Issues.ListComments(ctx, i.Owner, i.Repo, i.Number, opts)
		if err != nil {
			return nil, err
		}

		allComments = append(allComments, comments...)

		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	for _, comment := range allComments {
		if strings.Contains(*comment.Body, fmt.Sprintf("<!-- id: %s -->", i.Key)) {
			return comment, nil
		}
	}

	//nolint:nilnil
	return nil, nil
}
