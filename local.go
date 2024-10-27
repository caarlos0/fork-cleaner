package forkcleaner

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/google/go-github/v50/github"
)

// LocalRepoState tracks the git status cleanliness, git stash cleanliness, and
// for each branch, to which remote branch it has been merged, if any.
type LocalRepoState struct {
	Path           string
	repo           *git.Repository
	Size           int64
	StatusClean    bool
	StashClean     bool
	MergedOrigin   map[string]string
	MergedPR       map[string]*github.PullRequest
	Unmerged       map[string]struct{}
	RemotesChecked []string
}

func NewLocalRepoState(path string, client *github.Client, ctx context.Context) (*LocalRepoState, error) {
	lr := LocalRepoState{
		Path:         path,
		MergedOrigin: make(map[string]string),
		MergedPR:     make(map[string]*github.PullRequest),
		Unmerged:     make(map[string]struct{}),
	}
	var err error

	lr.Size, err = getDiskSpaceUsed(path)
	if err != nil {
		return nil, err
	}
	lr.repo, err = git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	// 1) check status
	if err := lr.checkLocalStatus(); err != nil {
		return nil, err
	}

	// 2) check stash
	if err := lr.checkLocalStash(); err != nil {
		return nil, err
	}

	// 3) check branches
	if err := lr.checkLocalBranches(client, ctx); err != nil {
		return nil, err
	}

	return &lr, nil

}

func (lr *LocalRepoState) checkLocalStatus() error {
	// git-go's status sometimes reports dirty when it should be clean. e.g.:
	// https://github.com/go-git/go-git/issues/691
	// https://github.com/go-git/go-git/issues/736
	/*
		w, err := lr.repo.Worktree()
		if err != nil {
			return err
		}
		status, err := w.Status()
		if err != nil {
			return err
		}
		log.Println("Status for ", lr.Path, status)
		lr.StatusClean = status.IsClean()
	*/

	cmd := exec.Command("git", "status", "--porcelain")
	var out bytes.Buffer
	cmd.Dir = lr.Path
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	cmd.Run()
	lr.StatusClean = out.String() == ""
	return nil
}

func (lr *LocalRepoState) checkLocalStash() error {
	// stash is not supported yet in the go library, so we run the git command
	// https://github.com/go-git/go-git/issues/606

	cmd := exec.Command("git", "stash", "list")
	var out bytes.Buffer
	cmd.Dir = lr.Path
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	cmd.Run()
	lr.StashClean = out.String() == ""
	return nil
}

func (b *LocalRepoState) AddMerged(local, remote string, pr *github.PullRequest) {
	b.MergedOrigin[local] = remote
	b.MergedPR[local] = pr
}

func (b *LocalRepoState) AddUnmerged(local string) {
	b.Unmerged[local] = struct{}{}
}

func (b *LocalRepoState) Clean() bool {
	return len(b.Unmerged) == 0 && b.StatusClean && b.StashClean
}

// does the git repository have any commits that are not pushed to the remote?
func (lr *LocalRepoState) checkLocalBranches(client *github.Client, ctx context.Context) error {
	// first get the local branches. they have a name like refs/heads/<branch name> or use b.Name().Short()
	branches, err := lr.repo.Branches()
	if err != nil {
		return err
	}

	// then get the commits for each branch and check if the commit is in the remote
	// it could be in a branch with the same name, or in a branch with a different name
	err = branches.ForEach(func(b *plumbing.Reference) error {
		var remotesFound int

		for _, remName := range []string{"origin", "upstream"} {
			rem, err := lr.repo.Remote(remName)
			if err == git.ErrRemoteNotFound {
				continue
			}
			if err != nil {
				return err
			}
			remotesFound++
			var found bool
			for _, existing := range lr.RemotesChecked {
				if existing == rem.Config().URLs[0] {
					found = true
				}
			}
			if !found {
				lr.RemotesChecked = append(lr.RemotesChecked, rem.Config().URLs[0])
			}

			found, pr, err := isCommitInRemote(ctx, client, rem, b.Hash())
			if err != nil {
				// if it's http 404, just continue
				if strings.Contains(err.Error(), "404 Not Found") {
					continue
				}
				// can't use an invalid or unsupported URL.
				// this usually means you use a host other than GitHub. these are not supported yet
				if strings.Contains(err.Error(), "unsupported remote URL") {
					continue
				}

				return err
			}
			if found { // note: pr might be nil if it was committed without a PR
				lr.AddMerged(b.Name().Short(), remName, pr)
				return nil
			}
		}
		lr.AddUnmerged(b.Name().Short())
		return nil
	})

	if err != nil {
		return fmt.Errorf("error while iterating over branches: %w", err)
	}

	return nil
}

func isCommitInRemote(ctx context.Context, client *github.Client, rem *git.Remote, commit plumbing.Hash) (bool, *github.PullRequest, error) {
	remoteUrl := rem.Config().URLs[0]

	owner, name, err := extractOwnerAndNameFromRemoteUrl(remoteUrl)
	if err != nil {
		return false, nil, err
	}
	opts := github.PullRequestListOptions{
		State: "closed",
	}
	prs, _, err := client.PullRequests.ListPullRequestsWithCommit(ctx, owner, name, commit.String(), &opts)
	if err != nil {
		if strings.Contains(err.Error(), "No commit found for SHA") {
			return false, nil, nil
		}
		return false, nil, err
	}

	if len(prs) > 0 {
		return true, prs[0], nil
	}
	// important! we are here because a commit was committed directly (without a PR)
	return true, nil, nil
}

// extractOwnerAndNameFromRemoteUrl extracts the owner and name from a remote URL
// the remoteURL is in the form of either git@github.com:<owner>/<name>.git or https://github.com/<owner>/<name>.git
func extractOwnerAndNameFromRemoteUrl(remoteUrl string) (string, string, error) {
	str := strings.TrimSuffix(remoteUrl, ".git")
	str = strings.TrimPrefix(str, "git@github.com:")
	str = strings.TrimPrefix(str, "https://github.com/")
	str = strings.TrimPrefix(str, "http://github.com/")
	str = strings.TrimPrefix(str, "git://github.com/")
	split := strings.Split(str, "/")
	if len(split) != 2 {
		return "", "", fmt.Errorf("unsupported remote URL: %s", remoteUrl)
	}
	return split[0], split[1], nil
}

// getDiskSpaceUsed calculates the amount of disk space used under `path`
// to achieve this we must fully, recursively walk the directory tree.
// note: this is not super accurate, we don't account for the size of the directories themselves
func getDiskSpaceUsed(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(base string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			//log.Println("info.Size():", info.Size(), path, base)
			size += info.Size()
		}
		return nil
	})
	//log.Println("total size", size)
	return size, err
}
