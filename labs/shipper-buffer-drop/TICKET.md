# INC-1015: New log lines never reach the aggregator, no errors anywhere

/var/log/app.log keeps getting new lines, but nothing new has landed in
/var/log/shipped-final.log in a while. ship-logs runs clean every time
with no error output. Get new lines actually shipped again.

Get a shell: docker exec -it breaklab-shipper-buffer-drop bash
