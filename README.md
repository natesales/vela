# VELA

VELA (the Virtual Ethernet Link Application) is a modern userspace pseudowire protocol. Unlike other IP encapsulation protocols, VELA circuits are natively dual-stack, support DNS resolution, and can be bidirectionally configured automatically from just 1 machine.



### Overview

- UDP Underlay
- 29-byte MTU reduction
- Up to 255 circuits per IP address
- Support for Jumbo Frames and nonstandard MTU sizes
- Autoconfiguration
- DNS resolution
- Seamless IPv6 compatibility



### Performance

VELA operates on a UDP underlay and uses a 1-byte circuit identifier, resulting in a 9-byte header. VELA accounts for the standard 20-byte IP header, so the final MTU will be decremented by 29 bytes. VELA has full support for jumbo frames or other nonstandard MTU sizes, and can automatically scale buffer size based on the interface MTU.

VELA offers very minimal overhead. Between 2 Linux hosts running a circuit on 1 gigabit interfaces, we can get very nearly line rate speed with a single iperf3 stream.

```
Connecting to host 172.16.16.1, port 5201
[  5] local 172.16.16.2 port 37302 connected to 172.16.16.1 port 5201
[ ID] Interval           Transfer     Bitrate         Retr  Cwnd
[  5]   0.00-1.00   sec   110 MBytes   921 Mbits/sec   36    932 KBytes       
[  5]   1.00-2.00   sec   110 MBytes   923 Mbits/sec    3    976 KBytes       
[  5]   2.00-3.00   sec   111 MBytes   933 Mbits/sec    8   1.02 MBytes       
[  5]   3.00-4.00   sec   111 MBytes   933 Mbits/sec   13   1.06 MBytes       
[  5]   4.00-5.00   sec   111 MBytes   933 Mbits/sec    2   1.12 MBytes       
[  5]   5.00-6.00   sec   111 MBytes   933 Mbits/sec   11   1.16 MBytes       
[  5]   6.00-7.00   sec   111 MBytes   933 Mbits/sec    3   1.23 MBytes       
[  5]   7.00-8.00   sec   111 MBytes   933 Mbits/sec    5    915 KBytes       
[  5]   8.00-9.00   sec   112 MBytes   944 Mbits/sec    5    984 KBytes       
[  5]   9.00-10.00  sec   110 MBytes   923 Mbits/sec   10   1.01 MBytes       
- - - - - - - - - - - - - - - - - - - - - - - - -
[ ID] Interval           Transfer     Bitrate         Retr
[  5]   0.00-10.00  sec  1.08 GBytes   931 Mbits/sec   96             sender
[  5]   0.00-10.01  sec  1.08 GBytes   929 Mbits/sec                  receiver

iperf Done.
```

