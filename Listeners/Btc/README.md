# GatewaysInfrastructure

## Btc listener

Run :

    go run main.go [parameters...] 
    
    Parameters:
       -config-path     path to config, by default "./config/config.yml"
       -config-path=./Listeners/Btc/config/config.yml for debug in JetBrains Goland
       
Docker :
    
    Build container:
    
    docker build --rm -t gateways-btc-listener:latest .
    
    Run container
    docker run --rm -d -p 5001:5001 --env-file=config/dev.env --name gateways-btc-listener gateways-btc-listener:latest
    
    Read logs from container
    docker logs gateways-btc-listener