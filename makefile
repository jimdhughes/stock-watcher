linux_64:
	env GOOS=linux ARCH=64 go build -v .
windows_64:
	env GOOS=windows ARCH=64 go build -v .
macos_64:
	env GOOS=darwin ARCH=64 go build -v
