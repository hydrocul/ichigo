# 必要なフィールドに絞ってタブ区切りに変換。読み仮名はひらがなに統一

# 表層形,左文脈ID,右文脈ID,コスト,品詞,品詞細分類1,品詞細分類2,品詞細分類3,活用形,活用型,原形,読み,発音
# ↓
# 表層形 左文脈ID 右文脈ID コスト 品詞名 原形 読み

cat var/dict1.txt |
perl -Mutf8 -MEncode -nle '
    @F = split(/\t/, $_);
    $kana = decode_utf8($F[11]);
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
    print "$F[0]\t$F[1]\t$F[2]\t$F[3]\t-\t$F[10]\t$hiragana";
'


