
################################################################################
# main target
################################################################################

ichigo-ipadic: go/bin/ichigo-ipadic
	cp go/bin/ichigo-ipadic ichigo-ipadic

all: test ichigo-ipadic


################################################################################
# test scripts
################################################################################

test: test1 test-ipadic bench-ipadic

test1: go/src/hydrocul/ichigo-test1/da.go go/src/hydrocul/ichigo-test1/da_test.go go/src/hydrocul/ichigo-test1/dict.go go/src/hydrocul/ichigo-test1/dict_test.go go/src/hydrocul/ichigo-test1/texts.go go/src/hydrocul/ichigo-test1/texts_test.go
	GOPATH=$(realpath .)/go go test hydrocul/ichigo-test1

test-ipadic: go/src/hydrocul/ichigo-test-ipadic/da.go go/src/hydrocul/ichigo-test-ipadic/dict.go go/src/hydrocul/ichigo-test-ipadic/dict_data.go go/src/hydrocul/ichigo-test-ipadic/texts.go go/src/hydrocul/ichigo-test-ipadic/pipe.go go/src/hydrocul/ichigo-test-ipadic/pipe_test.go go/src/hydrocul/ichigo-test-ipadic/posid.go
	GOPATH=$(realpath .)/go go test hydrocul/ichigo-test-ipadic

bench-ipadic: go/src/hydrocul/ichigo-test-ipadic/da.go go/src/hydrocul/ichigo-test-ipadic/dict.go go/src/hydrocul/ichigo-test-ipadic/dict_data.go go/src/hydrocul/ichigo-test-ipadic/texts.go go/src/hydrocul/ichigo-test-ipadic/pipe.go go/src/hydrocul/ichigo-test-ipadic/pipe_test.go go/src/hydrocul/ichigo-test-ipadic/posid.go
	GOPATH=$(realpath .)/go go test hydrocul/ichigo-test-ipadic -run none -bench . -benchtime 10s -benchmem



################################################################################
# making var directory
################################################################################

var/ipadic/mkdir:
	mkdir -p var/ipadic
	touch var/ipadic/mkdir



################################################################################
# building dictionary
################################################################################

var/ipadic/download.touch: var/ipadic/mkdir
	sh dict/ipadic/download.sh
	touch var/ipadic/download.touch

var/ipadic/matrix.txt: var/ipadic/mkdir var/ipadic/download.touch dict/mecab-matrix.sh
	sh dict/mecab-matrix.sh var/ipadic/download > var/ipadic/matrix.txt

var/ipadic/words-1.txt: var/ipadic/mkdir var/ipadic/download.touch dict/mecab-words-1.sh
	sh dict/mecab-words-1.sh var/ipadic/download > var/ipadic/words-1.txt

var/ipadic/words-2.txt: var/ipadic/mkdir var/ipadic/words-1.txt
	cat var/ipadic/words-1.txt | LC_ALL=C sort > var/ipadic/words-2.txt

var/ipadic/words-3.txt: var/ipadic/mkdir var/ipadic/words-2.txt dict/ipadic/words-3.sh
	cat var/ipadic/words-2.txt | sh dict/ipadic/words-3.sh > var/ipadic/words-3.txt

var/ipadic/dict.txt: var/ipadic/mkdir var/ipadic/words-3.txt
	cp var/ipadic/words-3.txt var/ipadic/dict.txt

var/ipadic/texts.txt: var/ipadic/mkdir var/ipadic/dict.txt dict/texts.sh
	cat var/ipadic/dict.txt | sh dict/texts.sh > var/ipadic/texts.txt

var/ipadic/dict.dat: var/ipadic/mkdir var/ipadic/matrix.txt var/ipadic/dict.txt var/ipadic/texts.txt go/bin/ichigo-build-ipadic
	go/bin/ichigo-build-ipadic var/ipadic/matrix.txt var/ipadic/dict.txt var/ipadic/texts.txt > var/ipadic/dict.dat.tmp
	mv var/ipadic/dict.dat.tmp var/ipadic/dict.dat

