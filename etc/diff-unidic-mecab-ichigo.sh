
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

cat $tmpdir/mecab.raw | perl -Mutf8 -MEncode -nle '
    @F = split(/\t/, $_);

    $surface = $F[0];
    if (@F == 1) {
        $posname = "";
        $lemma = "";
        $pron = "";
    } else {
        $posname = $F[4];
        $posname = "$posname:$F[5]" if ($F[5] ne "");
        $posname = "$posname:$F[6]" if ($F[6] ne "");
        $posname =~ s/-/,/g;
        $lemma = $F[3];
        $pron = $F[1];
    }

    $pron = decode_utf8($pron);
    $lenPron = length($pron);
    $hiragana = "";
    for ($i = 0; $i < $lenPron; $i++) {
      $c = substr($pron, $i, 1);
      $ch = ord($c);
      if ($ch >= 0x30A1 && $ch <= 0x30F6) {
        $hiragana = $hiragana . chr($ch - 0x60);
      } else {
        $hiragana = $hiragana . $c;
      }
    }
    $pron = encode_utf8($hiragana);

    print "$surface\t$posname\t$lemma\t$pron";
' > $tmpdir/mecab.txt

cat $tmpdir/ichigo.raw | perl -nle '
    @F = split(/\t/, $_);
    $surface = $F[2];
    $posname = $F[3];
    $lemma = $F[7];
    $pron = $F[6];
    print "$surface\t$posname\t$lemma\t$pron";
' > $tmpdir/ichigo.txt

diff -u $tmpdir/mecab.txt $tmpdir/ichigo.txt
RESULT=$?

rm -rf $tmpdir

exit $RESULT

