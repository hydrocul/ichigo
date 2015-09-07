
################################################################################
# main target
################################################################################

default: ichigo-ipadic ichigo-unidic

ichigo-ipadic: var/ipadic/boot.pl
	cp var/ipadic/boot.pl ichigo-ipadic
	chmod 755 ichigo-ipadic

ichigo-unidic: var/unidic/boot.pl
	cp var/unidic/boot.pl ichigo-unidic
	chmod 755 ichigo-unidic

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

var/ipadic/dict-normal.txt: var/ipadic/mkdir var/ipadic/formatted.txt
	cat var/ipadic/formatted.txt | LC_ALL=C sort > var/ipadic/dict-normal.txt

var/ipadic/texts.txt: var/ipadic/mkdir var/ipadic/dict-normal.txt dict/ipadic/texts.sh
	cat var/ipadic/dict-normal.txt | sh dict/ipadic/texts.sh > var/ipadic/texts.txt

var/ipadic/dict.dat: var/ipadic/mkdir var/ipadic/matrix.txt var/ipadic/dict-normal.txt var/ipadic/texts.txt go/bin/ichigo-build-ipadic
	go/bin/ichigo-build-ipadic var/ipadic/matrix.txt var/ipadic/texts.txt var/ipadic/dict-normal.txt > var/ipadic/dict.dat.tmp
	mv var/ipadic/dict.dat.tmp var/ipadic/dict.dat

var/ipadic/main: go/bin/ichigo-ipadic
	cp go/bin/ichigo-ipadic var/ipadic/main

var/ipadic/boot.pl: var/ipadic/mkdir boot.pl var/ipadic/dict.dat var/ipadic/main generate-boot.sh
	sh generate-boot.sh var/ipadic > var/ipadic/boot.pl



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

var/unidic/dict-normal.txt: var/unidic/mkdir var/unidic/formatted.txt
	cat var/unidic/formatted.txt | LC_ALL=C sort > var/unidic/dict-normal.txt

var/unidic/texts.txt: var/unidic/mkdir var/unidic/dict-normal.txt dict/unidic/texts.sh
	cat var/unidic/dict-normal.txt | sh dict/unidic/texts.sh > var/unidic/texts.txt

var/unidic/dict.dat: var/unidic/mkdir var/unidic/matrix.txt var/unidic/dict-normal.txt var/unidic/texts.txt go/bin/ichigo-build-unidic
	go/bin/ichigo-build-unidic var/unidic/matrix.txt var/unidic/texts.txt var/unidic/dict-normal.txt > var/unidic/dict.dat.tmp
	mv var/unidic/dict.dat.tmp var/unidic/dict.dat

var/unidic/main: go/bin/ichigo-unidic
	cp go/bin/ichigo-unidic var/unidic/main

var/unidic/boot.pl: var/unidic/mkdir boot.pl var/unidic/dict.dat var/unidic/main generate-boot.sh
	sh generate-boot.sh var/unidic > var/unidic/boot.pl



################################################################################
# compiling go sources - for test1
################################################################################

test1: \
	go/src/hydrocul/ichigo-test1/da.go \
	go/src/hydrocul/ichigo-test1/da_test.go \
	go/src/hydrocul/ichigo-test1/utf8.go \
	go/src/hydrocul/ichigo-test1/dict.go \
	go/src/hydrocul/ichigo-test1/dict_test.go \
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

go/src/hydrocul/ichigo-test1/utf8.go: utf8.go go/src/hydrocul/ichigo-test1/mkdir
	cp utf8.go go/src/hydrocul/ichigo-test1/utf8.go

go/src/hydrocul/ichigo-test1/dict.go: dict.go go/src/hydrocul/ichigo-test1/mkdir
	cp dict.go go/src/hydrocul/ichigo-test1/dict.go

go/src/hydrocul/ichigo-test1/dict_test.go: dict_test.go go/src/hydrocul/ichigo-test1/mkdir
	cp dict_test.go go/src/hydrocul/ichigo-test1/dict_test.go

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
	go/src/hydrocul/ichigo-build-ipadic/utf8.go \
	go/src/hydrocul/ichigo-build-ipadic/dict.go \
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

go/src/hydrocul/ichigo-build-ipadic/utf8.go: utf8.go go/src/hydrocul/ichigo-build-ipadic/mkdir
	cp utf8.go go/src/hydrocul/ichigo-build-ipadic/utf8.go

