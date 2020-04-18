# cognito-cli

This is the missing CLI tool for working with [AWS Cognito](https://aws.amazon.com/cognito), it provides a bunch of utility functions which are designed to make administering Cognito easier.

[![GitHub Actions status](https://github.com/wolfeidau/cognito-cli/workflows/Go/badge.svg?branch=master)](https://github.com/wolfeidau/cognito-cli/actions?query=workflow%3AGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/wolfeidau/cognito-cli)](https://goreportcard.com/report/github.com/wolfeidau/cognito-cli) [![Coverage Status](https://coveralls.io/repos/github/wolfeidau/cognito-cli/badge.svg?branch=master)](https://coveralls.io/github/wolfeidau/cognito-cli?branch=master)

# Installation

Install it using [gobinaries](https://gobinaries.com/).

```
curl -sf https://gobinaries.com/wolfeidau/cognito-cli | sh
```

Download the latest release from the [release page](https://github.com/wolfeidau/cognito-cli/releases).

Or if you have a working Go installation (Go-1.14 or higher) and want to build `cognito-cli` it fun the following command to install in your `$GOPATH/bin` folder.

```bash
GO111MODULE=off go get -u github.com/wolfeidau/cognito-cli
```

# Usage

I have pulled out the CLI help for each sub command and provided them in the order which they are best explored in.

To display a list of user pools.

```
Usage: cognito-cli ls

List pools.

Flags:
  --help                  Show context-sensitive help.
  --debug                 Enable debug mode.
  --disable-local-time    Disable localisation of times output.
  --version

  --csv                   Enable csv output.

```

To list the schema attributes of a user pool.

```
Usage: cognito-cli list-attributes --user-pool-id=STRING

List the schema attributes of the user pool.

Flags:
  --help                   Show context-sensitive help.
  --debug                  Enable debug mode.
  --disable-local-time     Disable localisation of times output.
  --version

  --user-pool-id=STRING
```

To find users in a pool.

```
Usage: cognito-cli find --user-pool-id=STRING

Find users.

Flags:
  --help                             Show context-sensitive help.
  --debug                            Enable debug mode.
  --disable-local-time               Disable localisation of times output.
  --version

  --user-pool-id=STRING
  --csv                              Enable csv output.
  --attributes=Username,email,...    Attributes to retrieve and output.
  --back-off=500                     Delay in ms used to backoff during paging of records
  --filter=KEY=VALUE;...             Filter users based on a set of patterns, supports '*' and '?' wildcards in either string.

```

To export users in a pool, filter and write the results in CSV format, I use this for analysis and to verify the integrity of user details stored.

```
Usage: cognito-cli export --user-pool-id=STRING

Export users, filter and write the results in CSV format.

Flags:
  --help                    Show context-sensitive help.
  --debug                   Enable debug mode.
  --disable-local-time      Disable localisation of times output.
  --version

  --user-pool-id=STRING
  --back-off=500            Delay in ms used to backoff during paging of records
  --filter=KEY=VALUE;...    Filter users based on a set of patterns, supports '*' and '?' wildcards in either string.

```

**NOTE:** This is not designed to directly feed into the [StartUserImportJob](https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_StartUserImportJob.html) operation, some transformations may be required.

To find users in a pool and log them out.

```
Usage: cognito-cli logout --user-pool-id=STRING

Find users and trigger a logout.

Flags:
  --help                    Show context-sensitive help.
  --debug                   Enable debug mode.
  --disable-local-time      Disable localisation of times output.
  --version

  --user-pool-id=STRING
  --back-off=500            Delay in ms used to backoff during paging of records
  --filter=KEY=VALUE;...    Filter users based on a set of patterns, supports '*' and '?' wildcards in either string.

```

**NOTE:** This calls [AdminUserGlobalSignOut](https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_AdminUserGlobalSignOut.html) for each user which revokes their refresh token, it won't actually log them out straight away, they will however have to log in the next time they attempt to refresh their access token. This can take up to 1 hour for users who have recently refreshed their token.

# License

This application is released under Apache 2.0 license and is copyright Mark Wolfe.
