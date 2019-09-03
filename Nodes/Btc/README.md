## Quick start

Requires that [Docker be installed](https://docs.docker.com/engine/installation/) on the host machine.

 Build container:
    
    docker build --rm -t bitcoind:latest .
```
# Create some directory where your bitcoin data will be stored.
$ mkdir /home/youruser/bitcoin_data

$  docker run --name bitcoind -d \
   --env 'BTC_RPCUSER=btcAdapter' \
   --env 'BTC_RPCPASSWORD=AdapterBTC' \
   --env 'BTC_TXINDEX=1' \
   --env 'BTC_DISABLEWALLET=0' \
   --volume /opt/data/bitcoin_data:/root/.bitcoin \
   -p 8332:8332 \
   --publish 8333:8333 \
   bitcoind:latest

$ docker logs -f bitcoind
[ ... ]
```


## Configuration

A custom `bitcoin.conf` file can be placed in the mounted data directory.
Otherwise, a default `bitcoin.conf` file will be automatically generated based
on environment variables passed to the container:

| name | default |
| ---- | ------- |
| BTC_RPCUSER | btc |
| BTC_RPCPASSWORD | changemeplz |
| BTC_RPCPORT | 8332 |
| BTC_RPCALLOWIP | 0.0.0.0/0|::/0 |
| BTC_RPCCLIENTTIMEOUT | 30 |
| BTC_DISABLEWALLET | 1 |
| BTC_TXINDEX | 0 |
| BTC_TESTNET | 0 |
| BTC_DBCACHE | 0 |
| BTC_ZMQPUBHASHTX | tcp://0.0.0.0:28333 |
| BTC_ZMQPUBHASHBLOCK | tcp://0.0.0.0:28333 |
| BTC_ZMQPUBRAWTX | tcp://0.0.0.0:28333 |
| BTC_ZMQPUBRAWBLOCK | tcp://0.0.0.0:28333 |


## Daemonizing

The smart thing to do if you're daemonizing is to use Docker's [builtin
restart
policies](https://docs.docker.com/config/containers/start-containers-automatically/#use-a-restart-policy),
but if you're insistent on using systemd, you could do something like

```
$ cat /etc/systemd/system/bitcoind.service

# bitcoind.service #######################################################################
[Unit]
Description=Bitcoind
After=docker.service
Requires=docker.service

[Service]
ExecStartPre=-/usr/bin/docker kill bitcoind
ExecStartPre=-/usr/bin/docker rm bitcoind
ExecStartPre=/usr/bin/docker pull jamesob/bitcoind
ExecStart=/usr/bin/docker run \
    --name bitcoind \
    -p 8333:8333 \
    -p 8332:8332 \
    -v /data/bitcoind:/root/.bitcoin \
    jamesob/bitcoind
ExecStop=/usr/bin/docker stop bitcoind
```

to ensure that bitcoind continues to run.