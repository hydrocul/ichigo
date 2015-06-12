
test: go/src/hydrocul/ichigo_test/da.go go/src/hydrocul/ichigo_test/da_test.go go/src/hydrocul/ichigo_test/dict.go go/src/hydrocul/ichigo_test/dict_test.go go/src/hydrocul/ichigo_test/texts.go go/src/hydrocul/ichigo_test/texts_test.go
	GOPATH=$(realpath .)/go go test hydrocul/ichigo_test

go/src/hydrocul/ichigo_test:
	mkdir -p go/src/hydrocul/ichigo_test

go/src/hydrocul/ichigo_test/da.go: da.go go/src/hydrocul/ichigo_test
	cp da.go go/src/hydrocul/ichigo_test/da.go

go/src/hydrocul/ichigo_test/da_test.go: da_test.go go/src/hydrocul/ichigo_test
	cp da_test.go go/src/hydrocul/ichigo_test/da_test.go

go/src/hydrocul/ichigo_test/dict.go: dict.go go/src/hydrocul/ichigo_test
	cp dict.go go/src/hydrocul/ichigo_test/dict.go

go/src/hydrocul/ichigo_test/dict_test.go: dict_test.go go/src/hydrocul/ichigo_test
	cp dict_test.go go/src/hydrocul/ichigo_test/dict_test.go

go/src/hydrocul/ichigo_test/texts.go: texts.go go/src/hydrocul/ichigo_test
	cp texts.go go/src/hydrocul/ichigo_test/texts.go

go/src/hydrocul/ichigo_test/texts_test.go: texts_test.go go/src/hydrocul/ichigo_test
	cp texts_test.go go/src/hydrocul/ichigo_test/texts_test.go



go/bin/ichigo_build: go/src/hydrocul/ichigo_build/main.go go/src/hydrocul/ichigo_build/da.go go/src/hydrocul/ichigo_build/dict.go go/src/hydrocul/ichigo_build/texts.go
	GOPATH=$(realpath .)/go go install hydrocul/ichigo_build

go/src/hydrocul/ichigo_build:
	mkdir -p go/src/hydrocul/ichigo_build

go/src/hydrocul/ichigo_build/main.go: build_main.go go/src/hydrocul/ichigo_build
	cp build_main.go go/src/hydrocul/ichigo_build/main.go

go/src/hydrocul/ichigo_build/da.go: da.go go/src/hydrocul/ichigo_build
	cp da.go go/src/hydrocul/ichigo_build/da.go

go/src/hydrocul/ichigo_build/dict.go: dict.go go/src/hydrocul/ichigo_build
	cp dict.go go/src/hydrocul/ichigo_build/dict.go

go/src/hydrocul/ichigo_build/texts.go: texts.go go/src/hydrocul/ichigo_build
	cp texts.go go/src/hydrocul/ichigo_build/texts.go



go/bin/ichigo: go/src/hydrocul/ichigo/main.go go/src/hydrocul/ichigo/da.go go/src/hydrocul/ichigo/dict.go go/src/hydrocul/ichigo/dict_data.go go/src/hydrocul/ichigo/texts.go
	GOPATH=$(realpath .)/go go install hydrocul/ichigo

go/src/hydrocul/ichigo:
	mkdir -p go/src/hydrocul/ichigo

go/src/hydrocul/ichigo/main.go: main.go go/src/hydrocul/ichigo
	cp main.go go/src/hydrocul/ichigo/main.go

go/src/hydrocul/ichigo/da.go: da.go go/src/hydrocul/ichigo
	cp da.go go/src/hydrocul/ichigo/da.go

go/src/hydrocul/ichigo/dict.go: dict.go go/src/hydrocul/ichigo
	cp dict.go go/src/hydrocul/ichigo/dict.go

go/src/hydrocul/ichigo/dict_data.go: var/dict_data.go go/src/hydrocul/ichigo
	cp var/dict_data.go go/src/hydrocul/ichigo/dict_data.go

go/src/hydrocul/ichigo/texts.go: texts.go go/src/hydrocul/ichigo
	cp texts.go go/src/hydrocul/ichigo/texts.go



var:
	mkdir -p var

var/ipadic:
	sh lib/download-ipadic.sh

var/matrix1.txt: var/ipadic lib/build-matrix1.sh
	sh lib/build-matrix1.sh > var/matrix1.txt

var/dict1.txt: var/ipadic lib/build-dict1.sh
	sh lib/build-dict1.sh > var/dict1.txt

var/dict2.txt: var/dict1.txt lib/build-dict2.sh
	sh lib/build-dict2.sh > var/dict2.txt

var/dict2-texts.txt: var/dict2.txt lib/build-dict2-texts.sh
	sh lib/build-dict2-texts.sh > var/dict2-texts.txt

var/dict3.dat: var/matrix1.txt var/dict2.txt var/dict2-texts.txt go/bin/ichigo_build
	go/bin/ichigo_build var/matrix1.txt var/dict2.txt var/dict2-texts.txt > var/dict3.dat

var/dict_data.go: dict_data.go var/dict3.dat lib/build-dict4.sh
	sh lib/build-dict4.sh > var/dict_data.go

