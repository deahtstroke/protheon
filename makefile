.PHONY: test build link clean

test:
	bin/protheon run --dsn --compress zstd "postgres://user:password@localhost:5432/mydb?sslmode=disable" --format jsonl --input ~/15780000000-15790000000.jsonl.zst --script transform.lua --table user

build:
	go build -o bin/protheon main.go

link: build
	sudo ln -sf $(PWD)/bin/protheon /usr/local/bin/protheon

clean:
	rm -rf bin/
