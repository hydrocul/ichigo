
make ichigo-ipadic || exit $?

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

cat $INPUT | ./ichigo-ipadic > $tmpdir/ichigo.raw || exit $?
cat $INPUT | ./etc/mecab-ipadic > $tmpdir/mecab.raw || exit $?

cat $tmpdir/mecab.raw | perl -nle '
    @F = split(/\t/, $_);
    print "$F[6]\t$F[11]\t$F[12]\t$F[13]\t$F[14]";
' > $tmpdir/mecab.txt

cat $tmpdir/ichigo.raw | perl -nle '
    @F = split(/\t/, $_);
    $F[11] =~ s/\((.+);.*\)$/$1/g;
    $F[11] =~ s/\((.+)\)$/$1/g;
    print "$F[6]\t$F[11]\t$F[12]\t$F[13]\t$F[14]";
' > $tmpdir/ichigo.txt

diff -u $tmpdir/mecab.txt $tmpdir/ichigo.txt
RESULT=$?

rm -rf $tmpdir

exit $RESULT

