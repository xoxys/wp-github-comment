package github

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-github/v63/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thegeeklab/wp-github-comment/github/mocks"
)

var ErrInternalServerError = errors.New("internal server error")

func TestGithubIssue_FindComment(t *testing.T) {
	tests := []struct {
		name     string
		issueOpt IssueOptions
		comments []*github.IssueComment
		want     *github.IssueComment
		wantErr  error
	}{
		{
			name: "no comments",
			issueOpt: IssueOptions{
				Key:   "test-key",
				Owner: "test-owner",
				Repo:  "test-repo",
			},
			wantErr: ErrCommentNotFound,
		},
		{
			name: "comment found",
			issueOpt: IssueOptions{
				Key:   "test-key",
				Owner: "test-owner",
				Repo:  "test-repo",
			},
			comments: []*github.IssueComment{
				{Body: github.String("<!-- id: test-key -->\ntest comment\n")},
			},
			want: &github.IssueComment{
				Body: github.String("<!-- id: test-key -->\ntest comment\n"),
			},
		},
		{
			name: "comment not found",
			issueOpt: IssueOptions{
				Key:   "test-key",
				Owner: "test-owner",
				Repo:  "test-repo",
			},
			comments: []*github.IssueComment{
				{Body: github.String("other comment")},
			},
			wantErr: ErrCommentNotFound,
		},
		{
			name: "multiple comments",
			issueOpt: IssueOptions{
				Key:   "test-key",
				Owner: "test-owner",
				Repo:  "test-repo",
			},
			comments: []*github.IssueComment{
				{Body: github.String("other comment")},
				{Body: github.String("<!-- id: test-key -->\ntest comment\n")},
				{Body: github.String("another comment")},
			},
			want: &github.IssueComment{Body: github.String("<!-- id: test-key -->\ntest comment\n")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := mocks.NewMockIssueService(t)
			issue := &Issue{
				client: mockClient,
				Opt:    tt.issueOpt,
			}

			mockClient.
				On("ListComments", mock.Anything, tt.issueOpt.Owner, tt.issueOpt.Repo, mock.Anything, mock.Anything).
				Return(tt.comments, nil, nil)

			got, err := issue.FindComment(context.Background())
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.want, got)

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGithubIssue_AddComment(t *testing.T) {
	tests := []struct {
		name     string
		issueOpt IssueOptions
		comments []*github.IssueComment
		want     *github.IssueComment
		wantErr  error
	}{
		{
			name: "create new comment",
			issueOpt: IssueOptions{
				Key:     "test-key",
				Owner:   "test-owner",
				Repo:    "test-repo",
				Message: "test message",
				Update:  false,
			},
			want: &github.IssueComment{
				Body: github.String("<!-- id: test-key -->\ntest message\n"),
			},
		},
		{
			name: "update existing comment",
			issueOpt: IssueOptions{
				Key:     "test-key",
				Owner:   "test-owner",
				Repo:    "test-repo",
				Message: "test message",
				Update:  true,
			},
			comments: []*github.IssueComment{
				{ID: github.Int64(123), Body: github.String("<!-- id: test-key -->\ntest message\n")},
			},
			want: &github.IssueComment{
				Body: github.String("<!-- id: test-key -->\ntest message\n"),
			},
		},
		{
			name: "update non-existing comment",
			issueOpt: IssueOptions{
				Key:     "test-key",
				Owner:   "test-owner",
				Repo:    "test-repo",
				Message: "test message",
				Update:  true,
			},
			want: &github.IssueComment{
				Body: github.String("<!-- id: test-key -->\ntest message\n"),
			},
		},
		{
			name: "create new comment with error",
			issueOpt: IssueOptions{
				Key:     "test-key",
				Owner:   "test-owner",
				Repo:    "test-repo",
				Message: "test message",
				Update:  false,
			},
			wantErr: ErrInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := mocks.NewMockIssueService(t)
			issue := &Issue{
				client: mockClient,
				Opt:    tt.issueOpt,
			}

			if tt.issueOpt.Update {
				mockClient.
					On("ListComments", mock.Anything, tt.issueOpt.Owner, tt.issueOpt.Repo, mock.Anything, mock.Anything).
					Return(tt.comments, nil, nil)
			}

			if tt.issueOpt.Update && tt.comments != nil {
				mockClient.
					On("EditComment", mock.Anything, tt.issueOpt.Owner, tt.issueOpt.Repo, mock.Anything, mock.Anything).
					Return(&github.IssueComment{
						Body: github.String(fmt.Sprintf("<!-- id: %s -->\n%s\n", tt.issueOpt.Key, tt.issueOpt.Message)),
					}, nil, nil)
			}

			if tt.comments == nil {
				var comment *github.IssueComment
				if tt.wantErr == nil {
					comment = &github.IssueComment{
						Body: github.String(fmt.Sprintf("<!-- id: %s -->\n%s\n", tt.issueOpt.Key, tt.issueOpt.Message)),
					}
				}

				mockClient.
					On("CreateComment", mock.Anything, tt.issueOpt.Owner, tt.issueOpt.Repo, mock.Anything, mock.Anything).
					Return(comment, nil, tt.wantErr)
			}

			got, err := issue.AddComment(context.Background())
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.want, got)

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
