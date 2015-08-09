# ipadicの辞書をダウンロード

if [ -e ./var/download/ipadic ]; then
    #rm -rvf ./var/download/ipadic
    exit 0
fi
if [ -e ./var/download/ipadic.tmp ]; then
    rm -rvf ./var/download/ipadic.tmp
fi
(
    mkdir -p ./var/download/ipadic.tmp || exit 1
    cd ./var/download/ipadic.tmp
    wget "http://mecab.googlecode.com/files/mecab-ipadic-2.7.0-20070801.tar.gz" >&2 || exit 1
    tar xvzf mecab-ipadic-2.7.0-20070801.tar.gz >&2 || exit 1
    rm mecab-ipadic-2.7.0-20070801.tar.gz
    for f in `ls mecab-ipadic-2.7.0-20070801/`; do
        echo "cat mecab-ipadic-2.7.0-20070801/$f | nkf -Ew > $f" >&2
        cat mecab-ipadic-2.7.0-20070801/$f | nkf -Ew > $f
    done
    cd ../../..
    mv var/download/ipadic.tmp var/download/ipadic || exit 1
) || exit 1


