### This is the first CovertChannel lab

It containts two programs that run separately on different devices.

The first is `./zakladka/zakladka.go` is a realisation of trafic modification program that sends ICMP packets depending on the rule to reproduce time covert channel.
`zakladka.go` has one flag `-trafic` which is false by deafult. When it is false the program will generate it's own trafic to create Covert Channel. In case if you will specify flag `trafic` then the program will wait for ICMP traffic to bufferize it and send to `receiver.go` in with some special timeouts.
The only argument that the program needs is message. Here is the example of `zakladka.go` run:

```
go run zakladka.go -traffic "H0 h0 h0 M3rry Xmas!"
```


`./receiver/receiver.go` is a program that should run on different device. It listens to ICMP traffic and converts it to the message.
The program does not need any arguments. Here is the example of `receiver.go` run:

```
go run receiver.go
```
