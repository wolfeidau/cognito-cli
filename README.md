# cognito-cli

This is the missing CLI tool for working with [AWS Cognito](https://aws.amazon.com/cognito), it provides a bunch of utility functions which are designed to make administering Cognito easier.

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

# TODO

The following sub commands enable the operation for all users, or a sub set Using a simple filter.

* [ ] Export all users and create a local json file containing the user data
* [x] Invoke global logout for those users
* [ ] Trigger a password reset for all matching users

# License

This application is released under Apache 2.0 license and is copyright Mark Wolfe.