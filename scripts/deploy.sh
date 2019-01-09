#!/bin/bash

iwallet account --import admin 2yquS3ySrGWPEKywCPzX4RTJugqRh7kJSo5aehsLYPEWkUxBWA39oMrZ7ZxuM4fgyXYs2cPwh5n8aNNpH5x2VyK1
iwallet \
 --expiration 1000000 --gas_limit 1000000 --gas_ratio 1 \
 --server 47.244.109.92:30002 \
 --account admin \
 --amount_limit '*:unlimited' \
 publish ../demos/$1.js ../demos/$1.abi
