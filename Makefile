.PHONY: run cpu mem

run:
	go run main.go

cpu:
	go tool pprof cpu.prof

mem:
	go tool pprof mem.prof