
if ! which perl >/dev/null; then
    echo "Not found perl!" >&2
    exit 1
fi

sha1=`cat $1/dict.dat $1/main $1/graph.pl | sha1sum -b | cut -b-40`

echo -n "#!"
which perl

echo

cat boot.pl | sed "s/##SHA1##/$sha1/g"

echo
echo "__DATA__"

cd $1

tar cz dict.dat main graph.pl

