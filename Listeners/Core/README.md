# GatewaysInfrastructure 

##Listener Core

Recompile protobuf :
    
    protoc -I grpc/ grpc/listener.proto --go_out=plugins=grpc:grpc
    protoc -I grpc/client grpc/client/callback_service.proto --go_out=plugins=grpc:grpc/client