go/src/hydrocul/ichigo-build-ipadic/dict.go: dict.go go/src/hydrocul/ichigo-build-ipadic/mkdir
	cp dict.go go/src/hydrocul/ichigo-build-ipadic/dict.go

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
	go/src/hydrocul/ichigo-build-unidic/utf8.go \
	go/src/hydrocul/ichigo-build-unidic/dict.go \
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

go/src/hydrocul/ichigo-build-unidic/utf8.go: utf8.go go/src/hydrocul/ichigo-build-unidic/mkdir
	cp utf8.go go/src/hydrocul/ichigo-build-unidic/utf8.go

go/src/hydrocul/ichigo-build-unidic/dict.go: dict.go go/src/hydrocul/ichigo-build-unidic/mkdir
	cp dict.go go/src/hydrocul/ichigo-build-unidic/dict.go

go/src/hydrocul/ichigo-build-unidic/posid.go: posid-unidic.go go/src/hydrocul/ichigo-build-unidic/mkdir
	cp posid-unidic.go go/src/hydrocul/ichigo-build-unidic/posid.go

go/src/hydrocul/ichigo-build-unidic/common.go: common.go go/src/hydrocul/ichigo-build-unidic/mkdir
	cp common.go go/src/hydrocul/ichigo-build-unidic/common.go



################################################################################
# compiling go sources - for test-ipadic
################################################################################

test-ipadic: \
	go/src/hydrocul/ichigo-test-ipadic/da.go \
	go/src/hydrocul/ichigo-test-ipadic/utf8.go \
	go/src/hydrocul/ichigo-test-ipadic/dict.go \
	go/src/hydrocul/ichigo-test-ipadic/dict_data.go \
	go/src/hydrocul/ichigo-test-ipadic/pipe.go \
	go/src/hydrocul/ichigo-test-ipadic/pipe_test.go \
	go/src/hydrocul/ichigo-test-ipadic/pipe_lib.go \
	go/src/hydrocul/ichigo-test-ipadic/shift.go \
	go/src/hydrocul/ichigo-test-ipadic/posid.go \
	go/src/hydrocul/ichigo-test-ipadic/common.go \
	var/ipadic/dict.dat
	GOPATH=$(realpath .)/go ICHIGO_DICTIONARY_PATH=$(realpath .)/var/ipadic/dict.dat go test hydrocul/ichigo-test-ipadic

bench-ipadic: \
	go/src/hydrocul/ichigo-test-ipadic/da.go \
	go/src/hydrocul/ichigo-test-ipadic/utf8.go \
	go/src/hydrocul/ichigo-test-ipadic/dict.go \
	go/src/hydrocul/ichigo-test-ipadic/dict_data.go \
	go/src/hydrocul/ichigo-test-ipadic/pipe.go \
	go/src/hydrocul/ichigo-test-ipadic/pipe_test.go \
	go/src/hydrocul/ichigo-test-ipadic/pipe_lib.go \
	go/src/hydrocul/ichigo-test-ipadic/shift.go \
	go/src/hydrocul/ichigo-test-ipadic/posid.go \
	go/src/hydrocul/ichigo-test-ipadic/common.go \
	var/ipadic/dict.dat
	GOPATH=$(realpath .)/go ICHIGO_DICTIONARY_PATH=$(realpath .)/var/ipadic/dict.dat go test hydrocul/ichigo-test-ipadic -run none -bench . -benchtime 3s -benchmem | tee var/bench-ipadic.txt

go/src/hydrocul/ichigo-test-ipadic/mkdir:
	mkdir -p go/src/hydrocul/ichigo-test-ipadic
	touch go/src/hydrocul/ichigo-test-ipadic/mkdir

go/src/hydrocul/ichigo-test-ipadic/da.go: da.go go/src/hydrocul/ichigo-test-ipadic/mkdir
	cp da.go go/src/hydrocul/ichigo-test-ipadic/da.go

go/src/hydrocul/ichigo-test-ipadic/utf8.go: utf8.go go/src/hydrocul/ichigo-test-ipadic/mkdir
	cp utf8.go go/src/hydrocul/ichigo-test-ipadic/utf8.go

go/src/hydrocul/ichigo-test-ipadic/dict.go: dict.go go/src/hydrocul/ichigo-test-ipadic/mkdir
	cp dict.go go/src/hydrocul/ichigo-test-ipadic/dict.go

go/src/hydrocul/ichigo-test-ipadic/dict_data.go: dict_data.go go/src/hydrocul/ichigo-test-ipadic/mkdir
	cp dict_data.go go/src/hydrocul/ichigo-test-ipadic/dict_data.go

