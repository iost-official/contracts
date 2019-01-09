# playground

This is a playground of demo contract of IOST

you can choose 2 ways to test & run demos

## Use simulator in this repo

to be continue

## Send to block chain

We prepared scripts to work well with test block chain. See scripts' name and description to use.

Example below:
``` bash
$ deploy.sh gobang          # deploy gobang using this
$ newAccount.sh player1     # new account and pledge; account name should between 5 and 11 chars in [0-9a-z_]
$ newAccount.sh player2    # another account

echo ~/.iwallet/player1_ed25519  # read the secure key of player1
```

We strongly recommend you using your own account to run tests.