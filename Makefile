test:
	go test -v ./...

test-local:
	sudo go test -v -local ./local_test.go

lint:
	golint -set_exit_status ./...

.PHONY: test test-local lint