package forkcleaner

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/google/go-github/v50/github"
)

func IsClean(ctx context.Context, path string, client *github.Client) (bool, error) {

	r, err := git.PlainOpen(path)
	if err != nil {
		return false, err
	}

	// 1) check status
	if !isLocalStatusClean(r) {
		return false, nil
	}

	// 2) check stash
	if !isLocalStashClean(path) {
		return false, nil
	}

	// 3) check branches
	ok, bms, err := isLocalBranchesClean(ctx, r, client)
	if err != nil {
		return false, err
	}
	bms.Dump()

	return ok, nil
}

func isLocalStatusClean(r *git.Repository) bool {
	w, err := r.Worktree()
	if err != nil {
		panic(err)
	}
	status, err := w.Status()
	if err != nil {
		panic(err)
	}
	fmt.Println("status clean", status.IsClean()) // WORKS
	return status.IsClean()
}

func isLocalStashClean(path string) bool {
	// stash is not supported yet in the go library, so we run the git command
	// https://github.com/go-git/go-git/issues/606

	cmd := exec.Command("git", "stash", "list")
	var out bytes.Buffer
	cmd.Dir = path
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	cmd.Run()
	return out.String() == ""
}

// BranchMergeState tracks for each branch, to which remote branch it has been merged, if any.
type BranchMergeState struct {
	MergedOrigin map[string]string
	MergedPR     map[string]*github.PullRequest
	Unmerged     map[string]struct{}
}

func NewBranchMergeState() *BranchMergeState {
	return &BranchMergeState{
		MergedOrigin: make(map[string]string),
		MergedPR:     make(map[string]*github.PullRequest),
		Unmerged:     make(map[string]struct{}),
	}
}

func (b *BranchMergeState) AddMerged(local, remote string, pr *github.PullRequest) {
	b.MergedOrigin[local] = remote
	b.MergedPR[local] = pr
}

func (b *BranchMergeState) AddUnmerged(local string) {
	b.Unmerged[local] = struct{}{}
}

func (b *BranchMergeState) Clean() bool {
	return len(b.Unmerged) == 0
}

func (b *BranchMergeState) Dump() {
	fmt.Println("Merged:")
	// get longest key
	longest := 0
	for k := range b.MergedOrigin {
		if len(k) > longest {
			longest = len(k)
		}
	}
	// define format string using longest length so they all align nicely
	format := fmt.Sprintf("%%-%ds -> %%s\n", longest)

	for k, v := range b.MergedOrigin {
		fmt.Printf(format, k, v)
		fmt.Printf("    PR #%5d : %s\n", b.MergedPR[k].GetNumber(), b.MergedPR[k].GetTitle())
		fmt.Printf("    Merged %t by %s\n", b.MergedPR[k].GetMerged(), b.MergedPR[k].GetMergedBy().GetLogin())
		fmt.Printf("    Head: %s\n", b.MergedPR[k].GetHead().GetLabel())
		fmt.Printf("    Base: %s\n", b.MergedPR[k].GetBase().GetLabel())
	}
	fmt.Println("Unmerged:")
	for k := range b.Unmerged {
		fmt.Printf("  %s\n", k)
	}
}

// does the git repository have any commits that are not pushed to the remote?
func isLocalBranchesClean(ctx context.Context, r *git.Repository, client *github.Client) (bool, *BranchMergeState, error) {
	// first get the local branches. they have a name like refs/heads/<branch name> or use b.Name().Short()
	branches, err := r.Branches()
	if err != nil {
		panic(err)
	}

	bms := NewBranchMergeState()

	// then get the commits for each branch and check if the commit is in the remote
	// it could be in a branch with the same name, or in a branch with a different name
	err = branches.ForEach(func(b *plumbing.Reference) error {
		fmt.Println("branch:", b.Name().Short(), "sha", b.Hash())

		var remotesFound int

		for _, remName := range []string{"origin", "upstream"} {
			rem, err := r.Remote(remName)
			if err == git.ErrRemoteNotFound {
				continue
			}
			if err != nil {
				return err
			}
			remotesFound++
			pr, err := isCommitInRemote(ctx, client, rem, b.Hash())
			if err != nil {
				return err
			}
			if pr != nil {
				bms.AddMerged(b.Name().Short(), remName, pr)
				return nil
			}
		}
		if remotesFound == 0 {
			return fmt.Errorf("no suitable upstream/origin remote found")
		}
		bms.AddUnmerged(b.Name().Short())
		return nil
	})

	if err != nil {
		return false, bms, fmt.Errorf("error while iterating over branches: %w", err)
	}

	return bms.Clean(), bms, nil
}

func isCommitInRemote(ctx context.Context, client *github.Client, rem *git.Remote, commit plumbing.Hash) (*github.PullRequest, error) {
	remoteUrl := rem.Config().URLs[0]

	owner, name, err := extractOwnerAndNameFromRemoteUrl(remoteUrl)
	if err != nil {
		return nil, err
	}
	opts := github.PullRequestListOptions{
		State: "closed",
	}
	prs, _, err := client.PullRequests.ListPullRequestsWithCommit(ctx, owner, name, commit.String(), &opts)
	if err != nil {
		if strings.Contains(err.Error(), "No commit found for SHA") {
			fmt.Println("commit", commit.String(), "is NOT in remote", rem.Config().Name)
			return nil, nil
		}
		return nil, err
	}

	if len(prs) > 0 {
		fmt.Println("commit", commit.String(), "is in remote", rem.Config().Name)
		return prs[0], nil
	}
	fmt.Println("commit", commit.String(), "is NOT in remote", rem.Config().Name)
	return nil, nil
}

// extractOwnerAndNameFromRemoteUrl extracts the owner and name from a remote URL
// the remoteURL is in the form of either git@github.com:<owner>/<name>.git or https://github.com/<owner>/<name>.git
func extractOwnerAndNameFromRemoteUrl(remoteUrl string) (string, string, error) {
	str := strings.TrimSuffix(remoteUrl, ".git")
	str = strings.TrimPrefix(str, "git@github.com:")
	str = strings.TrimPrefix(str, "https://github.com/")
	split := strings.Split(str, "/")
	if len(split) != 2 {
		return "", "", fmt.Errorf("invalid remote url: %s", remoteUrl)
	}
	return split[0], split[1], nil
}
