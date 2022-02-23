test:
	go test -v ./...

test-local:
	sudo go test -v -local ./local_test.go