go/src/hydrocul/ichigo-test-ipadic/pipe.go: pipe.go go/src/hydrocul/ichigo-test-ipadic/mkdir
	cp pipe.go go/src/hydrocul/ichigo-test-ipadic/pipe.go

go/src/hydrocul/ichigo-test-ipadic/pipe_test.go: pipe_test-ipadic.go go/src/hydrocul/ichigo-test-ipadic/mkdir
	cp pipe_test-ipadic.go go/src/hydrocul/ichigo-test-ipadic/pipe_test.go

go/src/hydrocul/ichigo-test-ipadic/pipe_lib.go: pipe_test-lib.go go/src/hydrocul/ichigo-test-ipadic/mkdir
	cp pipe_test-lib.go go/src/hydrocul/ichigo-test-ipadic/pipe_lib.go

go/src/hydrocul/ichigo-test-ipadic/shift.go: shift.go go/src/hydrocul/ichigo-test-ipadic/mkdir
	cp shift.go go/src/hydrocul/ichigo-test-ipadic/shift.go

go/src/hydrocul/ichigo-test-ipadic/posid.go: posid-ipadic.go go/src/hydrocul/ichigo-test-ipadic/mkdir
	cp posid-ipadic.go go/src/hydrocul/ichigo-test-ipadic/posid.go

go/src/hydrocul/ichigo-test-ipadic/common.go: common.go go/src/hydrocul/ichigo-test-ipadic/mkdir
	cp common.go go/src/hydrocul/ichigo-test-ipadic/common.go



################################################################################
# compiling go sources - for test-unidic
################################################################################

test-unidic: \
	go/src/hydrocul/ichigo-test-unidic/da.go \
	go/src/hydrocul/ichigo-test-unidic/utf8.go \
	go/src/hydrocul/ichigo-test-unidic/dict.go \
	go/src/hydrocul/ichigo-test-unidic/dict_data.go \
	go/src/hydrocul/ichigo-test-unidic/pipe.go \
	go/src/hydrocul/ichigo-test-unidic/pipe_test.go \
	go/src/hydrocul/ichigo-test-unidic/pipe_lib.go \
	go/src/hydrocul/ichigo-test-unidic/shift.go \
	go/src/hydrocul/ichigo-test-unidic/posid.go \
	go/src/hydrocul/ichigo-test-unidic/common.go \
	var/unidic/dict.dat
	GOPATH=$(realpath .)/go ICHIGO_DICTIONARY_PATH=$(realpath .)/var/unidic/dict.dat go test hydrocul/ichigo-test-unidic

bench-unidic: \
	go/src/hydrocul/ichigo-test-unidic/da.go \
	go/src/hydrocul/ichigo-test-unidic/utf8.go \
	go/src/hydrocul/ichigo-test-unidic/dict.go \
	go/src/hydrocul/ichigo-test-unidic/dict_data.go \
	go/src/hydrocul/ichigo-test-unidic/pipe.go \
	go/src/hydrocul/ichigo-test-unidic/pipe_test.go \
	go/src/hydrocul/ichigo-test-unidic/pipe_lib.go \
	go/src/hydrocul/ichigo-test-unidic/shift.go \
	go/src/hydrocul/ichigo-test-unidic/posid.go \
	go/src/hydrocul/ichigo-test-unidic/common.go \
	var/unidic/dict.dat
	GOPATH=$(realpath .)/go ICHIGO_DICTIONARY_PATH=$(realpath .)/var/unidic/dict.dat go test hydrocul/ichigo-test-unidic -run none -bench . -benchtime 3s -benchmem | tee var/bench-unidic.txt

go/src/hydrocul/ichigo-test-unidic/mkdir:
	mkdir -p go/src/hydrocul/ichigo-test-unidic
	touch go/src/hydrocul/ichigo-test-unidic/mkdir

go/src/hydrocul/ichigo-test-unidic/da.go: da.go go/src/hydrocul/ichigo-test-unidic/mkdir
	cp da.go go/src/hydrocul/ichigo-test-unidic/da.go

go/src/hydrocul/ichigo-test-unidic/utf8.go: utf8.go go/src/hydrocul/ichigo-test-unidic/mkdir
	cp utf8.go go/src/hydrocul/ichigo-test-unidic/utf8.go

go/src/hydrocul/ichigo-test-unidic/dict.go: dict.go go/src/hydrocul/ichigo-test-unidic/mkdir
	cp dict.go go/src/hydrocul/ichigo-test-unidic/dict.go

