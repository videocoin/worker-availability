# Docker build env

You can use docker to create an environment for generating protobuf files. Right now, Live Planet generates
protobuf files for C++, Go, and Python. This container will be used to unify different build environments so that
eventually each project that requires proto files will use the same version of protobuf and grpc.

# How to use

Pull docker image

- `make docker-pull-image`

Generate files

- `make docker-protoc`