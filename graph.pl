use strict;
use warnings;
use utf8;

my %nodes = ();
my $lastId = 0;

while (my $line = <STDIN>) {
    next unless $line =~ /^g\t/;
    $line =~ s/\A(.+)\s*\z/$1/g;
    my @cols = split(/\t/, $line);
    my $nexts = [];
    my $id = $cols[1];
    my $prev = $cols[2];
    my $text = $cols[3];
    my $left = $cols[4];
    my $right = $cols[5];
    my $wordCost = $cols[6];
    my $totalCost = $cols[7];
    my $flag = "";
    $nodes{$id} = [$nexts, $id, $prev, $text, $left, $right, $wordCost, $totalCost, $flag];
    $lastId = $id;
}

foreach my $id (sort { $a <=> $b } keys %nodes) {
    my $prev = $nodes{$id}->[2];
    if ($prev >= 0) {
        push(@{$nodes{$prev}->[0]}, $id);
    }
}

my $currId = $lastId;
while ($currId >= 0) {
    $nodes{$currId}->[8] = 1;
    $currId = $nodes{$currId}->[2];
}

my @lineIds = (0);

sub replaceLineIds {
    my ($id, $nexts) = @_;
    my @nexts = @$nexts;
    my @newLineIds = ();
    for (my $i = 0; $i < @lineIds; $i++) {
        unless ($lineIds[$i] == $id) {
            push(@newLineIds, $lineIds[$i]);
            next;
        }

        my $pos = $i;

        push(@newLineIds, @nexts);
        for ($i++; $i < @lineIds; $i++) {
            push(@newLineIds, $lineIds[$i]);
        }

        @lineIds = @newLineIds;

        return $pos;
    }
    die;
}

sub outputLines {
    my ($pos, $prevLineCount, $nextCount) = @_;
    if ($nextCount == 1) {
        return;
    } elsif ($nextCount == 0) {
        my $i = 0;
        for (; $i < $pos; $i++) {
            print "| ";
        }
        print " ";
        $i++;
        for (; $i < $prevLineCount; $i++) {
            print "/ ";
        }

        # | * | |
        # |  / /
        # | | |

        print "\n";
        return;
    } elsif ($nextCount == 2) {
        my $i = 0;
        for (; $i < $pos; $i++) {
            print "| ";
        }
        print "|\\ ";
        $i++;
        for (; $i < $prevLineCount; $i++) {
            print "\\ ";
        }

        # |\ \
        # | | |

        print "\n";
        return;
    } else {
        for (my $j = 1; $j < $nextCount; $j++) {
            my $i = 0;
            for (; $i < $pos; $i++) {
                print "| ";
            }
            print "|\\ ";
            $i++;
            for (; $i < $prevLineCount; $i++) {
                print "\\ ";
            }
            print "\n";

            $i = 0;
            for (; $i < $pos; $i++) {
                print "| ";
            }
            print "| \\ ";
            $i++;
            for (; $i < $prevLineCount; $i++) {
                print "\\ ";
            }
            print "\n";

            # |\ \
            # | \ \

            $prevLineCount++;
        }
    }
}

foreach my $id (sort { $a <=> $b } keys %nodes) {
    my @cols = @{$nodes{$id}};
    my $nexts = $cols[0];
    my $id = $cols[1];
    my $prev = $cols[2];
    my $text = $cols[3];
    my $left = $cols[4];
    my $right = $cols[5];
    my $wordCost = $cols[6];
    my $totalCost = $cols[7];
    my $flag = $cols[8];

    my $lineCount = @lineIds;
    my $pos = replaceLineIds($id, $nexts);
    for (my $i = 0; $i < $lineCount; $i++) {
        if ($i == $pos) {
            if ($flag) {
                print "% ";
            } elsif (@$nexts == 0) {
                print "x ";
            } else {
                print "* ";
            }
        } else {
            print "| ";
        }
    }

    my $nextsStr = join(",", @$nexts);
    print "$text $left $right $wordCost $totalCost\n";

    outputLines($pos, $lineCount, scalar @$nexts);
}

