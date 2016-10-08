# fork-cleaner

Cleans up old and inactive forks on your github account.

Create a personal access token with `repo` and `delete_repo`
permission, then export it as `GITHUB_TOKEN`, then, simply run then
binary. It will show you all repos that:

- are forks
- are not private
- have no forks
- have no stars
- had no activity in the last 1 month

Then, it will ask you if you want to delete them:

![screenshot](https://cloud.githubusercontent.com/assets/245435/19216454/a0201810-8d92-11e6-8edc-4e1fe156b5c2.png)

Read carefully the list, and, if you agree, type `y` and it will
finish the job for you.
