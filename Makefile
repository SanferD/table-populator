.PHONY: build run test coverage

build:
	go build -o build/table-populator

run:
	./build/table-populator

test:
	go test ./dataio ./locator ./logger ./application ./config

coverage:
	@for pkg in dataio locator logger application config; do \
		echo "Testing and coverage for $$pkg..."; \
		go test -coverprofile="build/.$$pkg.cover" "./$$pkg" && \
		grep -v 'mock_.*\.go' "build/.$$pkg.cover" > "build/.$$pkg.tmp" && \
		mv "build/.$$pkg.tmp" "build/.$$pkg.cover"; \
	done
	@gocovmerge build/.*.cover > build/.table-populator.cover
	@go tool cover -html=build/.table-populator.cover -o build/table-populator-coverage.html
	@rm -f build/.*.cover
	xdg-open build/table-populator-coverage.html

clean:
	rm -f build/*
