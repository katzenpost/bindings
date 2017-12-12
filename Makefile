all: python

python:
	GODEBUG=cgocheck=0 gopy bind -lang="py2" .