go/src/hydrocul/ichigo-test-unidic/dict_data.go: dict_data.go go/src/hydrocul/ichigo-test-unidic/mkdir
	cp dict_data.go go/src/hydrocul/ichigo-test-unidic/dict_data.go

go/src/hydrocul/ichigo-test-unidic/pipe.go: pipe.go go/src/hydrocul/ichigo-test-unidic/mkdir
	cp pipe.go go/src/hydrocul/ichigo-test-unidic/pipe.go

go/src/hydrocul/ichigo-test-unidic/pipe_test.go: pipe_test-unidic.go go/src/hydrocul/ichigo-test-unidic/mkdir
	cp pipe_test-unidic.go go/src/hydrocul/ichigo-test-unidic/pipe_test.go

go/src/hydrocul/ichigo-test-unidic/pipe_lib.go: pipe_test-lib.go go/src/hydrocul/ichigo-test-unidic/mkdir
	cp pipe_test-lib.go go/src/hydrocul/ichigo-test-unidic/pipe_lib.go

go/src/hydrocul/ichigo-test-unidic/shift.go: shift.go go/src/hydrocul/ichigo-test-unidic/mkdir
	cp shift.go go/src/hydrocul/ichigo-test-unidic/shift.go

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
	go/src/hydrocul/ichigo-ipadic/utf8.go \
	go/src/hydrocul/ichigo-ipadic/dict.go \
	go/src/hydrocul/ichigo-ipadic/dict_data.go \
	go/src/hydrocul/ichigo-ipadic/pipe.go \
	go/src/hydrocul/ichigo-ipadic/shift.go \
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

go/src/hydrocul/ichigo-ipadic/utf8.go: utf8.go go/src/hydrocul/ichigo-ipadic/mkdir
	cp utf8.go go/src/hydrocul/ichigo-ipadic/utf8.go

go/src/hydrocul/ichigo-ipadic/dict.go: dict.go go/src/hydrocul/ichigo-ipadic/mkdir
	cp dict.go go/src/hydrocul/ichigo-ipadic/dict.go

go/src/hydrocul/ichigo-ipadic/dict_data.go: dict_data.go go/src/hydrocul/ichigo-ipadic/mkdir
	cp dict_data.go go/src/hydrocul/ichigo-ipadic/dict_data.go

go/src/hydrocul/ichigo-ipadic/pipe.go: pipe.go go/src/hydrocul/ichigo-ipadic/mkdir
	cp pipe.go go/src/hydrocul/ichigo-ipadic/pipe.go

go/src/hydrocul/ichigo-ipadic/shift.go: shift.go go/src/hydrocul/ichigo-ipadic/mkdir
	cp shift.go go/src/hydrocul/ichigo-ipadic/shift.go

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
	go/src/hydrocul/ichigo-unidic/utf8.go \
	go/src/hydrocul/ichigo-unidic/dict.go \
	go/src/hydrocul/ichigo-unidic/dict_data.go \
	go/src/hydrocul/ichigo-unidic/pipe.go \
	go/src/hydrocul/ichigo-unidic/shift.go \
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

go/src/hydrocul/ichigo-unidic/utf8.go: utf8.go go/src/hydrocul/ichigo-unidic/mkdir
	cp utf8.go go/src/hydrocul/ichigo-unidic/utf8.go

go/src/hydrocul/ichigo-unidic/dict.go: dict.go go/src/hydrocul/ichigo-unidic/mkdir
	cp dict.go go/src/hydrocul/ichigo-unidic/dict.go

go/src/hydrocul/ichigo-unidic/dict_data.go: dict_data.go go/src/hydrocul/ichigo-unidic/mkdir
	cp dict_data.go go/src/hydrocul/ichigo-unidic/dict_data.go

go/src/hydrocul/ichigo-unidic/pipe.go: pipe.go go/src/hydrocul/ichigo-unidic/mkdir
	cp pipe.go go/src/hydrocul/ichigo-unidic/pipe.go

go/src/hydrocul/ichigo-unidic/shift.go: shift.go go/src/hydrocul/ichigo-unidic/mkdir
	cp shift.go go/src/hydrocul/ichigo-unidic/shift.go

go/src/hydrocul/ichigo-unidic/posid.go: posid-unidic.go go/src/hydrocul/ichigo-unidic/mkdir
	cp posid-unidic.go go/src/hydrocul/ichigo-unidic/posid.go

go/src/hydrocul/ichigo-unidic/common.go: common.go go/src/hydrocul/ichigo-unidic/mkdir
	cp common.go go/src/hydrocul/ichigo-unidic/common.go



################################################################################


