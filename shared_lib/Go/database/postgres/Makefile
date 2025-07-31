coverage_test:
	go test -coverpkg=./... -coverprofile=coverage.out ./test/...	&& \
			go tool cover -html=coverage.out -o coverage.html
