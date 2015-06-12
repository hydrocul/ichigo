
cat ./var/ipadic/matrix.def |
perl -nle '/^(\d+)\s+(\d+)\s+(-?\d+)$/ and print "$1\t$2\t$3"'

