
# テキスト一覧を作成

# 0      1        2        3      4      5    6        7    8
# 表層形 左文脈ID 右文脈ID コスト 品詞名 原形 ふりがな 発音 代表表記

(
    cat $1 | perl -nle '
        @F = split(/\t/, $_);
        print $F[0];
        print $F[4];
        print $F[5];
        print $F[6];
        print $F[7];
        print $F[8];
    '
    cat $2
    cat $3
) | LC_ALL=C sort | LC_ALL=C uniq


