linux:
	env GOOS=linux ARCH=64 go build -v .
windows:
	env GOOS=windows ARCH=64 go build -v .
