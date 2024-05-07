package plugin

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-github/v61/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/stretchr/testify/assert"
)

func TestGithubIssue_FindComment(t *testing.T) {
	ctx := context.Background()
	key := "test-key"
	keyPattern := "<!-- id: " + key + " -->"
	owner := "test-owner"
	repo := "test-repo"
	number := 1

	tests := []struct {
		name     string
		comments []*github.IssueComment
		want     *github.IssueComment
		wantErr  error
	}{
		{
			name:    "no comments",
			want:    nil,
			wantErr: ErrCommentNotFound,
		},
		{
			name: "comment found",
			comments: []*github.IssueComment{
				{Body: github.String(keyPattern)},
			},
			want: &github.IssueComment{Body: github.String(keyPattern)},
		},
		{
			name: "comment not found",
			comments: []*github.IssueComment{
				{Body: github.String("other comment")},
			},
			want:    nil,
			wantErr: ErrCommentNotFound,
		},
		{
			name: "multiple comments",
			comments: []*github.IssueComment{
				{Body: github.String("other comment")},
				{Body: github.String(keyPattern)},
				{Body: github.String("another comment")},
			},
			want: &github.IssueComment{Body: github.String(keyPattern)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockedHTTPClient := mock.NewMockedHTTPClient(
				mock.WithRequestMatch(
					mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
					tt.comments,
				),
			)

			client := github.NewClient(mockedHTTPClient)
			issue := &GithubIssue{
				Client: client,
				Owner:  owner,
				Repo:   repo,
				Number: number,
				Key:    key,
			}

			got, err := issue.FindComment(ctx)
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
	ctx := context.Background()
	key := "test-key"
	keyPattern := "<!-- id: " + key + " -->"
	owner := "test-owner"
	repo := "test-repo"
	number := 1
	message := "test message"

	tests := []struct {
		name        string
		update      bool
		existingKey string
		comments    []*github.IssueComment
		want        *github.IssueComment
		wantErr     bool
	}{
		{
			name:   "create new comment",
			update: false,
			want:   &github.IssueComment{Body: github.String("test message\n<!-- id: test-key -->\n")},
		},
		{
			name:   "update existing comment",
			update: true,
			comments: []*github.IssueComment{
				{ID: github.Int64(123), Body: github.String(keyPattern)},
			},
			want: &github.IssueComment{Body: github.String("test message\n<!-- id: test-key -->\n")},
		},
		{
			name:   "update non-existing comment",
			update: true,
			want:   &github.IssueComment{Body: github.String("test message\n<!-- id: test-key -->\n")},
		},
		{
			name:    "create new comment with error",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockedHTTPClient := mock.NewMockedHTTPClient(
				mock.WithRequestMatch(
					mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
					tt.comments,
				),
				mock.WithRequestMatchHandler(
					mock.PostReposIssuesCommentsByOwnerByRepoByIssueNumber,
					http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
						if tt.wantErr {
							mock.WriteError(w, http.StatusInternalServerError, "internal server error")
						} else {
							_, _ = w.Write(mock.MustMarshal(
								&github.IssueComment{
									Body: github.String(fmt.Sprintf("%s\n%s\n", message, keyPattern)),
								},
							))
						}
					}),
				),
				mock.WithRequestMatchHandler(
					mock.PatchReposIssuesCommentsByOwnerByRepoByCommentId,
					http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
						if tt.wantErr {
							mock.WriteError(w, http.StatusInternalServerError, "internal server error")
						} else {
							_, _ = w.Write(mock.MustMarshal(
								&github.IssueComment{
									Body: github.String(fmt.Sprintf("%s\n%s\n", message, keyPattern)),
								},
							))
						}
					}),
				),
			)

			client := github.NewClient(mockedHTTPClient)
			issue := &GithubIssue{
				Client:  client,
				Owner:   owner,
				Repo:    repo,
				Number:  number,
				Key:     key,
				Message: message,
				Update:  tt.update,
			}

			got, err := issue.AddComment(ctx)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.want, got)

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
