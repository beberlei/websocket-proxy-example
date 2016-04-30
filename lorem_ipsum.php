<?php

$counterUsers = $argv[1];

$words = explode(" ", "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque lacus felis, condimentum vitae enim quis, dictum mollis lorem. Sed non dignissim ligula. Quisque luctus porttitor dui, a dapibus diam gravida non. Nulla sit amet hendrerit ligula. Nulla at est rutrum, hendrerit ipsum sed, eleifend dolor. Sed dolor diam, fermentum sit amet tincidunt vel, cursus sed enim. Ut sagittis tellus dolor, in scelerisque velit porttitor id. Morbi quis volutpat erat, et volutpat enim. Donec ut odio imperdiet, pharetra orci eu, lacinia nisi. Nulla quis ullamcorper turpis. Nam commodo orci nec justo faucibus sagittis. Donec nec ante sollicitudin, finibus mi eu, blandit felis. Cras eget purus sem. Curabitur pretium tempor lacinia.");
$colors = ["green", "red", "blue", "pink", "yellow", "orange", "black", "cyan"];
$count = 0;

while (true) {
    $fp = @stream_socket_client("udp://udp:9292", $errno, $errstr, 1);
    for ($userId = 1; $userId <= $counterUsers; $userId++) {
        $message = $words[$count % count($words)];
        $color = $colors[$count % count($colors)];

        fwrite($fp, json_encode(['UserId' => $userId, 'Name' => 'Message', 'Payload' => ['text' => $message]]));
        fwrite($fp, json_encode(['UserId' => $userId, 'Name' => 'Background', 'Payload' => ['color' => $color]]));
    }
    fclose($fp);
    sleep(1);
    $count++;
}
