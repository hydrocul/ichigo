
# 必要なフィールドに絞り、読み仮名はひらがなに統一

# 0      1        2        3      4     5     6     7     8        9      10   11       12
# 表層形 左文脈ID 右文脈ID コスト 品詞1 品詞2 品詞3 品詞4 活用種類 活用形 原形 ふりがな 発音
# ↓
# 0      1        2        3      4      5    6        7    8
# 表層形 左文脈ID 右文脈ID コスト 品詞名 原形 ふりがな 発音 代表表記

perl -Mutf8 -MEncode -nle '
    @F = split(/\t/, $_);
    $kana = $F[11];
    $pron = $F[12];

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

    $posname = "$F[4]/$F[5]/$F[6]/$F[7]";
    $posname = $1 while $posname =~ /^(.+)\/\*$/;
    $posname = "$posname:$F[8]" if ($F[8] ne "*");
    $posname = "$posname:$F[9]" if ($F[9] ne "*");
    print "$F[0]\t$F[1]\t$F[2]\t$F[3]\t$posname\t$F[10]\t$kana\t$pron\t";
'


