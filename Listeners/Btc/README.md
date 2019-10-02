# GatewaysInfrastructure

## Btc listener

Development :

    for local development with modules just add
    
    replace github.com/wavesplatform/GatewaysInfrastructure => ../../../GatewaysInfrastructure
    
    to go.mod file

Run :

    go run main.go [parameters...] 
    
    Parameters:
       -config-path     path to config, by default "./config/config.yml"
       -config-path=./Listeners/Btc/config/config.yml for debug in JetBrains Goland
       
Docker :
    
    Build container:
    
    docker build --rm -t gateways-btc-listener:latest .
    
    Run container 
    docker run --rm -d --env-file=config/dev.env -p 5001:5001 --name gateways-btc-listener gateways-btc-listener:latest
    
    Read logs from container
    docker logs gateways-btc-listener