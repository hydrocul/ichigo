
################################################################################
# main target
################################################################################

default: ichigo-ipadic ichigo-unidic

ichigo-ipadic: go/bin/ichigo-ipadic
	cp go/bin/ichigo-ipadic ichigo-ipadic

ichigo-unidic: go/bin/ichigo-unidic
	cp go/bin/ichigo-unidic ichigo-unidic

all: test ichigo-ipadic ichigo-unidic

test: test1 test-ipadic test-unidic bench-ipadic bench-unidic



################################################################################
# etc
################################################################################

var/download/mkdir:
	mkdir -p var/download
	touch var/download/mkdir



################################################################################
# building ipadic dictionary
################################################################################

var/ipadic/mkdir:
	mkdir -p var/ipadic
	touch var/ipadic/mkdir

var/download/ipadic.touch: var/download/mkdir dict/ipadic/download.sh
	sh dict/ipadic/download.sh
	touch var/download/ipadic.touch

var/ipadic/matrix.txt: var/ipadic/mkdir var/download/ipadic.touch dict/mecab-matrix.sh
	sh dict/mecab-matrix.sh var/download/ipadic > var/ipadic/matrix.txt

var/ipadic/normalized.txt: var/ipadic/mkdir var/download/ipadic.touch dict/mecab-normalize.sh
	sh dict/mecab-normalize.sh var/download/ipadic > var/ipadic/normalized.txt

var/ipadic/formatted.txt: var/ipadic/mkdir var/ipadic/normalized.txt dict/ipadic/format.sh
	cat var/ipadic/normalized.txt | sh dict/ipadic/format.sh > var/ipadic/formatted.txt

var/ipadic/dict.txt: var/ipadic/mkdir var/ipadic/formatted.txt
	cat var/ipadic/formatted.txt | LC_ALL=C sort > var/ipadic/dict.txt

var/ipadic/texts.txt: var/ipadic/mkdir var/ipadic/dict.txt dict/ipadic/texts.sh
	cat var/ipadic/dict.txt | sh dict/ipadic/texts.sh > var/ipadic/texts.txt

var/ipadic/dict.dat: var/ipadic/mkdir var/ipadic/matrix.txt var/ipadic/dict.txt var/ipadic/texts.txt go/bin/ichigo-build-ipadic
	#go/bin/ichigo-build-ipadic var/ipadic/matrix.txt var/ipadic/dict.txt var/ipadic/texts.txt > var/ipadic/dict.dat.tmp
	go/bin/ichigo-build-ipadic var/ipadic/matrix.txt var/ipadic/dict.txt var/ipadic/texts.txt | gzip - > var/ipadic/dict.dat.tmp
	mv var/ipadic/dict.dat.tmp var/ipadic/dict.dat

var/ipadic/dict_data.go: var/ipadic/mkdir dict_data.go var/ipadic/dict.dat dict/to-go-source.sh
	sh dict/to-go-source.sh var/ipadic/dict.dat > var/ipadic/dict_data.go.tmp
	mv var/ipadic/dict_data.go.tmp var/ipadic/dict_data.go



################################################################################
# building unidic dictionary
################################################################################

var/unidic/mkdir:
	mkdir -p var/unidic
	touch var/unidic/mkdir

var/download/unidic.touch: var/download/mkdir dict/unidic/download.sh
	sh dict/unidic/download.sh
	touch var/download/unidic.touch

var/unidic/matrix.txt: var/unidic/mkdir var/download/unidic.touch dict/mecab-matrix.sh
	sh dict/mecab-matrix.sh var/download/unidic > var/unidic/matrix.txt

var/unidic/normalized.txt: var/unidic/mkdir var/download/unidic.touch dict/mecab-normalize.sh
	sh dict/mecab-normalize.sh var/download/unidic > var/unidic/normalized.txt

var/unidic/formatted.txt: var/unidic/mkdir var/unidic/normalized.txt dict/unidic/format.sh
	cat var/unidic/normalized.txt | sh dict/unidic/format.sh > var/unidic/formatted.txt

