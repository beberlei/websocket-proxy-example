<?php

$fp = @stream_socket_client("udp://udp:9292", $errno, $errstr, 1);
fwrite($fp, json_encode([
    'UserId' => (int)$argv[1],
    'Name' => 'Message',
    'Payload' => ['text' => $argv[2]],
]));
fclose($fp);
usleep(10);
