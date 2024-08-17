.PHONY: run cpu mem server

run:
	go run main.go

cpu:
	go tool pprof cpu.prof

mem:
	go tool pprof mem.prof


server:
	go tool pprof -http 127.0.0.1:8080 cpu.prof