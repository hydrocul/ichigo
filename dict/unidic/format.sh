
# 必要なフィールドに絞る

# 0      1        2        3      4     5     6     7     8        9      10                   11           12     13   14   15       16   17         18         19         20
# 表層形 左文脈ID 右文脈ID コスト 品詞1 品詞2 品詞3 品詞4 活用種類 活用形 原型代表表記ふりがな 原型代表表記 表層形 発音 原型 原型発音 語種 語頭変化型 語頭変化形 語末変化型 語末変化形
# ↓
# 0      1        2        3      4      5    6        7    8
# 表層形 左文脈ID 右文脈ID コスト 品詞名 原形 ふりがな 発音 代表表記

# 語種: ※,不明,和,固,外,混,漢,記号


perl -Mutf8 -MEncode -nle '
    @F = split(/\t/, $_);
    $posname = "$F[4]/$F[5]/$F[6]/$F[7]";
    $posname = $1 while $posname =~ /^(.+)\/\*$/;
    print "$F[0]\t$F[1]\t$F[2]\t$F[3]\t$posname\t$F[14]\t\t$F[13]\t$F[11]";
'

