FROM ubuntu:xenial

WORKDIR /app

ARG PROTOC_VER=3.9.1
ARG GO_VER=1.12.10

ARG BUILD_PACKAGES="automake autoconf git make g++ \
             libtool-bin python3-pip curl unzip "

# Add repository to be able to install latest git
RUN    apt-get update && apt-get install -y --no-install-recommends wget software-properties-common && apt-add-repository ppa:git-core/ppa \
    && apt-get update && apt-get install -y --no-install-recommends $BUILD_PACKAGES \
    && apt-get clean && rm -rfv /var/cache/apt /var/lib/apt/lists/* /tmp/* /var/tmp/* /usr/share/doc

# Download in the home directory (root/go)
RUN wget https://dl.google.com/go/go${GO_VER}.linux-amd64.tar.gz && tar -xvf *.tar.gz -C ~ && rm -rf *.gz
RUN cd / && wget https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VER}/protoc-${PROTOC_VER}-linux-x86_64.zip && unzip *.zip && rm *.zip

ENV GOPATH="/go_workspace"
ENV PATH="/root/go/bin:${GOPATH}/bin:${PATH}"

RUN git clone --depth 1 --branch v1.24.2 https://github.com/grpc/grpc && \
    cd grpc && \
    git submodule update --init && \
    make -j3 install && \
    make clean && \
    make grpc_python_plugin && \
    cd third_party/protobuf && ./autogen.sh && ./configure --disable-shared && make -j3 install && \
    cd /app && \
    cp /app/grpc/bins/opt/grpc_python_plugin /usr/bin/grpc_python_plugin && \
    rm -rf grpc

RUN pip3 install \
    grpcio==1.11.0 \
    grpcio-tools==1.11.0 \
    protobuf==3.6.1


COPY ./scripts/get_go_repo.sh /app/
# get required packages based on version and reinstall if necessary
# some packages do not version so we need to checkout out based on commit
# some packages do not need to be installed as they are just source files (sometimes generated ones i.e genproto)
# the order of these packages matter so that each dependency is only downloaded once
RUN bash get_go_repo.sh github.com/grpc-ecosystem/grpc-gateway                          no  v1.11.3
RUN bash get_go_repo.sh google.golang.org/grpc                                          no  v1.24.0
RUN bash get_go_repo.sh github.com/gogo/googleapis                                      no  v1.3.0
RUN bash get_go_repo.sh github.com/gogo/protobuf                                        no  v1.3.1
RUN bash get_go_repo.sh github.com/gogo/protobuf/proto                                  no  v1.3.1
RUN bash get_go_repo.sh github.com/gogo/protobuf/gogofast                               no  v1.3.1
RUN bash get_go_repo.sh github.com/gogo/protobuf/protoc-gen-gogofast                    yes v1.3.1
RUN bash get_go_repo.sh github.com/gogo/protobuf/gogoproto                              no  v1.3.1
RUN bash get_go_repo.sh github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway  yes v1.12.2
RUN bash get_go_repo.sh github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger       yes v1.12.2

COPY ./scripts/proto_gen.bash   /app/
ENTRYPOINT ["bash", "/app/proto_gen.bash"]
