
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

cat $tmpdir/mecab.raw | perl -Mutf8 -MEncode -nle '
    @F = split(/\t/, $_);
    @G = split(/,/, $F[1]);
    $surface = $F[0];
    if ($F[1] eq "") {
        $posname = "";
        $base = "";
        $kana = "";
        $pron = "";
    } else {
        $posname = "$G[0],$G[1],$G[2],$G[3]";
        $posname = $1 while $posname =~ /^(.+),\*$/;
        $posname = "$posname:$G[4]" if ($G[4] ne "*");
        $posname = "$posname:$G[5]" if ($G[5] ne "*");
        $base = $G[6];
        $kana = $G[7];
        $pron = $G[8];
    }

    $kana = decode_utf8($kana);
    $lenKana = length($kana);
    $hiragana = "";
    for ($i = 0; $i < $lenKana; $i++) {
      $c = substr($kana, $i, 1);
      $ch = ord($c);
      if ($ch >= 0x30A1 && $ch <= 0x30F6) {
        $hiragana = $hiragana . chr($ch - 0x60);
      } else {
        $hiragana = $hiragana . $c;
      }
    }
    $kana = encode_utf8($hiragana);

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

    print "$surface\t$posname\t$base\t$kana\t$pron";
' > $tmpdir/mecab.txt

cat $tmpdir/ichigo.raw | perl -nle '
    @F = split(/\t/, $_);
    $surface = $F[2];
    $posname = $F[3];
    $base = $F[4];
    $kana = $F[5];
    $pron = $F[6];
    print "$surface\t$posname\t$base\t$kana\t$pron";
' > $tmpdir/ichigo.txt

diff -u $tmpdir/mecab.txt $tmpdir/ichigo.txt
RESULT=$?

rm -rf $tmpdir

exit $RESULT

