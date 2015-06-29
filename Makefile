
################################################################################
# main target
################################################################################

ichigo: go/bin/ichigo
	cp go/bin/ichigo ichigo



################################################################################
# test scripts
################################################################################

test: test1 test2

test1: go/src/hydrocul/ichigo-test1/da.go go/src/hydrocul/ichigo-test1/da_test.go go/src/hydrocul/ichigo-test1/dict.go go/src/hydrocul/ichigo-test1/dict_test.go go/src/hydrocul/ichigo-test1/texts.go go/src/hydrocul/ichigo-test1/texts_test.go
	GOPATH=$(realpath .)/go go test hydrocul/ichigo-test1

test2: go/src/hydrocul/ichigo-test2/da.go go/src/hydrocul/ichigo-test2/dict.go go/src/hydrocul/ichigo-test2/dict_data.go go/src/hydrocul/ichigo-test2/texts.go go/src/hydrocul/ichigo-test2/pipe.go go/src/hydrocul/ichigo-test2/pipe_test.go go/src/hydrocul/ichigo-test2/posid.go
	GOPATH=$(realpath .)/go go test hydrocul/ichigo-test2 -bench .



################################################################################
# making var directory
################################################################################

var/mkdir:
	mkdir -p var
	touch var/mkdir



################################################################################
# building dictionary
################################################################################

var/ipadic: var/mkdir
	sh dict/download-ipadic.sh

var/matrix.txt: var/mkdir var/ipadic dict/matrix.sh
	sh dict/matrix.sh > var/matrix.txt

var/ipadic1.txt: var/mkdir var/ipadic dict/ipadic1.sh
	sh dict/ipadic1.sh > var/ipadic1.txt

var/ipadic2.txt: var/mkdir var/ipadic1.txt dict/ipadic1.sh
	cat var/ipadic1.txt | sh dict/ipadic2.sh > var/ipadic2.txt

var/dict.txt: var/ipadic2.txt
	cp var/ipadic2.txt var/dict.txt

var/texts.txt: var/mkdir var/dict.txt dict/texts.sh
	cat var/dict.txt | sh dict/texts.sh > var/texts.txt

var/dict.dat: var/mkdir var/matrix.txt var/dict.txt var/texts.txt go/bin/ichigo-build
	go/bin/ichigo-build var/matrix.txt var/dict.txt var/texts.txt > var/dict.dat.tmp
	mv var/dict.dat.tmp var/dict.dat

var/dict_data.go: var/mkdir dict_data.go var/dict.dat dict/to-go-source.sh
	sh dict/to-go-source.sh var/dict.dat > var/dict_data.go



################################################################################
# compiling go sources - for test1
################################################################################

go/src/hydrocul/ichigo-test1:
	mkdir -p go/src/hydrocul/ichigo-test1

go/src/hydrocul/ichigo-test1/da.go: da.go go/src/hydrocul/ichigo-test1
	cp da.go go/src/hydrocul/ichigo-test1/da.go

go/src/hydrocul/ichigo-test1/da_test.go: da_test.go go/src/hydrocul/ichigo-test1
	cp da_test.go go/src/hydrocul/ichigo-test1/da_test.go

go/src/hydrocul/ichigo-test1/dict.go: dict.go go/src/hydrocul/ichigo-test1
	cp dict.go go/src/hydrocul/ichigo-test1/dict.go

go/src/hydrocul/ichigo-test1/dict_test.go: dict_test.go go/src/hydrocul/ichigo-test1
	cp dict_test.go go/src/hydrocul/ichigo-test1/dict_test.go

go/src/hydrocul/ichigo-test1/texts.go: texts.go go/src/hydrocul/ichigo-test1
	cp texts.go go/src/hydrocul/ichigo-test1/texts.go

go/src/hydrocul/ichigo-test1/texts_test.go: texts_test.go go/src/hydrocul/ichigo-test1
	cp texts_test.go go/src/hydrocul/ichigo-test1/texts_test.go



################################################################################
# compiling go sources - for building dictionary
################################################################################

go/bin/ichigo-build: go/src/hydrocul/ichigo-build/main.go go/src/hydrocul/ichigo-build/da.go go/src/hydrocul/ichigo-build/dict.go go/src/hydrocul/ichigo-build/texts.go
	GOPATH=$(realpath .)/go go install hydrocul/ichigo-build

go/src/hydrocul/ichigo-build:
	mkdir -p go/src/hydrocul/ichigo-build

go/src/hydrocul/ichigo-build/main.go: build_main.go go/src/hydrocul/ichigo-build
	cp build_main.go go/src/hydrocul/ichigo-build/main.go

go/src/hydrocul/ichigo-build/da.go: da.go go/src/hydrocul/ichigo-build
	cp da.go go/src/hydrocul/ichigo-build/da.go

go/src/hydrocul/ichigo-build/dict.go: dict.go go/src/hydrocul/ichigo-build
	cp dict.go go/src/hydrocul/ichigo-build/dict.go

go/src/hydrocul/ichigo-build/texts.go: texts.go go/src/hydrocul/ichigo-build
	cp texts.go go/src/hydrocul/ichigo-build/texts.go



################################################################################
# compiling go sources - for test2
################################################################################

go/src/hydrocul/ichigo-test2:
	mkdir -p go/src/hydrocul/ichigo-test2

go/src/hydrocul/ichigo-test2/da.go: da.go go/src/hydrocul/ichigo-test2
	cp da.go go/src/hydrocul/ichigo-test2/da.go

go/src/hydrocul/ichigo-test2/dict.go: dict.go go/src/hydrocul/ichigo-test2
	cp dict.go go/src/hydrocul/ichigo-test2/dict.go

go/src/hydrocul/ichigo-test2/dict_data.go: var/dict_data.go go/src/hydrocul/ichigo-test2
	cp var/dict_data.go go/src/hydrocul/ichigo-test2/dict_data.go

go/src/hydrocul/ichigo-test2/texts.go: texts.go go/src/hydrocul/ichigo-test2
	cp texts.go go/src/hydrocul/ichigo-test2/texts.go

go/src/hydrocul/ichigo-test2/pipe.go: pipe.go go/src/hydrocul/ichigo-test2
	cp pipe.go go/src/hydrocul/ichigo-test2/pipe.go

go/src/hydrocul/ichigo-test2/pipe_test.go: pipe_test.go go/src/hydrocul/ichigo-test2
	cp pipe_test.go go/src/hydrocul/ichigo-test2/pipe_test.go

go/src/hydrocul/ichigo-test2/posid.go: posid.go go/src/hydrocul/ichigo-test2
	cp posid.go go/src/hydrocul/ichigo-test2/posid.go



################################################################################
# compiling go sources - for main binary
################################################################################

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



################################################################################


