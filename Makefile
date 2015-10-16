
################################################################################
# main target
################################################################################

default: ichigo-ipadic ichigo-unidic

all: test ichigo-ipadic ichigo-unidic

ichigo-ipadic: var/ipadic/boot.pl
	cp var/ipadic/boot.pl ichigo-ipadic
	chmod 755 ichigo-ipadic

ichigo-unidic: var/unidic/boot.pl
	cp var/unidic/boot.pl ichigo-unidic
	chmod 755 ichigo-unidic

test: \
	test1 \
	test-ipadic \
	test-unidic \
	bench-ipadic \
	bench-unidic \
	wagahai-ipadic \
	wagahai-unidic



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

var/ipadic/graph.pl: graph.pl
	cp graph.pl var/ipadic/graph.pl

var/ipadic/boot.pl: var/ipadic/mkdir boot.pl var/ipadic/dict.dat var/ipadic/main var/ipadic/graph.pl generate-boot.sh
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

var/unidic/graph.pl: graph.pl
	cp graph.pl var/unidic/graph.pl

var/unidic/boot.pl: var/unidic/mkdir boot.pl graph.pl var/unidic/dict.dat var/unidic/main var/unidic/graph.pl generate-boot.sh
	sh generate-boot.sh var/unidic > var/unidic/boot.pl



################################################################################
# compiling go sources - for test1
################################################################################

test1: go/src/hydrocul/ichigo-test1/mkdir
	GOPATH=$(realpath .)/go go test hydrocul/ichigo-test1

go/src/hydrocul/ichigo-test1/mkdir: \
	da.go \
	da_test.go \
	utf8.go \
	dict.go \
	dict_test.go \
	posid.go \
	common.go \
	preprocess.pl
	mkdir -p go/src/hydrocul/ichigo-test1
	perl preprocess.pl test1 go/src/hydrocul/ichigo-test1 \
		da.go \
		da_test.go \
		utf8.go \
		dict.go \
		dict_test.go \
		posid.go \
		common.go



################################################################################
# compiling go sources - for building ipadic dictionary
################################################################################

go/bin/ichigo-build-ipadic: go/src/hydrocul/ichigo-build-ipadic/mkdir
	GOPATH=$(realpath .)/go go install hydrocul/ichigo-build-ipadic

go/src/hydrocul/ichigo-build-ipadic/mkdir: \
	build_main.go \
	da.go \
	utf8.go \
	dict.go \
	posid.go \
	common.go \
	preprocess.pl
	mkdir -p go/src/hydrocul/ichigo-build-ipadic
	perl preprocess.pl ipadic go/src/hydrocul/ichigo-build-ipadic \
		build_main.go \
		da.go \
		utf8.go \
		dict.go \
		posid.go \
		common.go



################################################################################
# compiling go sources - for building unidic dictionary
################################################################################

go/bin/ichigo-build-unidic: go/src/hydrocul/ichigo-build-unidic/mkdir
	GOPATH=$(realpath .)/go go install hydrocul/ichigo-build-unidic

go/src/hydrocul/ichigo-build-unidic/mkdir: \
	build_main.go \
	da.go \
	utf8.go \
	dict.go \
	posid.go \
	common.go \
	preprocess.pl
	mkdir -p go/src/hydrocul/ichigo-build-unidic
	perl preprocess.pl unidic go/src/hydrocul/ichigo-build-unidic \
		build_main.go \
		da.go \
		utf8.go \
		dict.go \
		posid.go \
		common.go



################################################################################
# compiling go sources - for test-ipadic
################################################################################

test-ipadic: go/src/hydrocul/ichigo-test-ipadic/mkdir var/ipadic/dict.dat
	GOPATH=$(realpath .)/go ICHIGO_DICTIONARY_PATH=$(realpath .)/var/ipadic/dict.dat go test hydrocul/ichigo-test-ipadic

var/ichigo-test-ipadic.test: go/src/hydrocul/ichigo-test-ipadic/mkdir
	GOPATH=$(realpath .)/go go test -c hydrocul/ichigo-test-ipadic -o var/ichigo-test-ipadic.test

bench-ipadic: var/ipadic/dict.dat var/ichigo-test-ipadic.test
	GOPATH=$(realpath .)/go ICHIGO_DICTIONARY_PATH=$(realpath .)/var/ipadic/dict.dat ./var/ichigo-test-ipadic.test -test.cpuprofile=var/cpuprofile-ipadic.out -test.bench . -test.benchtime 30s -test.benchmem | tee var/bench-ipadic.txt
	GOPATH=$(realpath .)/go go tool pprof -list=hydrocul var/ichigo-test-ipadic.test var/cpuprofile-ipadic.out > var/cpuprofile-ipadic-list.txt

