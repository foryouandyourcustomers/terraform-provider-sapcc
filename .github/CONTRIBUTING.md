Contributing
============
If you want to contribute to a project and make it better, your help is very welcome. Contributing is also a great way to learn more about social coding on Github, new technologies and their ecosystems and how to make constructive, helpful bug reports, feature requests and the noblest of all contributions: a good, clean pull request.

### Commit messages
In general, always write your commit messages in the present tense. Your commit message should describe what the commit, when applied, does to the code – not what you did to the code. This repository uses [Semantic Releases](https://semantic-release.gitbook.io/semantic-release/) for automatic releases and change logs. This helps keep to the release cycle tidy and useful in long-term maintenance. This is possible by keeping the commit messages clean and structured. 


### Structure of commit messages 
The commit messages must be structured to as conventional-commit:

In a nut-shell, here's the format:

```shell
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

Read the full specification [here](https://www.conventionalcommits.org/en/v1.0.0/#specification)

Make sure the type `feat` actually reflects a new feature being introduced and `fix` actually fixes a bug in the code. For eveything else, one can use `chore` for regular scoped maintenance. Please look at the [CHANGELOG](../CHANGELOG.md) for the existing `scope`s before introducing one. Please link an issue in the commit message body if there's already a linked issue, for eg, `Refs: #<issuse_id>` or `Closes: #<issue_id>` etc.

### How to make a clean pull request

Look for a project's contribution instructions. If there are any, follow them.

- Create a personal fork of the project on Github.
- Clone the fork on your local machine. Your remote repo on Github is called `origin`.
- Add the original repository as a remote called `upstream`.
- If you created your fork a while ago be sure to pull upstream changes into your local repository.
- Create a new branch to work on! Branch from `master`
- Implement/fix your feature, comment your code.
- Follow the code style of the project, including indentation.
- If the project has tests run them!
- Write or adapt tests as needed.
- Add or change the documentation as needed.
- Squash your commits into a single commit with git's [interactive rebase](https://help.github.com/articles/interactive-rebase). Create a new branch if necessary.
- Push your branch to your fork on Github, the remote `origin`.
- From your fork open a pull request in the correct branch. Target the project's `develop` branch if there is one, else go for `master`!
- If the maintainer requests further changes just push them to your branch. The PR will be updated automatically.
- Once the pull request is approved and merged you can pull the changes from `upstream` to your local repo and delete
  your extra branch(es).

_**Note: If the changes introduce breaking changes, please create a PR on `beta` branch not on the master!**_ 
