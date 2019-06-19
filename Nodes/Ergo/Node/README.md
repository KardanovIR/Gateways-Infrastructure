    
    docker run --rm -d \
           -p 9020:9020 \
           -p 9052:9052 \
           -v ergo:/home/ergo/.ergo \
           -v test.conf:/etc/test.conf \
           ergoplatform/ergo:v2.0.3 --testnet -c /etc/test.conf