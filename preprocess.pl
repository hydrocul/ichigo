use strict;
use warnings;

my @definitions = split(/,/, $ARGV[0]);

my $dstDir = $ARGV[1];

my @sources = @ARGV[2..$#ARGV];

sub condMatch {
    my @conds = @_;
    foreach my $cond (@conds) {
        if (grep {$_ eq $cond} @definitions) {
            return 1;
        }
    }
    return '';
}

sub preprocess {
    my ($source) = @_;

    my $dst = "$dstDir/$source";
    if ($source eq "build_main.go") {
        $dst = "$dstDir/main.go";
    }

    if (-e $dst && (stat($source))[9] < (stat($dst))[9]) {
        return '' ;
    }

    print "preprocess: $dst\n";

    open(IN, '<', $source) or die;
    open(OUT, '>', $dst) or die;

    my @ifStack = ();

    while (my $line = <IN>) {
        $line =~ s/^(.*)\s*$/$1/;
        if ($line =~ /\A\s*#ifdef\s+([\sa-zA-Z0-9]+)\z/) {
            my @conds = split(/\s+/, $1);
            if (@ifStack && !$ifStack[0]) {
                unshift(@ifStack, '');
            } elsif (condMatch(@conds)) {
                unshift(@ifStack, 1);
            } else {
                unshift(@ifStack, '');
            }
            print OUT "\n";
        } elsif ($line =~ /\A\s*#endif\s*\z/) {
            die unless (@ifStack);
            shift(@ifStack);
            print OUT "\n";
        } else {
            if (@ifStack && !$ifStack[0]) {
                print OUT "\n";
            } else {
                print OUT "$line\n";
            }
        }
    }

    close(OUT);
    close(IN);

    return 1
}

my $touchFlag = '';
foreach my $source (@sources) {
    if (preprocess($source)) {
        $touchFlag = 1;
    }
}
if ($touchFlag) {
    open(TOUCH, '>', "$dstDir/mkdir");
    close(TOUCH);
}