var/ipadic/dict_data.go: var/ipadic/mkdir dict_data.go var/ipadic/dict.dat dict/to-go-source.sh
	sh dict/to-go-source.sh var/ipadic/dict.dat > var/ipadic/dict_data.go.tmp
	mv var/ipadic/dict_data.go.tmp var/ipadic/dict_data.go



################################################################################
# compiling go sources - for test1
################################################################################

go/src/hydrocul/ichigo-test1/mkdir:
	mkdir -p go/src/hydrocul/ichigo-test1
	touch go/src/hydrocul/ichigo-test1/mkdir

go/src/hydrocul/ichigo-test1/da.go: da.go go/src/hydrocul/ichigo-test1/mkdir
	cp da.go go/src/hydrocul/ichigo-test1/da.go

go/src/hydrocul/ichigo-test1/da_test.go: da_test.go go/src/hydrocul/ichigo-test1/mkdir
	cp da_test.go go/src/hydrocul/ichigo-test1/da_test.go

go/src/hydrocul/ichigo-test1/dict.go: dict.go go/src/hydrocul/ichigo-test1/mkdir
	cp dict.go go/src/hydrocul/ichigo-test1/dict.go

go/src/hydrocul/ichigo-test1/dict_test.go: dict_test.go go/src/hydrocul/ichigo-test1/mkdir
	cp dict_test.go go/src/hydrocul/ichigo-test1/dict_test.go

go/src/hydrocul/ichigo-test1/texts.go: texts.go go/src/hydrocul/ichigo-test1/mkdir
	cp texts.go go/src/hydrocul/ichigo-test1/texts.go

go/src/hydrocul/ichigo-test1/texts_test.go: texts_test.go go/src/hydrocul/ichigo-test1/mkdir
	cp texts_test.go go/src/hydrocul/ichigo-test1/texts_test.go



################################################################################
# compiling go sources - for building ipadic dictionary
################################################################################

go/bin/ichigo-build-ipadic: go/src/hydrocul/ichigo-build-ipadic/main.go go/src/hydrocul/ichigo-build-ipadic/da.go go/src/hydrocul/ichigo-build-ipadic/dict.go go/src/hydrocul/ichigo-build-ipadic/texts.go go/src/hydrocul/ichigo-build-ipadic/posid.go
	GOPATH=$(realpath .)/go go install hydrocul/ichigo-build-ipadic

go/src/hydrocul/ichigo-build-ipadic/mkdir:
	mkdir -p go/src/hydrocul/ichigo-build-ipadic
	touch go/src/hydrocul/ichigo-build-ipadic/mkdir

go/src/hydrocul/ichigo-build-ipadic/main.go: build_main.go go/src/hydrocul/ichigo-build-ipadic/mkdir
	cp build_main.go go/src/hydrocul/ichigo-build-ipadic/main.go

go/src/hydrocul/ichigo-build-ipadic/da.go: da.go go/src/hydrocul/ichigo-build-ipadic/mkdir
	cp da.go go/src/hydrocul/ichigo-build-ipadic/da.go

go/src/hydrocul/ichigo-build-ipadic/dict.go: dict.go go/src/hydrocul/ichigo-build-ipadic/mkdir
	cp dict.go go/src/hydrocul/ichigo-build-ipadic/dict.go

go/src/hydrocul/ichigo-build-ipadic/texts.go: texts.go go/src/hydrocul/ichigo-build-ipadic/mkdir
	cp texts.go go/src/hydrocul/ichigo-build-ipadic/texts.go

go/src/hydrocul/ichigo-build-ipadic/posid.go: posid-ipadic.go go/src/hydrocul/ichigo-build-ipadic/mkdir
	cp posid-ipadic.go go/src/hydrocul/ichigo-build-ipadic/posid.go



################################################################################
# compiling go sources - for test-ipadic
################################################################################

go/src/hydrocul/ichigo-test-ipadic/mkdir:
	mkdir -p go/src/hydrocul/ichigo-test-ipadic
	touch go/src/hydrocul/ichigo-test-ipadic/mkdir

