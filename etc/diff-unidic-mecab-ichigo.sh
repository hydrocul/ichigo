
make ichigo-unidic || exit $?

tmpdir=`mktemp -d`

INPUT=

while [ -n "$1" ]; do
    if [ -z "$INPUT" ]; then
        INPUT=$1
    fi
    shift
done

if [ -z "$INPUT" ]; then
    cat > $tmpdir/input.txt
    INPUT=$tmpdir/input.txt
fi

cat $INPUT | ./ichigo-unidic > $tmpdir/ichigo.raw || exit $?
cat $INPUT | ./etc/mecab-unidic > $tmpdir/mecab.raw || exit $?

cat $tmpdir/mecab.raw | perl -nle '
    @F = split(/\t/, $_);
    print "$F[6]\t$F[11]\t$F[14]\t$F[15]";
' > $tmpdir/mecab.txt

cat $tmpdir/ichigo.raw | perl -nle '
    @F = split(/\t/, $_);
    $F[11] =~ s/\((.+);.*\)$/$1/g;
    $F[11] =~ s/\((.+)\)$/$1/g;
    print "$F[6]\t$F[11]\t$F[14]\t$F[15]";
' > $tmpdir/ichigo.txt

diff -u $tmpdir/mecab.txt $tmpdir/ichigo.txt
RESULT=$?

rm -rf $tmpdir

exit $RESULT

