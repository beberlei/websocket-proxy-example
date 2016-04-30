<?php
$serverParts = explode(':', $_SERVER['HTTP_HOST'], 2);
$serverName = $serverParts[0];
?>
<html>
    <head>
        <title>Websockets</title>
        <style type="text/css">
            body { background-color: blue; color: #fff; }
            h1 { text-align: center; margin-top: 20px; }
        </style>
    </head>

    <body>
        <h1 id="message"></h1>

        <script type="text/javascript" language="javascript">
            var ws = new WebSocket("ws://<?php echo $serverName; ?>:9191/ws")
            ws.onmessage = function (message) {
                var event = JSON.parse(message.data)

                switch (event.Name) {
                    case "Message":
                        document.getElementById("message").innerHTML = event.Payload.text;
                        break;
                    case "Background":
                        document.getElementsByTagName("body")[0].style.backgroundColor = event.Payload.color;
                        break;
                }
            }
        </script>
    </body>
</html>
