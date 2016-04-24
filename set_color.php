<?php

$fp = @stream_socket_client("udp://127.0.0.1:8081", $errno, $errstr, 1);
fwrite($fp, json_encode([
    'UserId' => (int)$argv[1],
    'Name' => 'Background',
    'Payload' => ['color' => $argv[2]],
]));
fclose($fp);
usleep(10);
