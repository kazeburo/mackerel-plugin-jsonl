#!/usr/bin/perl

use JSON;

my $file = $ARGV[0] || "json.log";
my $i = $ARGV[1] || 1;
my @status = (200, 201, 202, 204, 301, 302, 400, 401, 403, 404, 499, 500, 502, 503);
srand(time ^ $$);



open(my $fh, ">>", $file) or die $!;
while ($i>0) {
    my $log = JSON::encode_json({
        "time" => "08/Mar/2017:14:12:40 +0900",
        "status" => $status[int(rand(@status))],
        "reqtime" => 0.030,
        "host" => "10.20.30.40",
        "req" => "GET /example/path HTTP/1.1",
        "method" => "GET",
        "size" => 941,
        "ua" => "Mozilla/5.0 (Linux; Android 4.4.2; SO-01F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.90 Mobile Safari/537.36",
    });
    print $fh "$log\n";
    $i--;
}