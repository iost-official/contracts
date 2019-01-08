#!/usr/bin/env bash

iwallet account --import admin 2yquS3ySrGWPEKywCPzX4RTJugqRh7kJSo5aehsLYPEWkUxBWA39oMrZ7ZxuM4fgyXYs2cPwh5n8aNNpH5x2VyK1
iwallet \
 --expiration 10000 --gas_limit 1000000 --gas_ratio 1 \
 --server localhost:30002 \
 --account admin \
 --amount_limit '*:unlimited' \
 call "ContractEm6NRHWGJ9f2XWBayKMQfffXpGL7vpfF49HhqLs8tD63" 'read' '[]'
 #call "token.iost" "transfer" '["iost","someone","me","10000.00","trasfer from someone not exist"]'