var/unidic/dict.txt: var/unidic/mkdir var/unidic/formatted.txt
	cat var/unidic/formatted.txt | LC_ALL=C sort > var/unidic/dict.txt

var/unidic/texts.txt: var/unidic/mkdir var/unidic/dict.txt dict/unidic/texts.sh
	cat var/unidic/dict.txt | sh dict/unidic/texts.sh > var/unidic/texts.txt

var/unidic/dict.dat: var/unidic/mkdir var/unidic/matrix.txt var/unidic/dict.txt var/unidic/texts.txt go/bin/ichigo-build-unidic
	#go/bin/ichigo-build-unidic var/unidic/matrix.txt var/unidic/dict.txt var/unidic/texts.txt > var/unidic/dict.dat.tmp
	go/bin/ichigo-build-unidic var/unidic/matrix.txt var/unidic/dict.txt var/unidic/texts.txt | gzip - > var/unidic/dict.dat.tmp
	mv var/unidic/dict.dat.tmp var/unidic/dict.dat

var/unidic/dict_data.go: var/unidic/mkdir dict_data.go var/unidic/dict.dat dict/to-go-source.sh
	sh dict/to-go-source.sh var/unidic/dict.dat > var/unidic/dict_data.go.tmp
	mv var/unidic/dict_data.go.tmp var/unidic/dict_data.go



################################################################################
# compiling go sources - for test1
################################################################################

test1: \
	go/src/hydrocul/ichigo-test1/da.go \
	go/src/hydrocul/ichigo-test1/da_test.go \
	go/src/hydrocul/ichigo-test1/dict.go \
	go/src/hydrocul/ichigo-test1/dict_test.go \
	go/src/hydrocul/ichigo-test1/texts.go \
	go/src/hydrocul/ichigo-test1/texts_test.go \
	go/src/hydrocul/ichigo-test1/posid.go \
	go/src/hydrocul/ichigo-test1/common.go
	GOPATH=$(realpath .)/go go test hydrocul/ichigo-test1

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

go/src/hydrocul/ichigo-test1/posid.go: posid-test1.go go/src/hydrocul/ichigo-test1/mkdir
	cp posid-test1.go go/src/hydrocul/ichigo-test1/posid.go

go/src/hydrocul/ichigo-test1/common.go: common.go go/src/hydrocul/ichigo-test1/mkdir
	cp common.go go/src/hydrocul/ichigo-test1/common.go



################################################################################
# compiling go sources - for building ipadic dictionary
################################################################################

go/bin/ichigo-build-ipadic: \
	go/src/hydrocul/ichigo-build-ipadic/main.go \
	go/src/hydrocul/ichigo-build-ipadic/da.go \
	go/src/hydrocul/ichigo-build-ipadic/dict.go \
	go/src/hydrocul/ichigo-build-ipadic/texts.go \
	go/src/hydrocul/ichigo-build-ipadic/posid.go \
	go/src/hydrocul/ichigo-build-ipadic/common.go
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

go/src/hydrocul/ichigo-build-ipadic/common.go: common.go go/src/hydrocul/ichigo-build-ipadic/mkdir
	cp common.go go/src/hydrocul/ichigo-build-ipadic/common.go



################################################################################
# compiling go sources - for building unidic dictionary
################################################################################

go/bin/ichigo-build-unidic: \
	go/src/hydrocul/ichigo-build-unidic/main.go \
	go/src/hydrocul/ichigo-build-unidic/da.go \
	go/src/hydrocul/ichigo-build-unidic/dict.go \
	go/src/hydrocul/ichigo-build-unidic/texts.go \
	go/src/hydrocul/ichigo-build-unidic/posid.go \
	go/src/hydrocul/ichigo-build-unidic/common.go
	GOPATH=$(realpath .)/go go install hydrocul/ichigo-build-unidic

