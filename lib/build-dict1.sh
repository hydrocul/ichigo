# ipadicのフォーマットで、区切り文字をタブに変更して、Unicode正規化をし、単一ファイルを作成

# 表層形,左文脈ID,右文脈ID,コスト,品詞,品詞細分類1,品詞細分類2,品詞細分類3,活用形,活用型,原形,読み,発音

cat ./var/ipadic/*.csv |
perl -Mutf8 -MEncode -MUnicode::Normalize -nle '
    if (/^[^#]/) {
        my $line = decode_utf8($_);
        $line =~ s/,/\t/g;
        $line = Unicode::Normalize::NFKC($line);
        print encode_utf8($line);
    }
' |
LC_ALL=C sort

