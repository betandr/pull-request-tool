# Pull Request Tool

[![Build Status](https://drone.andr.io/api/badges/betandr/pull-request-tool/status.svg)](https://drone.andr.io/betandr/pull-request-tool)

_`prt` is a tool that allows you to list pull requests on GitHub
from the command line._

*NB: PRT is not complete and is not ready to actually use yet! :)*

## Usage

You'll need a GitHub access token from https://github.com/settings/tokens to use
`prt`.

`export OAUTH_TOKEN={yourtokenhere}`

List:
`prt list owner/repo (optional: --all)`

Get:
`prt get owner/repo {number}`

Create:
`prt create owner/repo {branch} {base} {title}`

Merge:
`prt merge owner/repo {number} (optional: {title} {message} {method} (merge, squash or rebase))`

## TODO

- Add comment
- List Review comments: https://api.github.com/repos/owner/repo/pulls/number/comments
- Reply to comment
- Approve
- Request changes
- Review
- Reply to Review
- Resolve review

_betandr/prt is licensed under the_
_*GNU Affero General Public License v3.0*_
