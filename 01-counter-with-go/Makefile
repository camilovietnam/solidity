abi:
	solc --optimize --overwrite --abi ./counter.sol -o build
bin: 
	solc --optimize --bin ./counter.sol -o build
go: 
	abigen --abi=./build/counter.abi --bin=./build/counter.bin --pkg=api --out=./api/counter.go
run:
	go run .