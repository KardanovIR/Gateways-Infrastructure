# GatewaysInfrastructure 

##Listener Core

Recompile protobuf :
    
    protoc -I grpc/ grpc/listener.proto --go_out=plugins=grpc:grpc
    protoc -I grpc/client grpc/client/callback_service.proto --go_out=plugins=grpc:grpc/client


To make changes available for another module:

Push it to master branch

Add submodule tag:  git tag Listeners/Core/v1.3.2

git push origin Listeners/Core/v1.3.2