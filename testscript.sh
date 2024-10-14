#!/bin/sh

dir=$(mktemp -d)
cd $dir
pwd

function mk () {
	echo "making $1"
	git clone -q git@github.com:caarlos0/fork-cleaner.git $1
}

mk repo-clean

mk repo-dirty-new-file
echo dirty-file > repo-dirty-new-file/dirty-file

# note: we can consider this as clean, if there's no other changes
mk repo-clean-new-dir
mkdir repo-clean-new-dir/dirty-dir

mk    repo-dirty-new-file-in-new-dir
mkdir repo-dirty-new-file-in-new-dir/dirty-dir
echo dirty-file > repo-dirty-new-file-in-new-dir/dirty-dir/dirty-file


mk repo-dirty-stash
cd repo-dirty-stash
echo "some change" >> README.md
git stash
cd ..

mk repo-dirty-index
cd repo-dirty-index
echo "some change" >> README.md
git add README.md
cd ..

mk repo-dirty-commit-to-main
cd repo-dirty-commit-to-main
echo "some change" >> README.md
git commit README.md -m 'extra line'
cd ..

mk repo-dirty-removed-file-in-main
cd repo-dirty-removed-file-in-main
rm README.md
cd ..

mk repo-dirty-commit-to-existing-branch
cd repo-dirty-commit-to-existing-branch
git checkout list # some pre-existing branch from github. this may stop working in the future
echo "some change" >> README.md
git commit README.md -m 'extra line'
cd ..

# same, but a bit more hidden as we check main out again
mk repo-dirty-commit-to-existing-branch-back-to-main
cd repo-dirty-commit-to-existing-branch-back-to-main
git checkout list # some pre-existing branch from github. this may stop working in the future
echo "some change" >> README.md
git commit README.md -m 'extra line'
git checkout main
cd ..



mk repo-clean-other-branch
cd repo-clean-other-branch
git checkout -b some-other-branch
cd ..

mk repo-dirty-commit-to-new-branch
cd repo-dirty-commit-to-new-branch
git checkout -b some-other-branch
echo "some change" >> README.md
git commit README.md -m 'extra line'
cd ..

mk repo-dirty-commit-to-new-branch-back-to-main
cd repo-dirty-commit-to-new-branch-back-to-main
git checkout -b some-other-branch
echo "some change" >> README.md
git commit README.md -m 'extra line'
git checkout main
cd ..

# special cases we don't need to do anything special for:
# have a local branch that:
# does/does not exist in any remote (we don't have to differentiate between upstream and our git fork, because fork-cleaner in classic 'github mode' will figure out whether our fork contains anything special or not)
# has been merged as a squash commit (these cases are probably to hard to figure out with a script and we may need to leave these marked as 'dirty' and manually look into them)
