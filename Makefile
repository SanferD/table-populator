.PHONY: build run test coverage clean echo_test_pkgs

IGNORE_PATTERN := "^github\.com/SanferD/table-populator$$\|^github\.com/SanferD/table-populator/ioutil$$"
TEST_PKGS := $(shell go list ./... | grep -vE "$(IGNORE_PATTERN)" | awk -F'/' '{print $$NF}')

build:
	go build -o build/table-populator

run:
	./build/table-populator

test:
	@for pkg in $(TEST_PKGS); do \
		echo "Testing local package \"$$pkg\""; \
		go test ./$$pkg > /dev/null; \
	done

coverage:
	@for pkg in $(TEST_PKGS); do \
		echo "Testing and coverage for local package \"$$pkg\""; \
		go test -coverprofile="build/$$pkg.cover" ./$$pkg > /dev/null 2>&1 && \
		grep -v 'mock_.*\.go' "build/$$pkg.cover" > "build/$$pkg.tmp" && \
		mv "build/$$pkg.tmp" "build/$$pkg.cover"; \
	done
	@gocovmerge build/*.cover > build/.table-populator.cover
	@go tool cover -html=build/.table-populator.cover -o build/table-populator-coverage.html
	@rm -f build/*.cover
	@xdg-open build/table-populator-coverage.html || true

clean:
	rm -f build/*
