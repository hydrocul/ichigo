# テキスト一覧を作成

# 表層形,左文脈ID,右文脈ID,コスト,品詞,品詞細分類1,品詞細分類2,品詞細分類3,活用形,活用型,原形,読み,発音
# ↓
# 表層形 左文脈ID 右文脈ID コスト 原形 読み

cat var/dict2.txt |
perl -nle '
    @F = split(/\t/, $_);
    print $F[0];
    print $F[4];
    print $F[5];
' | LC_ALL=C sort | LC_ALL=C uniq


