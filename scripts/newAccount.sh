#! /bin/bash

if [[ ! -f "~/.wallet/admin_ed25519" ]]; then
  iwallet account --import admin 2yquS3ySrGWPEKywCPzX4RTJugqRh7kJSo5aehsLYPEWkUxBWA39oMrZ7ZxuM4fgyXYs2cPwh5n8aNNpH5x2VyK1
fi

iwallet \
  account --create $1 \   # create a account named player0, NOTICE: duplicated name will failed
  --initial_ram 4096  \  # ram in byte
  --initial_gas_pledge 10000 \
  --initial_balance 10000 \
  --amount_limit '*:unlimited' \
  --account admin \       # set account to publish new account transaction
  --server 47.244.109.92:30002