go/src/hydrocul/ichigo-test-ipadic/mkdir: \
	main.go \
	da.go \
	utf8.go \
	dict.go \
	dict_data.go \
	pipe.go \
	pipe_test.go \
	shift.go \
	posid.go \
	common.go \
	preprocess.pl
	mkdir -p go/src/hydrocul/ichigo-test-ipadic
	perl preprocess.pl ipadic go/src/hydrocul/ichigo-test-ipadic \
		main.go \
		da.go \
		utf8.go \
		dict.go \
		dict_data.go \
		pipe.go \
		pipe_test.go \
		shift.go \
		posid.go \
		common.go # TODO main.go は escapeForOutput があるからとりあえず



################################################################################
# compiling go sources - for test-unidic
################################################################################

test-unidic: go/src/hydrocul/ichigo-test-unidic/mkdir var/unidic/dict.dat
	GOPATH=$(realpath .)/go ICHIGO_DICTIONARY_PATH=$(realpath .)/var/unidic/dict.dat go test hydrocul/ichigo-test-unidic

var/ichigo-test-unidic.test: go/src/hydrocul/ichigo-test-unidic/mkdir
	GOPATH=$(realpath .)/go go test -c hydrocul/ichigo-test-unidic -o var/ichigo-test-unidic.test

bench-unidic: var/unidic/dict.dat var/ichigo-test-unidic.test
	GOPATH=$(realpath .)/go ICHIGO_DICTIONARY_PATH=$(realpath .)/var/unidic/dict.dat ./var/ichigo-test-unidic.test -test.cpuprofile=var/cpuprofile-unidic.out -test.bench . -test.benchtime 30s -test.benchmem | tee var/bench-unidic.txt
	GOPATH=$(realpath .)/go go tool pprof -list=hydrocul var/ichigo-test-unidic.test var/cpuprofile-unidic.out > var/cpuprofile-unidic-list.txt

go/src/hydrocul/ichigo-test-unidic/mkdir: \
	main.go \
	da.go \
	utf8.go \
	dict.go \
	dict_data.go \
	pipe.go \
	pipe_test.go \
	shift.go \
	posid.go \
	common.go \
	preprocess.pl
	mkdir -p go/src/hydrocul/ichigo-test-unidic
	perl preprocess.pl unidic go/src/hydrocul/ichigo-test-unidic \
		main.go \
		da.go \
		utf8.go \
		dict.go \
		dict_data.go \
		pipe.go \
		pipe_test.go \
		shift.go \
		posid.go \
		common.go # TODO main.go は escapeForOutput があるからとりあえず



################################################################################
# compiling go sources - for main binary ichigo-ipadic
################################################################################

go/bin/ichigo-ipadic: go/src/hydrocul/ichigo-ipadic/mkdir
	GOPATH=$(realpath .)/go go install hydrocul/ichigo-ipadic

go/src/hydrocul/ichigo-ipadic/mkdir: \
	main.go \
	da.go \
	utf8.go \
	dict.go \
	dict_data.go \
	pipe.go \
	shift.go \
	posid.go \
	common.go \
	preprocess.pl
	mkdir -p go/src/hydrocul/ichigo-ipadic
	perl preprocess.pl ipadic go/src/hydrocul/ichigo-ipadic \
		main.go \
		da.go \
		utf8.go \
		dict.go \
		dict_data.go \
		pipe.go \
		shift.go \
		posid.go \
		common.go



################################################################################
# compiling go sources - for main binary ichigo-unidic
################################################################################

go/bin/ichigo-unidic: go/src/hydrocul/ichigo-unidic/mkdir
	GOPATH=$(realpath .)/go go install hydrocul/ichigo-unidic

go/src/hydrocul/ichigo-unidic/mkdir: \
	main.go \
	da.go \
	utf8.go \
	dict.go \
	dict_data.go \
	pipe.go \
	shift.go \
	posid.go \
	common.go \
	preprocess.pl
	mkdir -p go/src/hydrocul/ichigo-unidic
	perl preprocess.pl unidic go/src/hydrocul/ichigo-unidic main.go \
		da.go \
		utf8.go \
		dict.go \
		dict_data.go \
		pipe.go \
		shift.go \
		posid.go \
		common.go



################################################################################

wagahai-ipadic: etc/wagahai/ipadic-expected.txt var/wagahai-ipadic-actual.txt
	git diff --no-index etc/wagahai/ipadic-expected.txt var/wagahai-ipadic-actual.txt && echo OK

var/wagahai-ipadic-actual.txt: etc/wagahai/wagahai.txt ichigo-ipadic etc/wagahai/filter.sh
	cat etc/wagahai/wagahai.txt | ./ichigo-ipadic | sh ./etc/wagahai/filter.sh > var/wagahai-ipadic-actual.txt

wagahai-unidic: etc/wagahai/unidic-expected.txt var/wagahai-unidic-actual.txt
	git diff --no-index etc/wagahai/unidic-expected.txt var/wagahai-unidic-actual.txt && echo OK

var/wagahai-unidic-actual.txt: etc/wagahai/wagahai.txt ichigo-unidic etc/wagahai/filter.sh
	cat etc/wagahai/wagahai.txt | ./ichigo-unidic | sh ./etc/wagahai/filter.sh > var/wagahai-unidic-actual.txt



################################################################################
