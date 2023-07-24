#!/bin/sh

dir=$(mktemp -d)
cd $dir
pwd

function mk () {
	echo "making $1"
	git clone -q git@github.com:caarlos0/fork-cleaner.git $1
}

mk repo-clean

mk repo-new-file
echo dirty-file > repo-new-file/dirty-file

# note: we can consider this as clean, if there's no other changes
mk repo-new-dir
mkdir repo-new-dir/dirty-dir

mk    repo-new-file-in-new-dir
mkdir repo-new-file-in-new-dir/dirty-dir
echo dirty-file > repo-new-file-in-new-dir/dirty-dir/dirty-file


mk repo-stash
cd repo-stash
echo "some change" >> README.md
git stash
cd ..

mk repo-dirty-index
cd repo-dirty-index
echo "some change" >> README.md
git add README.md
cd ..

mk repo-commit-to-main
cd repo-commit-to-main
echo "some change" >> README.md
git commit README.md -m 'extra line'
cd ..

# TODO more test cases:
# have a local branch that:
# does/does not exist in any remote (we don't have to differentiate between upstream and our git fork, because fork-cleaner in classic 'github mode' will figure out whether our fork contains anything special or not)
# has a commit ahead of main or the remote branch
# has been merged as a squash commit (these cases are probably to hard to figure out with a script and we may need to leave these marked as 'dirty' and manually look into them)
