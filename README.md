# git-replicator

A CLI tool for checking out multiple copies of the same repository in order to have multiple coding agents work simultaneously.

ref. https://www.anthropic.com/engineering/claude-code-best-practices

## Installation

```
go install github.com/terakoya76/git-replicator@latest
```

## Usage

```sh
$ git-replicator get https://github.com/terakoya76/git-replicator-test
Enumerating objects: 5, done.
Counting objects: 100% (5/5), done.
Compressing objects: 100% (4/4), done.
Total 5 (delta 0), reused 0 (delta 0), pack-reused 0 (from 0)

# path to be cloned
$ cd ~/git-replicator/github.com/terakoya76/git-replicator-test

# initial cloned status
$ ls
base
$ pushd base
$ git branch
* main
$ popd 

# checkout another independent directory
$ git-replicator switch foo
Enumerating objects: 5, done.
Counting objects: 100% (5/5), done.
Compressing objects: 100% (4/4), done.
Total 5 (delta 0), reused 0 (delta 0), pack-reused 0 (from 0)
cloned branch: foo to dir: /home/$USER/git-replicator/github.com/terakoya76/git-replicator-test/foo

$ pushd foo/
$ git branch
* foo
  main
$ popd

# From now on, you can let the agent use this directory as it pleases.
# e.g. 
# $ cursor ./foo/
```

## Features
- Clone a git repository into a structured local directory (`get <url>`)
- List all managed repositories (`list`)
- Clone current repo into a new branch directory (`switch <branch>`, like `git switch`)

## Development

- Build: `make build`
- Lint: `make lint`
- Test: `make test`

## License

MIT