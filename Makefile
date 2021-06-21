run:
	go run cmd/svr/main.go
run-with-docker:
	docker-compose -f ./docker/docker-compose.yaml up -d --build --force-recreate
test:
	go test ./...
coverage:
	go test -failfast=true ./... -coverprofile cover.out
	go tool cover -html=cover.out
	rm cover.out
mocks:
	mockery --name=NoteDaoHandler --recursive=true --case=underscore --output=./pkg/testhelper/mocks;
	mockery --name=ExtAPIHandler --recursive=true --case=underscore --output=./pkg/testhelper/mocks;
	mockery --name=NoteServiceHandler --recursive=true --case=underscore --output=./pkg/testhelper/mocks;
	mockery --name=Requester --recursive=true --case=underscore --output=./pkg/testhelper/mocks;
