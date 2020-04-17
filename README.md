# cognito-cli

This is the missing CLI tool for working with [AWS Cognito](https://aws.amazon.com/cognito), it provides a bunch of utility functions which are designed to make administering Cognito easier.

[![Coverage Status](https://coveralls.io/repos/github/wolfeidau/cognito-cli/badge.svg?branch=master)](https://coveralls.io/github/wolfeidau/cognito-cli?branch=master)

# usage

To list pools.

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

To list the attributes which are configured on a user pool.

```
Usage: cognito-cli list-attributes --user-pool-id=STRING

List the user attributes for a pool.

Flags:
  --help                   Show context-sensitive help.
  --debug                  Enable debug mode.
  --disable-local-time     Disable localisation of times output.
  --version

  --user-pool-id=STRING
```

To find users in a pool and export the results in CSV format.

```
Usage: cognito-cli export --user-pool-id=STRING

Find users and export in CSV format.

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