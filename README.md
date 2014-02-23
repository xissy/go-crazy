# go-crazy
> a hassle-free WAN acclerator.

** This project is in progress now! Wait for v1.0 please. **


## Why I've built this
Recently I’ve been heard about Golang that has strong weapons for concurrent programming — `goroutine` & `channel`. A goroutine is a hassle-free lightweight thread. And a channel is a way of communicating between goroutines. Using channels sharing variables could be removed that means we don’t need to worry about a race-condition anymore. Sounds great.

I majored in computer science and have 10 years programming experience in diverse languages and environments. I’m familiar with Java concurrent package and Node.js’ event driven model but sometimes I’ve felt that these are not perfect solutions to manage a multi-core processor. So I decided to taste the power of Golang and implemented a complex network application — WAN accelerator. Because, IMHO, one of the best ways to approach concurrent programming is using network sockets.


## WAN accelerator?
Whenever transferring files on the internet we usually use TCP protocol. TCP guarantees to arrive data in order and lossless. But for this, TCP should send a set of packets and receive an acknowledgment, otherwise TCP doesn’t send next packets. This mechanism can lead a low-usability of network bandwidth on a WAN environment of long delayed Round Trip Time(RTT).

I’ve designed a reliable UDP protocol with a new flow control to alternate TCP protocol especially optimized for transferring bulk sized files on WAN. I know WAN accelerator could include duplicated data cache or end-point hardwares but I focused to implement an application layer protocol first. This approach can be found at Aspera, FileCatalyst, Tsunami, UDT, etc.


## Protocol design


## How to use


## Limits


## Roadmap
  * Encryption
  * Compression
  * GUI tools for various platforms
  * Performance tuning
  * etc.

## License
