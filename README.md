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