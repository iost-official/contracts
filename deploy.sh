~/bin/iwallet account --import admin 2yquS3ySrGWPEKywCPzX4RTJugqRh7kJSo5aehsLYPEWkUxBWA39oMrZ7ZxuM4fgyXYs2cPwh5n8aNNpH5x2VyK1
~/bin/iwallet \
 --expiration 10000 --gaslimit 100000000 --gasprice 100 \
 --server 47.244.109.92:30002 \
 --account admin \
 compile demos/gobang.js demos/gobang.abi
