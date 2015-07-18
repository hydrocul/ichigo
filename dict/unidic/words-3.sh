# 必要なフィールドに絞り、読み仮名はひらがなに統一

# 0      1        2        3      4     5     6     7     8        9      10       11       12     13   14   15       16   17         18         19         20
# 表層形 左文脈ID 右文脈ID コスト 品詞1 品詞2 品詞3 品詞4 活用種類 活用形 ふりがな 代表表記 表層形 発音 原型 原型発音 語種 語頭変化型 語頭変化形 語末変化型 語末変化形
# ↓
# 表層形 左文脈ID 右文脈ID コスト 品詞名 原形 読み

# 語種: ※,不明,和,固,外,混,漢,記号


perl -Mutf8 -MEncode -nle '
    @F = split(/\t/, $_);
    $kana = decode_utf8($F[10]);
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
    $hiragana = encode_utf8($hiragana);
    $posname = "$F[4]/$F[5]/$F[6]/$F[7]";
    $posname = $1 while $posname =~ /^(.+)\/\*$/;
    print "$F[0]\t$F[1]\t$F[2]\t$F[3]\t$posname\t$F[14]\t$hiragana";
'


