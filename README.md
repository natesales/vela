# Vela

VELA (the Virtual Ethernet Link Application) is a logical layer 2 link in userspace. Unlike other IP encapsulation protocols, VELA circuits are natively dual-stack, support DNS resolution, and can be bidirectionally configured automatically from just 1 machine.



Up to 255 VELA circuits can be created per unique IP address.



#### Bidirectional autoconfiguration



#### DNS Resolution


#### Updates

- Transport preprocessing

  - Compression: theory.stanford.edu/~matias/papers/infcm_cr.pdf
  - Encryption



###### Control Codes

VC (Vela Control) codes are the first byte in the packet. Most packets will be tagged with a VC-NOP byte followed by the IP packet.

| Code | Definition                             | Function |
| ---- | :-------------------------------------------- | - |
| 0    | No Operation (NOP)                            | Do nothing |
| 1    | Session Initialization Request (IREQ)         | Request a new connection |
| 2    | Session Initialization Acknowledgement (IACK) | Accept a new connection |
| 3    | Session Initialization Confirmation (ICON)    | Confirm new connection is established |
| 4   | Session Closure Request (CREQ) | Notify remote that session is closing |
| 5  |                                               |                                       |

VELA Session Handshake

VELA operates on a UDP underlay, but incorporates stateful session management for auto-configuration.



1. router1 sends it's IP address and next available vid to router2
2. If the session doesn't already exist, router2 responds with an awknowldegement message, otherwise it sends a connection ready message.
3. router1 recieves the ACK and marks the session as established. router1 now starts sending keepalive messages

```
router1       -----IREQ---->       router2  // router1 asks router2 to connect
router1       <----IACK-----       router2  // router2 awknowledges the request
router1       -----ICON---->       router2  // router1 responds 
```





Vela stores a session table including source IP address, VELA ID (vid) and remote IP address.

When a packet comes in, the remote IP and vid are used to look up the remote address, and the packet is sent to that remote address.