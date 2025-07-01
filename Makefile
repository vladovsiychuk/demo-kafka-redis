tidy-forwarder:
	(cd forwarder && go mod tidy)

tidy-diff-calculator:
	(cd diff-calculator && go mod tidy)

tidy:
	make tidy-forwarder
	make tidy-diff-calculator

test-forwarder:
	(cd forwarder && go test ./...)

test-diff-calculator:
	(cd diff-calculator && go test ./...)

test:
	make test-forwarder
	make test-diff-calculator

run-forwarder:
	(cd forwarder && go run cmd/forwarder/main.go)

run-diff-calculator:
	(cd diff-calculator && go run cmd/diff-calculator/main.go)

run-apps:
	make run-forwarder & \
	make run-diff-calculator