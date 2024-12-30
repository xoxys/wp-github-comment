package github

import (
	"context"

	"github.com/google/go-github/v68/github"
)

// APIClient is an interface that wraps the GitHub API client.
//
//nolint:lll
type IssueService interface {
	CreateComment(ctx context.Context, owner, repo string, number int, comment *github.IssueComment) (*github.IssueComment, *github.Response, error)
	EditComment(ctx context.Context, owner, repo string, commentID int64, comment *github.IssueComment) (*github.IssueComment, *github.Response, error)
	ListComments(ctx context.Context, owner, repo string, number int, opts *github.IssueListCommentsOptions) ([]*github.IssueComment, *github.Response, error)
}

type IssueServiceImpl struct {
	client *github.Client
}

// CreateComment wraps the CreateComment method of the github.IssuesService.
//
//nolint:lll
func (s *IssueServiceImpl) CreateComment(ctx context.Context, owner, repo string, number int, comment *github.IssueComment) (*github.IssueComment, *github.Response, error) {
	return s.client.Issues.CreateComment(ctx, owner, repo, number, comment)
}

// EditComment wraps the EditComment method of the github.IssuesService.
//
//nolint:lll
func (s *IssueServiceImpl) EditComment(ctx context.Context, owner, repo string, commentID int64, comment *github.IssueComment) (*github.IssueComment, *github.Response, error) {
	return s.client.Issues.EditComment(ctx, owner, repo, commentID, comment)
}

// ListComments wraps the ListComments method of the github.IssuesService.
//
//nolint:lll
func (s *IssueServiceImpl) ListComments(ctx context.Context, owner, repo string, number int, opts *github.IssueListCommentsOptions) ([]*github.IssueComment, *github.Response, error) {
	return s.client.Issues.ListComments(ctx, owner, repo, number, opts)
}
