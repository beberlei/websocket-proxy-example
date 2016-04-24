<?php

$counterUsers = $argv[1];
$winner = rand(1, $counterUsers);

$fp = @stream_socket_client("udp://127.0.0.1:8081", $errno, $errstr, 1);
for ($userId = 1; $userId <= $counterUsers; $userId++) {
    $message = ($userId == $winner) ? "YOU WON" : "YOU LOST";
    $color = ($userId == $winner) ? "green" : "red";

    fwrite($fp, json_encode(['UserId' => $userId, 'Name' => 'Message', 'Payload' => ['text' => $message]]));
    fwrite($fp, json_encode(['UserId' => $userId, 'Name' => 'Background', 'Payload' => ['color' => $color]]));
}
fclose($fp);
usleep(10);