go/src/hydrocul/ichigo-build-unidic/mkdir:
	mkdir -p go/src/hydrocul/ichigo-build-unidic
	touch go/src/hydrocul/ichigo-build-unidic/mkdir

go/src/hydrocul/ichigo-build-unidic/main.go: build_main.go go/src/hydrocul/ichigo-build-unidic/mkdir
	cp build_main.go go/src/hydrocul/ichigo-build-unidic/main.go

go/src/hydrocul/ichigo-build-unidic/da.go: da.go go/src/hydrocul/ichigo-build-unidic/mkdir
	cp da.go go/src/hydrocul/ichigo-build-unidic/da.go

go/src/hydrocul/ichigo-build-unidic/dict.go: dict.go go/src/hydrocul/ichigo-build-unidic/mkdir
	cp dict.go go/src/hydrocul/ichigo-build-unidic/dict.go

go/src/hydrocul/ichigo-build-unidic/texts.go: texts.go go/src/hydrocul/ichigo-build-unidic/mkdir
	cp texts.go go/src/hydrocul/ichigo-build-unidic/texts.go

go/src/hydrocul/ichigo-build-unidic/posid.go: posid-unidic.go go/src/hydrocul/ichigo-build-unidic/mkdir
	cp posid-unidic.go go/src/hydrocul/ichigo-build-unidic/posid.go

go/src/hydrocul/ichigo-build-unidic/common.go: common.go go/src/hydrocul/ichigo-build-unidic/mkdir
	cp common.go go/src/hydrocul/ichigo-build-unidic/common.go



################################################################################
# compiling go sources - for test-ipadic
################################################################################

test-ipadic: \
	go/src/hydrocul/ichigo-test-ipadic/da.go \
	go/src/hydrocul/ichigo-test-ipadic/dict.go \
	go/src/hydrocul/ichigo-test-ipadic/dict_data.go \
	go/src/hydrocul/ichigo-test-ipadic/texts.go \
	go/src/hydrocul/ichigo-test-ipadic/pipe.go \
	go/src/hydrocul/ichigo-test-ipadic/pipe_test.go \
	go/src/hydrocul/ichigo-test-ipadic/pipe_lib.go \
	go/src/hydrocul/ichigo-test-ipadic/posid.go \
	go/src/hydrocul/ichigo-test-ipadic/common.go
	GOPATH=$(realpath .)/go go test hydrocul/ichigo-test-ipadic

bench-ipadic: \
	go/src/hydrocul/ichigo-test-ipadic/da.go \
	go/src/hydrocul/ichigo-test-ipadic/dict.go \
	go/src/hydrocul/ichigo-test-ipadic/dict_data.go \
	go/src/hydrocul/ichigo-test-ipadic/texts.go \
	go/src/hydrocul/ichigo-test-ipadic/pipe.go \
	go/src/hydrocul/ichigo-test-ipadic/pipe_test.go \
	go/src/hydrocul/ichigo-test-ipadic/pipe_lib.go \
	go/src/hydrocul/ichigo-test-ipadic/posid.go \
	go/src/hydrocul/ichigo-test-ipadic/common.go
	GOPATH=$(realpath .)/go go test hydrocul/ichigo-test-ipadic -run none -bench . -benchtime 3s -benchmem | tee var/bench-ipadic.txt

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

go/src/hydrocul/ichigo-test-ipadic/pipe_test.go: pipe_test-ipadic.go go/src/hydrocul/ichigo-test-ipadic/mkdir
	cp pipe_test-ipadic.go go/src/hydrocul/ichigo-test-ipadic/pipe_test.go

go/src/hydrocul/ichigo-test-ipadic/pipe_lib.go: pipe_test-lib.go go/src/hydrocul/ichigo-test-ipadic/mkdir
	cp pipe_test-lib.go go/src/hydrocul/ichigo-test-ipadic/pipe_lib.go

go/src/hydrocul/ichigo-test-ipadic/posid.go: posid-ipadic.go go/src/hydrocul/ichigo-test-ipadic/mkdir
	cp posid-ipadic.go go/src/hydrocul/ichigo-test-ipadic/posid.go

