
cat dict_data.go | while read line; do

    if echo $line | grep COMPRESSED-DICT-DATA >/dev/null; then
        echo -n '[]uint8("'

        cat $1 |
        perl -e '
            sub conv {
                my ($ch) = @_;
                my $c = ord($ch);
                if ($c == 0x0A) {
                    "\\n";
                } elsif ($c >= 0x20 && $c <= 0x7E && $c != 0x22 && $c != 0x5c) {
                    $ch;
                } else {
                    sprintf("\\x%02x", $c);
                }
            }
            while (<STDIN>) {
                s/(.)/conv($1)/egs;
                print "$_";
            }
        '

        echo '")'
    else
        echo "$line"
    fi

done


