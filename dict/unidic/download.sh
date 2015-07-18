# unidicの辞書をダウンロード

if [ -e ./var/unidic/download ]; then
    rm -rvf ./var/unidic/download
fi
if [ -e ./var/unidic/download.tmp ]; then
    rm -rvf ./var/unidic/download.tmp
fi
(
    mkdir -p ./var/unidic/download.tmp || exit 1
    cd ./var/unidic/download.tmp
    wget -O unidic-mecab-2.1.2_src.zip "http://osdn.jp/frs/redir.php?m=iij&f=%2Funidic%2F58338%2Funidic-mecab-2.1.2_src.zip" >&2 || exit 1
    unzip unidic-mecab-2.1.2_src.zip || exit 1
    rm unidic-mecab-2.1.2_src.zip
    cp -v unidic-mecab-2.1.2_src/* ./
    cd ../../..
    mv var/unidic/download.tmp var/unidic/download || exit 1
) || exit 1


