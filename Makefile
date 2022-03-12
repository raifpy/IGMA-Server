installocr:
	python -m "import easyocr" || pip install easyocr
	chmod +x assets/easyocrbin
	sudo cp assets/easyocrbin /usr/bin/ #?is sudo ok
	echo "OK!"
build:
	env CGO_ENABLED=0 go build -ldflags="-w -s" -x -v .