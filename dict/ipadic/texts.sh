
# テキスト一覧を作成

# 0      1        2        3      4      5    6
# 表層形 左文脈ID 右文脈ID コスト 品詞名 原形 ふりがな

perl -nle '
    @F = split(/\t/, $_);
    print $F[0];
    print $F[4];
    print $F[5];
    print $F[6];
' | LC_ALL=C sort | LC_ALL=C uniq


