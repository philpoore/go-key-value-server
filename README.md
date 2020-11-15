# Go Key Value Server

Very simple, very basic, key value tcp server written in go.

Server starts on port 20000.


# Protocall
```
# Set value
-> SET [KEY] [VALUE] 
<- OK

# Get value
-> GET [KEY]
<- [VALUE]

# Quit
-> QUIT
```

# Examples
```
$ nc 127.0.0.1 20000
SET foo 123
OK
GET foo
123
SET foo 321
OK
GET foo
321
```