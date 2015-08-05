
# MeCabの辞書フォーマットで、区切り文字をタブに変更して、Unicode正規化をし、単一ファイルを作成

cat $1/*.csv |
perl -Mutf8 -MEncode -MUnicode::Normalize -nle '
    if (/^[^#]/) {
        my $line = decode_utf8($_);
        $line =~ s/,/\t/g;
        $line = Unicode::Normalize::NFKC($line);
        print encode_utf8($line);
    }
'

