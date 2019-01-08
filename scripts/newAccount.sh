#!/bin/bash

iwallet account --import admin 2yquS3ySrGWPEKywCPzX4RTJugqRh7kJSo5aehsLYPEWkUxBWA39oMrZ7ZxuM4fgyXYs2cPwh5n8aNNpH5x2VyK1

iwallet \
  account --create $1 \
  --initial_ram 65536  \
  --initial_gas_pledge 1000000 \
  --initial_balance 10000 \
  --amount_limit '*:unlimited' \
  --account admin \
  --server 47.244.109.92:30002