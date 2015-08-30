use strict;
use warnings;
use File::Path;

my $tmpdir = "/tmp/ichigo-dictionary";
my $sha1 = "##SHA1##";

my $dictionaryPath = "$tmpdir/$sha1/dict.dat";
if ( ! -e $dictionaryPath) {
    mkpath("$tmpdir/$sha1");
    my $command = "tar xz";
    my $READER;
    my $WRITER;
    pipe($READER, $WRITER);
    my $pid = fork;
    if ($pid) {
        close($READER);
        my $blockSize = 4096;
        while () {
            my $buf;
            my $l = read(DATA, $buf, $blockSize);
            last if ($l <= 0);
            print $WRITER $buf;
        }
        close($WRITER);
        wait;
    } elsif (defined $pid) {
        close($WRITER);
        chdir "$tmpdir/$sha1";
        open(STDIN, '<&=', fileno($READER));
        open(STDOUT, '>&=', fileno(STDERR));
        exec($command);
    } else {
        die;
    }
}

$ENV{"ICHIGO_DICTIONARY_PATH"} = "$tmpdir/$sha1/dict.dat";

# TODO help表示のパスが一時ファイルになってしまっている

exec("$tmpdir/$sha1/main", @ARGV);

