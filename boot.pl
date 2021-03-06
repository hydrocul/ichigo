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

# コマンドライン引数の解釈
# TODO まだ厳密ではない
my $graphFlag = '';
my $i = 0;
while ($i < @ARGV) {
    my $arg = $ARGV[$i];
    if ($arg eq "--graph") {
        $graphFlag = 1;
    }
    $i++;
}

# TODO help表示のパスが一時ファイルになってしまっている
# TODO help表示だったらグラフ表示はしない

if ($graphFlag) {

    my $CHILD_READER;
    my $PARENT_WRITER;
    pipe($CHILD_READER, $PARENT_WRITER);

    my $pid1 = fork;
    if ($pid1) {
        # 親プロセス
        close $CHILD_READER;
        open(STDOUT, '>&=', fileno($PARENT_WRITER));
        exec("$tmpdir/$sha1/main", @ARGV);
    } else {
        # 子プロセス
        die unless defined $pid1;
        close $PARENT_WRITER;
        open(STDIN, '<&=', fileno($CHILD_READER));
        exec('perl', "$tmpdir/$sha1/graph.pl");
    }

} else {
    exec("$tmpdir/$sha1/main", @ARGV);
}

