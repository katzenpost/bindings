all: katzenpost.so

katzenpost.so: python/*.go
	GODEBUG=cgocheck=0 gopy bind -lang="py2" ./python
