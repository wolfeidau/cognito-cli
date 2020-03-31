# cognito-cli

This is the missing CLI tool for working with [AWS Cognito](https://aws.amazon.com/cognito), it provides a bunch of utility functions which are designed to make administering Cognito easier.

# usage

To list pools.

```
Usage: cognito-cli ls

List pools.

Flags:
  --help     Show context-sensitive help.
  --debug    Enable debug mode.

  --csv      Enable csv output.
```

To find users in a pool.

```
Usage: cognito-cli find --user-pool-id=STRING

Find users.

Flags:
  --help                    Show context-sensitive help.
  --debug                   Enable debug mode.

  --user-pool-id=STRING
  --csv                     Enable csv output.
  --attributes=email,...    Attributes to retrieve and output.
  --back-off=500            Delay in ms used to backoff during paging of records

```

# TODO

The following sub commands enable the operation for all users, or a sub set Using a simple filter.

* [ ] Export all users and create a local json file containing the user data
* [ ] Invoke global logout for those users
* [ ] Trigger a password reset for all matching users

# License

This application is released under Apache 2.0 license and is copyright Mark Wolfe.