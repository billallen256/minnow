build:
	go build -v -o minnow cmd/minnow/main.go

example_data:
	python3 examples/make_ingest_data.py examples/ingest/

clean_example:
	rm examples/ingest/*
	rm -rf examples/work/*

run: build clean_example example_data
	./minnow examples/config.properties

tests:
	cd internal/minnow; go test
