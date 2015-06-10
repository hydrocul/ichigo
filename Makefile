
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

