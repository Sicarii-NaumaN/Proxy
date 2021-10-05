# Burp

## How to build and run 

    $ go mod tidy 
    $ go run . |--https|

<b>Notice:</b> https scheme doesn't work correctly (<b>use</b> go run .)

    [:]:8080 - proxy address
    [:]:8081 - repeater address
    [:]:8082 - param-miner address

## Features and examples

To test functionality try:

1) Proxy
    
    
    curl -x http://127.0.0.1:8080 http://mail.ru 

Output:
        
        <html>
        <head><title>301 Moved Permanently</title></head>
        <body bgcolor="white">
        <center><h1>301 Moved Permanently</h1></center>
        <hr><center>nginx/1.14.1</center>
        </body>
        </html>

In directory proxy/history you can see history of sent requests

    proxy/history/last_request_mail.ru.txt:

    GET / HTTP/1.1
    Host: mail.ru
    User-Agent: curl/7.64.1
    Accept: */*

2) Repeater
 
You can send any request to repeater from history via query parameter:

    curl -i -X PUT http://127.0.0.1:8081\?request=proxy/history/last_request_mail.ru.txt

    proxy/repeater.txt:

    GET / HTTP/1.1
    Host: mail.ru
    User-Agent: curl/7.64.1
    Accept: */*

To repeat request from repeater run:

    curl -i -X POST http://127.0.0.1:8081

3) Param-miner

Param-miner takes target from request in repeater, to run it:

    curl -i http://127.0.0.1:8082

Example of output:

    tatus: 301 ----- length: 178 ----- param: { val }





