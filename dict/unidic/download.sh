# unidicの辞書をダウンロード

if [ -e ./var/download/unidic ]; then
    rm -rvf ./var/download/unidic
fi
if [ -e ./var/download/unidic.tmp ]; then
    rm -rvf ./var/download/unidic.tmp
fi
(
    mkdir -p ./var/download/unidic.tmp || exit 1
    cd ./var/download/unidic.tmp
    wget -O unidic-mecab-2.1.2_src.zip "http://osdn.jp/frs/redir.php?m=iij&f=%2Funidic%2F58338%2Funidic-mecab-2.1.2_src.zip" >&2 || exit 1
    unzip unidic-mecab-2.1.2_src.zip || exit 1
    rm unidic-mecab-2.1.2_src.zip
    cp -v unidic-mecab-2.1.2_src/* ./
    cd ../../..
    mv var/download/unidic.tmp var/download/unidic || exit 1
) || exit 1


