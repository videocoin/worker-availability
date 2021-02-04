
docker-mock-accounts:
	${DOCKER_MOCK_COMMAND} -package mock \
	 -source /go_workspace/src/github.com/videocoin/cloud-api/accounts/v1/account_service.pb.go \
	 -destination /go_workspace/src/github.com/videocoin/cloud-api/mock/accounts.go

docker-mock-miners:
	${DOCKER_MOCK_COMMAND} -package mock \
	 -source /go_workspace/src/github.com/videocoin/cloud-api/miners/v1/miner_service.pb.go \
	 -destination /go_workspace/src/github.com/videocoin/cloud-api/mock/miners.go

docker-mock-users:
	${DOCKER_MOCK_COMMAND} -package mock \
	 -source /go_workspace/src/github.com/videocoin/cloud-api/users/v1/user_service.pb.go \
	 -destination /go_workspace/src/github.com/videocoin/cloud-api/mock/users.go

docker-mock-streams:
	${DOCKER_MOCK_COMMAND} -package mock \
	 -source /go_workspace/src/github.com/videocoin/cloud-api/streams/private/v1/streams_service.pb.go \
	 -destination /go_workspace/src/github.com/videocoin/cloud-api/mock/streams.go

docker-mock-%:
	${DOCKER_MOCK_COMMAND} -package mock \
	 -source /go_workspace/src/github.com/videocoin/cloud-api/$*/v1/$*_service.pb.go \
	 -destination /go_workspace/src/github.com/videocoin/cloud-api/mock/$*.go