go/src/hydrocul/ichigo-test-ipadic/common.go: common.go go/src/hydrocul/ichigo-test-ipadic/mkdir
	cp common.go go/src/hydrocul/ichigo-test-ipadic/common.go



################################################################################
# compiling go sources - for test-unidic
################################################################################

test-unidic: \
	go/src/hydrocul/ichigo-test-unidic/da.go \
	go/src/hydrocul/ichigo-test-unidic/dict.go \
	go/src/hydrocul/ichigo-test-unidic/dict_data.go \
	go/src/hydrocul/ichigo-test-unidic/texts.go \
	go/src/hydrocul/ichigo-test-unidic/pipe.go \
	go/src/hydrocul/ichigo-test-unidic/pipe_test.go \
	go/src/hydrocul/ichigo-test-unidic/pipe_lib.go \
	go/src/hydrocul/ichigo-test-unidic/posid.go \
	go/src/hydrocul/ichigo-test-unidic/common.go
	GOPATH=$(realpath .)/go go test hydrocul/ichigo-test-unidic

bench-unidic: \
	go/src/hydrocul/ichigo-test-unidic/da.go \
	go/src/hydrocul/ichigo-test-unidic/dict.go \
	go/src/hydrocul/ichigo-test-unidic/dict_data.go \
	go/src/hydrocul/ichigo-test-unidic/texts.go \
	go/src/hydrocul/ichigo-test-unidic/pipe.go \
	go/src/hydrocul/ichigo-test-unidic/pipe_test.go \
	go/src/hydrocul/ichigo-test-unidic/pipe_lib.go \
	go/src/hydrocul/ichigo-test-unidic/posid.go \
	go/src/hydrocul/ichigo-test-unidic/common.go
	GOPATH=$(realpath .)/go go test hydrocul/ichigo-test-unidic -run none -bench . -benchtime 3s -benchmem | tee var/bench-unidic.txt

go/src/hydrocul/ichigo-test-unidic/mkdir:
	mkdir -p go/src/hydrocul/ichigo-test-unidic
	touch go/src/hydrocul/ichigo-test-unidic/mkdir

go/src/hydrocul/ichigo-test-unidic/da.go: da.go go/src/hydrocul/ichigo-test-unidic/mkdir
	cp da.go go/src/hydrocul/ichigo-test-unidic/da.go

go/src/hydrocul/ichigo-test-unidic/dict.go: dict.go go/src/hydrocul/ichigo-test-unidic/mkdir
	cp dict.go go/src/hydrocul/ichigo-test-unidic/dict.go

go/src/hydrocul/ichigo-test-unidic/dict_data.go: var/unidic/dict_data.go go/src/hydrocul/ichigo-test-unidic/mkdir
	cp var/unidic/dict_data.go go/src/hydrocul/ichigo-test-unidic/dict_data.go

go/src/hydrocul/ichigo-test-unidic/texts.go: texts.go go/src/hydrocul/ichigo-test-unidic/mkdir
	cp texts.go go/src/hydrocul/ichigo-test-unidic/texts.go

go/src/hydrocul/ichigo-test-unidic/pipe.go: pipe.go go/src/hydrocul/ichigo-test-unidic/mkdir
	cp pipe.go go/src/hydrocul/ichigo-test-unidic/pipe.go

go/src/hydrocul/ichigo-test-unidic/pipe_test.go: pipe_test-unidic.go go/src/hydrocul/ichigo-test-unidic/mkdir
	cp pipe_test-unidic.go go/src/hydrocul/ichigo-test-unidic/pipe_test.go

go/src/hydrocul/ichigo-test-unidic/pipe_lib.go: pipe_test-lib.go go/src/hydrocul/ichigo-test-unidic/mkdir
	cp pipe_test-lib.go go/src/hydrocul/ichigo-test-unidic/pipe_lib.go

