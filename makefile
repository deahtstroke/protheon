.PHONY test:
	bin/protheon run --dsn --compress zstd "postgres://user:password@localhost:5432/mydb?sslmode=disable" --format jsonl --input ~/15780000000-15790000000.jsonl.zst --script transform.lua --table user

