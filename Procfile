node1: ./build/bin/./geth --unlock 0x07Bc3e4c91a3428573E33A259b62705d71991C9f --config datadir/node1/config.toml --mine --password datadir/password.txt --verbosity 4 1>datadir/logs/node1.log 2>datadir/logs/node1_err.log
node2: ./build/bin/./geth --unlock 0xfbca26D4dfe3d70A4AE021eD0a2266B492fD4C9d --config datadir/node2/config.toml --password datadir/password.txt --verbosity 4 1>datadir/logs/node2.log 2>datadir/logs/node2_err.log
node3: ./build/bin/./geth --unlock 0x87b42f0debdddeacfa791e78b86417d684282984 --config datadir/node3/config.toml --password datadir/password.txt --verbosity 4 1>datadir/logs/node3.log 2>datadir/logs/node3_err.log
node4: ./build/bin/./geth --config datadir/node4/config.toml  --verbosity 4  1>datadir/logs/node4.log 2>datadir/logs/node4_err.log
bootstrap: ./build/bin/bootnode -nodekey datadir/boot.key -addr :30305 -verbosity 1
