My Supply Audit
===========

Yet another independent way to check the supply of ETH.

This is a Go program that uses Ethereum's Geth client, both as library & as an RPC node, 
to check the total number of issued Ether on Ethereum.

It works by grabbing each block, then summing up the rewards for each block.

Yup, that easy.

Building:
======


Assuming you have Go installed

```


$ git clone https://github.com/currencytycoon/mysupplyaudit
$ cd mysupplyaudit
$ cp config.json.dist config.json

(Edit your config.json and change your RPC URL if needed, default is usually good)

$ cd cmd/mysupplyaudit
$ go build


```

This will produce a program named `mysupplyaudit` in that directory


Configuration
========
So far the config file is just one option `rpc_url`
This can be a http url, a websocket or a unix socket path.

url example (default): `http://127.0.0.1:8545`

websocket example: `ws://localhost:8546`
(Run Geth with the `--ws` option)

unix socket: `/home/fred/.ethereum/geth.ipc`

By far, using the unix socket is the fastest as it has the least overhead

Running
=================

Make sure your Geth node is all synced to the latest block, running with the `--rpc` option, then

`$ ./mysupplyaudit audit -c ../../config.json`

optional, specify the highest block number to stop at (eg  block number`10000000`)

`$ ./mysupplyaudit audit 10000000 -c ../../config.json`


The program `mysupplyaudit` should be in your cmd/mysupplyaudit (See Building instructions above)

It will take a while...

License
====

Copyright 2020 @0xTycoon

The MIT license