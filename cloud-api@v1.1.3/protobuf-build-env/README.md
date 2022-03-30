# Protobuf build env

This docker file is used to create an environment for generating protobuf files. Right now, Live Planet generates
protobuf files for C++, Go, and Python. This container will be used to unify different build environments so that
eventually each project that requires proto files will use the same version of protobuf and grpc.

# How to build

- `make build && make push`

# Limitations

- all of the output directories need to be somewhere inside the current working directory when invoking docker
- the generated files are currently created with root permissions and require root to delete

# Milestones

- [x] First version of Docker image
- [x] Integration of Docker image file into build process of various repositories/projects
- [ ] Update C++ to 3.7.1 (or whatever the latest protoc version used by Go and Python)
- [ ] Remove generated proto files and update build processes to always pull the docker image and generate the files
