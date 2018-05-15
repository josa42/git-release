# git-release

**Work in progress**

## Installation

**Homebrew (macOS)**

```
brew tap josa42/homebrew-git-tools
brew install git-release
```

**Other**

```
go get github.com/josa42/git-release
```

## Usage

```
Usage:
  git-release [--major|--minor|--patch] [--stable|--beta|--rc] [--dirty] [--force] [--do-not-push] [--no-empty-commit]
  git-release --stable|--beta|--rc                             [--dirty] [--force] [--do-not-push] [--no-empty-commit]
  git-release <version>                                        [--dirty] [--force] [--do-not-push] [--no-empty-commit]
  git-release --help
  git-release --version

Options:
  -h --help          Show this screen.
  --version          Show version.
  --dirty            Include changed files in release commit.
  --force            Force new commit even thought the latest commit is already tagged.
  --do-not-push      Do not push commit and tags
  --no-empty-commit  Do not commit if nothing changed
```
## Config

- **Custom commit message"
  `git config git-release.message "🎉  Release {version}"`

## Update version files

Versions in these files are updated automatically

- `package.json`
- `bower.json`

For more complex version updates, you can provide a script. Which will be called with the new version as argument:

```
.git-release/set-version.sh
```


## License

MIT (See [license.md](license.md))
