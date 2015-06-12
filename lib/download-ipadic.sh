# ipadicの辞書をダウンロード

if [ -e ./var/ipadic ]; then
    rm -rvf ./var/ipadic
fi
(
    mkdir -p ./var/ipadic.tmp || exit 1
    cd ./var/ipadic.tmp
    wget "http://mecab.googlecode.com/files/mecab-ipadic-2.7.0-20070801.tar.gz" >&2 || exit 1
    tar xvzf mecab-ipadic-2.7.0-20070801.tar.gz >&2 || exit 1
    rm mecab-ipadic-2.7.0-20070801.tar.gz
    for f in `ls mecab-ipadic-2.7.0-20070801/`; do
        echo "cat mecab-ipadic-2.7.0-20070801/$f | nkf -Ew > $f" >&2
        cat mecab-ipadic-2.7.0-20070801/$f | nkf -Ew > $f
    done
    cd ../..
    mv var/ipadic.tmp var/ipadic || exit 1
) || exit 1


