ichigo: go/bin/ichigo
	cp go/bin/ichigo ichigo

test: test1 test2



test1: go/src/hydrocul/ichigo_test1/da.go go/src/hydrocul/ichigo_test1/da_test.go go/src/hydrocul/ichigo_test1/dict.go go/src/hydrocul/ichigo_test1/dict_test.go go/src/hydrocul/ichigo_test1/texts.go go/src/hydrocul/ichigo_test1/texts_test.go
	GOPATH=$(realpath .)/go go test hydrocul/ichigo_test1

test2: go/src/hydrocul/ichigo_test2/da.go go/src/hydrocul/ichigo_test2/dict.go go/src/hydrocul/ichigo_test2/dict_data.go go/src/hydrocul/ichigo_test2/texts.go go/src/hydrocul/ichigo_test2/pipe.go go/src/hydrocul/ichigo_test2/pipe_test.go go/src/hydrocul/ichigo_test2/posid.go
	GOPATH=$(realpath .)/go go test hydrocul/ichigo_test2 -bench .



var/mkdir:
	mkdir -p var
	touch var/mkdir

var/ipadic: var/mkdir
	sh lib/download-ipadic.sh

var/matrix1.txt: var/mkdir var/ipadic lib/build-matrix1.sh
	sh lib/build-matrix1.sh > var/matrix1.txt

var/dict1.txt: var/mkdir var/ipadic lib/build-dict1.sh
	sh lib/build-dict1.sh > var/dict1.txt

var/dict2.txt: var/mkdir var/dict1.txt lib/build-dict2.sh
	sh lib/build-dict2.sh > var/dict2.txt

var/dict2-texts.txt: var/mkdir var/dict2.txt lib/build-dict2-texts.sh
	sh lib/build-dict2-texts.sh > var/dict2-texts.txt

var/dict3.dat: var/mkdir var/matrix1.txt var/dict2.txt var/dict2-texts.txt go/bin/ichigo_build
	go/bin/ichigo_build var/matrix1.txt var/dict2.txt var/dict2-texts.txt > var/dict3.dat

var/dict_data.go: var/mkdir dict_data.go var/dict3.dat lib/build-dict4.sh
	sh lib/build-dict4.sh > var/dict_data.go



go/bin/ichigo: go/src/hydrocul/ichigo/main.go go/src/hydrocul/ichigo/da.go go/src/hydrocul/ichigo/dict.go go/src/hydrocul/ichigo/dict_data.go go/src/hydrocul/ichigo/texts.go go/src/hydrocul/ichigo/pipe.go go/src/hydrocul/ichigo/posid.go
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

go/src/hydrocul/ichigo/pipe.go: pipe.go go/src/hydrocul/ichigo
	cp pipe.go go/src/hydrocul/ichigo/pipe.go

go/src/hydrocul/ichigo/posid.go: posid.go go/src/hydrocul/ichigo
	cp posid.go go/src/hydrocul/ichigo/posid.go



go/src/hydrocul/ichigo_test1:
	mkdir -p go/src/hydrocul/ichigo_test1

go/src/hydrocul/ichigo_test1/da.go: da.go go/src/hydrocul/ichigo_test1
	cp da.go go/src/hydrocul/ichigo_test1/da.go

go/src/hydrocul/ichigo_test1/da_test.go: da_test.go go/src/hydrocul/ichigo_test1
	cp da_test.go go/src/hydrocul/ichigo_test1/da_test.go

go/src/hydrocul/ichigo_test1/dict.go: dict.go go/src/hydrocul/ichigo_test1
	cp dict.go go/src/hydrocul/ichigo_test1/dict.go

go/src/hydrocul/ichigo_test1/dict_test.go: dict_test.go go/src/hydrocul/ichigo_test1
	cp dict_test.go go/src/hydrocul/ichigo_test1/dict_test.go

go/src/hydrocul/ichigo_test1/texts.go: texts.go go/src/hydrocul/ichigo_test1
	cp texts.go go/src/hydrocul/ichigo_test1/texts.go

go/src/hydrocul/ichigo_test1/texts_test.go: texts_test.go go/src/hydrocul/ichigo_test1
	cp texts_test.go go/src/hydrocul/ichigo_test1/texts_test.go



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



go/src/hydrocul/ichigo_test2:
	mkdir -p go/src/hydrocul/ichigo_test2

go/src/hydrocul/ichigo_test2/da.go: da.go go/src/hydrocul/ichigo_test2
	cp da.go go/src/hydrocul/ichigo_test2/da.go

go/src/hydrocul/ichigo_test2/dict.go: dict.go go/src/hydrocul/ichigo_test2
	cp dict.go go/src/hydrocul/ichigo_test2/dict.go

go/src/hydrocul/ichigo_test2/dict_data.go: var/dict_data.go go/src/hydrocul/ichigo_test2
	cp var/dict_data.go go/src/hydrocul/ichigo_test2/dict_data.go

go/src/hydrocul/ichigo_test2/texts.go: texts.go go/src/hydrocul/ichigo_test2
	cp texts.go go/src/hydrocul/ichigo_test2/texts.go

go/src/hydrocul/ichigo_test2/pipe.go: pipe.go go/src/hydrocul/ichigo_test2
	cp pipe.go go/src/hydrocul/ichigo_test2/pipe.go

go/src/hydrocul/ichigo_test2/pipe_test.go: pipe_test.go go/src/hydrocul/ichigo_test2
	cp pipe_test.go go/src/hydrocul/ichigo_test2/pipe_test.go

go/src/hydrocul/ichigo_test2/posid.go: posid.go go/src/hydrocul/ichigo_test2
	cp posid.go go/src/hydrocul/ichigo_test2/posid.go





