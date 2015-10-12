
perl -nle '
    @F = split(/\t/, $_);
    print "$F[0]\t$F[1]\t$F[6]\t$F[11]\t$F[12]\t$F[13]\t$F[14]\t$F[15]";
'