go/src/hydrocul/ichigo-test-ipadic/da.go: da.go go/src/hydrocul/ichigo-test-ipadic/mkdir
	cp da.go go/src/hydrocul/ichigo-test-ipadic/da.go

go/src/hydrocul/ichigo-test-ipadic/dict.go: dict.go go/src/hydrocul/ichigo-test-ipadic/mkdir
	cp dict.go go/src/hydrocul/ichigo-test-ipadic/dict.go

go/src/hydrocul/ichigo-test-ipadic/dict_data.go: var/ipadic/dict_data.go go/src/hydrocul/ichigo-test-ipadic/mkdir
	cp var/ipadic/dict_data.go go/src/hydrocul/ichigo-test-ipadic/dict_data.go

go/src/hydrocul/ichigo-test-ipadic/texts.go: texts.go go/src/hydrocul/ichigo-test-ipadic/mkdir
	cp texts.go go/src/hydrocul/ichigo-test-ipadic/texts.go

go/src/hydrocul/ichigo-test-ipadic/pipe.go: pipe.go go/src/hydrocul/ichigo-test-ipadic/mkdir
	cp pipe.go go/src/hydrocul/ichigo-test-ipadic/pipe.go

go/src/hydrocul/ichigo-test-ipadic/pipe_test.go: pipe_test.go go/src/hydrocul/ichigo-test-ipadic/mkdir
	cp pipe_test.go go/src/hydrocul/ichigo-test-ipadic/pipe_test.go

go/src/hydrocul/ichigo-test-ipadic/posid.go: posid-ipadic.go go/src/hydrocul/ichigo-test-ipadic/mkdir
	cp posid-ipadic.go go/src/hydrocul/ichigo-test-ipadic/posid.go



################################################################################
# compiling go sources - for main binary ichigo-ipadic
################################################################################

go/bin/ichigo-ipadic: go/src/hydrocul/ichigo-ipadic/main.go go/src/hydrocul/ichigo-ipadic/da.go go/src/hydrocul/ichigo-ipadic/dict.go go/src/hydrocul/ichigo-ipadic/dict_data.go go/src/hydrocul/ichigo-ipadic/texts.go go/src/hydrocul/ichigo-ipadic/pipe.go go/src/hydrocul/ichigo-ipadic/posid.go
	GOPATH=$(realpath .)/go go install hydrocul/ichigo-ipadic

go/src/hydrocul/ichigo-ipadic/mkdir:
	mkdir -p go/src/hydrocul/ichigo-ipadic
	touch go/src/hydrocul/ichigo-ipadic/mkdir

go/src/hydrocul/ichigo-ipadic/main.go: main.go go/src/hydrocul/ichigo-ipadic/mkdir
	cp main.go go/src/hydrocul/ichigo-ipadic/main.go

go/src/hydrocul/ichigo-ipadic/da.go: da.go go/src/hydrocul/ichigo-ipadic/mkdir
	cp da.go go/src/hydrocul/ichigo-ipadic/da.go

go/src/hydrocul/ichigo-ipadic/dict.go: dict.go go/src/hydrocul/ichigo-ipadic/mkdir
	cp dict.go go/src/hydrocul/ichigo-ipadic/dict.go

go/src/hydrocul/ichigo-ipadic/dict_data.go: var/ipadic/dict_data.go go/src/hydrocul/ichigo-ipadic/mkdir
	cp var/ipadic/dict_data.go go/src/hydrocul/ichigo-ipadic/dict_data.go

go/src/hydrocul/ichigo-ipadic/texts.go: texts.go go/src/hydrocul/ichigo-ipadic/mkdir
	cp texts.go go/src/hydrocul/ichigo-ipadic/texts.go

go/src/hydrocul/ichigo-ipadic/pipe.go: pipe.go go/src/hydrocul/ichigo-ipadic/mkdir
	cp pipe.go go/src/hydrocul/ichigo-ipadic/pipe.go

go/src/hydrocul/ichigo-ipadic/posid.go: posid-ipadic.go go/src/hydrocul/ichigo-ipadic/mkdir
	cp posid-ipadic.go go/src/hydrocul/ichigo-ipadic/posid.go



################################################################################


