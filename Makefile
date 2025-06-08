


TEST_DIR ?= ./example/...
test:
	@echo "▶️  Running tests with coverage for: $(TEST_DIR)"
	@go test -coverprofile=coverage.out $(TEST_DIR)

clean-coverage:
	@rm -f coverage.out coverage.html

diff:
	@echo "▶️  Running diff for coverage"
	@git diff 67e52b27 4518e2cc > pr.diff