go/src/hydrocul/ichigo-test-unidic/posid.go: posid-unidic.go go/src/hydrocul/ichigo-test-unidic/mkdir
	cp posid-unidic.go go/src/hydrocul/ichigo-test-unidic/posid.go

go/src/hydrocul/ichigo-test-unidic/common.go: common.go go/src/hydrocul/ichigo-test-unidic/mkdir
	cp common.go go/src/hydrocul/ichigo-test-unidic/common.go



################################################################################
# compiling go sources - for main binary ichigo-ipadic
################################################################################

go/bin/ichigo-ipadic: \
	go/src/hydrocul/ichigo-ipadic/main.go \
	go/src/hydrocul/ichigo-ipadic/da.go \
	go/src/hydrocul/ichigo-ipadic/dict.go \
	go/src/hydrocul/ichigo-ipadic/dict_data.go \
	go/src/hydrocul/ichigo-ipadic/texts.go \
	go/src/hydrocul/ichigo-ipadic/pipe.go \
	go/src/hydrocul/ichigo-ipadic/posid.go \
	go/src/hydrocul/ichigo-ipadic/common.go
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

go/src/hydrocul/ichigo-ipadic/common.go: common.go go/src/hydrocul/ichigo-ipadic/mkdir
	cp common.go go/src/hydrocul/ichigo-ipadic/common.go



################################################################################
# compiling go sources - for main binary ichigo-unidic
################################################################################

go/bin/ichigo-unidic: \
	go/src/hydrocul/ichigo-unidic/main.go \
	go/src/hydrocul/ichigo-unidic/da.go \
	go/src/hydrocul/ichigo-unidic/dict.go \
	go/src/hydrocul/ichigo-unidic/dict_data.go \
	go/src/hydrocul/ichigo-unidic/texts.go \
	go/src/hydrocul/ichigo-unidic/pipe.go \
	go/src/hydrocul/ichigo-unidic/posid.go \
	go/src/hydrocul/ichigo-unidic/common.go
	GOPATH=$(realpath .)/go go install hydrocul/ichigo-unidic

go/src/hydrocul/ichigo-unidic/mkdir:
	mkdir -p go/src/hydrocul/ichigo-unidic
	touch go/src/hydrocul/ichigo-unidic/mkdir

go/src/hydrocul/ichigo-unidic/main.go: main.go go/src/hydrocul/ichigo-unidic/mkdir
	cp main.go go/src/hydrocul/ichigo-unidic/main.go

go/src/hydrocul/ichigo-unidic/da.go: da.go go/src/hydrocul/ichigo-unidic/mkdir
	cp da.go go/src/hydrocul/ichigo-unidic/da.go

go/src/hydrocul/ichigo-unidic/dict.go: dict.go go/src/hydrocul/ichigo-unidic/mkdir
	cp dict.go go/src/hydrocul/ichigo-unidic/dict.go

go/src/hydrocul/ichigo-unidic/dict_data.go: var/unidic/dict_data.go go/src/hydrocul/ichigo-unidic/mkdir
	cp var/unidic/dict_data.go go/src/hydrocul/ichigo-unidic/dict_data.go

go/src/hydrocul/ichigo-unidic/texts.go: texts.go go/src/hydrocul/ichigo-unidic/mkdir
	cp texts.go go/src/hydrocul/ichigo-unidic/texts.go

go/src/hydrocul/ichigo-unidic/pipe.go: pipe.go go/src/hydrocul/ichigo-unidic/mkdir
	cp pipe.go go/src/hydrocul/ichigo-unidic/pipe.go

go/src/hydrocul/ichigo-unidic/posid.go: posid-unidic.go go/src/hydrocul/ichigo-unidic/mkdir
	cp posid-unidic.go go/src/hydrocul/ichigo-unidic/posid.go

go/src/hydrocul/ichigo-unidic/common.go: common.go go/src/hydrocul/ichigo-unidic/mkdir
	cp common.go go/src/hydrocul/ichigo-unidic/common.go



################################################################################


