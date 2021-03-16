# 基于Web Socket的Web实时通讯研究

## 摘要

针对传统的Web实时通信解决方案如轮询、长轮询、flash插件等的局限性，提出将新兴的Web套接字技术应用于Web实时通信领域，介绍了Web套接字技术的特点，分析了Web套接字协议与HTTP协议的区别，提出了一种在客户端和服务器端实现WebSocket的方法，通过实验证明WebSocket可以大大降低网络流量和延迟，并对WebSocket在Web实时通信中的应用前景进行了展望。



**关键字**: Web套接字协议; 全双工;HTTP流;长轮询;延迟



## 引言

在飞速发展的信息时代，互联网已经成为人们生活中不可或缺的一部分。人们对互联网的需求已经从Web1.0时代的信息可访问性转变为Web2.0时代的信息交互，并在越来越多的定价系统、电子商务系统和新闻发布系统中观察到当前的即时交互。

目前，客户端浏览器和服务器之间的通信是基于超文本传输协议（HTTP）的，HTTP是一种基于请求-响应的无状态应用层协议。HTTP客户机启动请求。它建立了一个传输控制协议（TCP）连接。在接收到客户端的请求消息后，服务器将消息作为响应发送回，并终止连接。在这种模式下，服务器不能向客户端发送实时数据。因此，Flash、Comet和Ajax长轮询等技术被用来实现客户机和服务器之间的实时通信。然而，这些技术不能满足实时通信的要求，因为其中一些技术需要在浏览器上安装插件，其中一些技术会给服务器带来沉重的负载。HTML5的出现和WebSocket协议的出现，实现了基于Web系统的实时数据传输，目前认为是解决这一问题的最佳解决方案。

## 传统Web实时通讯的解决方案

轮询、长轮询和HTTP流是过去Web开发人员用来实现浏览器和服务器之间实时通信的主要解决方案。

用自动运行程序代替人工刷新页面的轮询方法是最早应用于浏览器的实时通信解决方案。此解决方案的最大优点是易于实现，并且对客户端和服务器没有额外的要求。然而，这种解决方案存在一些明显的缺点，很难计算出数据更新的频率，因此浏览器无法及时获取最新的数据。另外，在一段时间内没有数据更新的情况下，浏览器的频繁请求会产生不必要的网络流量，给服务器造成不必要的负担。

为了使服务器能够随时与浏览器进行通信，Web开发人员设计了一种新的访问机制被称为长轮询或Comet，通过该机制，服务器可以将来自浏览器的新请求保存一段时间，而不是立即发送响应。如果在此期间发生数据更新，服务器将用新的数据响应浏览器，浏览器将在收到响应时发出另一个请求<sup>[1]</sup>。通过这种机制，浏览器可以及时获取服务器端的最新数据。但是，如果发生大量并发，那么维护这些活动的HTTP连接将极大地消耗服务器内存和计算能力。

开发人员还尝试了“HTTP流”访问机制。它的主要区别是服务器永远不会关闭浏览器发起的连接，服务器会随时使用这个连接发送消息。在这种情况下，由于服务器不会发出连接完成的信号，来自服务器的响应可能会被网络中的防火墙和代理服务器缓冲，从而导致浏览器在接收数据时出现一些错误。

## WebSocket简介

Websocket作为HTML5的一个新特性，被定义为一种使Web页面能够使用WebSocket协议与远程主机进行全双工通信的技术。它引入了Web套接字接口，并定义了一个全双工通信信道，该信道通过Web上的单个套接字进行操作<sup>[2]</sup>。Html5 WebSocket以最小的开销有效地提供到互联网的套接字连接。与Ajax轮询和Comet解决方案相比，它大大减少了网络通信量和延迟，Ajax轮询和Comet解决方案通常用于通过两个HTTP连接传输实时数据来模拟全双工通信。因此，它是构建可扩展的、实时的web通信系统的理想技术。

要使用Html5 Websocket将一个Web客户机与另一个远程端点连接起来，应使用表示要连接的远程端点的有效URL初始化一个新的Web套接字。Websocket将ws://和wss://方案分别定义为websocket和安全websocket连接。在客户机和服务器之间的初始握手期间，将HTTP协议更新为Web套接字协议时，将建立Web套接字连接。

Web套接字连接使用标准HTTP端口（80和443），因此被称为“代理服务器和防火墙友好协议”<sup>[3]</sup>。因此，Html5 websocket不需要安装任何新的硬件。不需要任何中间服务器（代理或反向代理服务器、防火墙、负载平衡路由器等），只要客户机和服务器都支持websocket协议，就可以成功地建立新的websocket连接。

##  比较WebSocket连接和HTTP连接

客户端和服务器之间的通信通常是基于HTTP连接的，HTTP连接要求在客户端的请求和服务器的响应之间附加头，根据HTTP协议的定义，这些头包含一些传输控制信息，如协议类型、协议版本、浏览器类型、传输语言、编码类型、输出等时间，Cookie和会话。在Firebug等软件的帮助下，打开Live HTTP报头，可以清晰地观察到请求和响应的报头。一个请求和响应的头的示例定义如下：

从客户端（浏览器）到服务端:

```http
GET /long-polling HTTP/1.1 
Host: www.kaazing.com 
User-Agent: Mozilla/5.0 (X11; U; Linux x86_64; en-US; rv:1.9) Gecko/2008061017 Firefox/3.0 
Accept: text/html,application/xhtml+xml,application/xml;q = 0.9, */*; q = 0.8 
Accept-Language: en-us,en;q = 0.5 
Accept-Encoding: gzip,deflate 
Accept-Charset: ISO-8859-1,utf-8;q = 0.7,*;q = 0.7 
Keep-Alive: 300 
Connection: keep-alive 
Cache-Control: max-age = 0 
Referer: http://www.example.com/ 
```

从服务端到客户端（浏览器）:

```Http
Date: Tue, 16 Aug 2008 00:00:00 GMT 
Server: Apache/2.2.9 (Unix) 
Content-Type: text/plain
Content-Length: 12

Hello world
```

从以上两个header可以看出，除了数据“Hello World”之外，这些header中的大部分数据在客户端和服务器之间的交互过程中对最终用户来说都是无用的，更不用说Cookie和Session了（在大多数网站中，这两个项目中包含的信息通常多于header中的控制信息）。此外，这些类型的头将包含在每个交互中。因此，如果采用轮询和Comet解决方案，必然会浪费大量的带宽，产生大量的网络流量。此外，构建和分析报头将占用一些用于处理请求和响应的时间，并导致一定程度的延迟。polling和Comet的这些缺点表明，这两种技术在未来必须被其他实时通信技术所取代。让我们转到Web套接字连接。
Web套接字使用HTTP升级机制升级到Web套接字协议<sup>[4]</sup>。websocket的握手机制与HTTP兼容。因此，HTTP服务器可以与Web套接字服务器共享默认的HTTP和HTTPS端口（80和443）。为了建立一个新的websocket连接，在客户端和服务器的初始握手过程中，HTTP协议将升级为websocket协议。一旦建立了连接，Web套接字将基于全双工模式在客户机和服务器之间来回传输。初始握手的标头如下所示：

从客户端（浏览器）到服务端:

```http
GET /text HTTP/1.1 
Upgrade: WebSocket 
Connection: Upgrade
Host: www.websocket.org
```

从服务端到客户端（浏览器）:

```Http
HTTP/1.1 101 WebSocket Protocol Handshake
Upgrade: WebSocket
Connection: Upgrade

Hello world
```

很明显，Web套接字连接头中包含的控制信息远小于HTTP连接头中的控制信息。另一方面，根据websocket协议的规范，Web套接字连接头不允许Cookie和Session，首先，连接成功建立后，客户端可以与服务器自由通信，并且只有两位控制信息附加到终端用户所需的数据上，该数据由UTF-8编码，一位位于起始位置的“\x00”，另一位为位于末端的“\xFF”。这种websocket的定义大大减少了websocket连接中处理头文件所消耗的带宽和时间，从而减少了网络流量和延迟。这就是为什么Web套接字比轮询和Comet更适合于基于Web的实时通信的确切原因。
从安全角度看，websocket协议和HTTP协议都能实现安全传输。Wss和https是它们各自的安全传输协议。因此，Web套接字是网络流量、延时和安全性等方面实时通信的理想技术。

## WebSocket技术的实现

### 客户端websocket的实现

客户端的实现相对简单。以下是W3C工作组<sup>[5]</sup>给出的Web套接字接口定义：

```javascript
[Constructor(in DOMString url, in optional DOMString protocol)]
interface WebSocket {
    readonly attribute DOMString URL;
    // ready state
    const unsigned short CONNECTING = 0;
    const unsigned short OPEN = 1;
    const unsigned short CLOSED = 2;
    readonly attribute unsigned short readyState;
    readonly attribute unsigned long bufferedAmount;
    // networking
        attribute Function onopen;
        attribute Function onmessage;
        attribute Function onclose;
    boolean send(in DOMString data);
    void close();
};
```

根据接口和构造函数的定义，一个新的websocket实例可以由两个参数初始化，一个必要的参数是有效的网络地址，另一个可选的参数是protocol type。在浏览器中，websocket对象由JavaScript操作。一段简单的代码可以创建一个新的Web套接字实例：

```javascript
var myWebSocket = new WebSocket(“ws://www.websocket.org”);
```

前缀“ws”表示Web套接字连接，与此相关的是“wss”，它表示安全的Web套接字连接。初始化前，需要检查客户端浏览器是否支持Web套接字技术：

```javascript
if (“WebSocket” in window)
{var ws = new WebSocket
(“ws://example.com/service”);}
else
{alert(“WebSockets NOT supported here”);}
```

在发送消息之前，需要注册几个函数来处理Web套接字连接过程中的一系列事件，如成功建立连接、接收消息、关闭消息等。

```javascript
myWebSocket.onopen = function(evt){alert(“Connection open ...”); };
myWebSocket.onmessage = function(evt) {alert( “Received Message: “+ evt.data);};
myWebSocket.onclose = function(evt){alert(“Connection closed.”);};
```

要发送消息，只需调用post message方法，然后调用message content作为其默认参数。消息发送后，调用disconnect方法终止连接。

```javascript
myWebSocket.postMessage("Hello Web Socket!");
myWebSocket.disconnect(); 
stockTickerWebSocket.disconnect();
```

### 服务端WebSocket的实现

与客户端的实现相比，服务器端的实现要复杂得多，客户端的大多数操作，如生成头文件、分析头文件、提取有用数据等，都是由浏览器自动完成的，但这些都不是在服务器端实现的，需要开发人员手工完成。服务器端Web套接字的实现主要依赖于套接字的实现，对于C#语言、java语言和C++语言都是常见的。本文采用C#语言实现了服务器端的websocket功能。
首先，应该创建一个侦听器来监视网络中的新请求。

```c#
private Socket serverListener = new Socket(AddressFamily.InterNetwork,SocketType.Stream, protocolType.IP);
```

accept函数负责监听新来的请求，因此，应该将它放入一个循环中，该循环一直运行，以便随时获取客户机请求。

```c#
while (true) 
{ 
Socket sc = serverListener.Accept();
//get a new connection
if (sc != null) { … }//process the request
}
```

同样，当接收到新的连接请求时，需要注册一些服务器函数来处理在通信期间发生的事件，如接收消息、发送消息、关闭连接。

```c#
ci.ReceiveData += new ClientSocketEvent (Ci_ReceiveData);
ci.BroadcastMessage += new BroadcastEvent (ci.SendMessage);
ci.DisConnection += new ClientSocketEvent (Ci_DisConnection)
```

然后调用BeginReceive方法接收客户端请求消息，并尝试与客户端浏览器握手，如果握手成功，则可以启动全双工通信。

```c#
ci.ClientSocket.BeginReceive(ci.receivedDataBuffer,0, ci.receivedDataBuffer.Length, 0,new syncCallback(ci.StartHandshake), ci.ClientSocket.Available);
```

在上面的代码中，StartHandshake方法负责根据客户机请求生成握手信息。该方法从请求头中提取“Sec-Web-Socket-Key1”和“Sec-WebSocket-Key2”两个键的值，并根据这两个值进行Web Socket协议中定义的MD5计算，最后返回结果，MD5结果是在握手过程中保护数据的一种手段<sup>[6]</sup>。如果握手成功，新的连接将被放入连接轮询中，以便下次重用。

```c#
listConnection.Add(ci);
```

这里需要注意的是，在发送消息时，将字符“\x00”放在消息的开头，“\xFF”放在消息的结尾，在阅读消息时，将这两个字符重新移动。另外，在使用之前，消息应该由UTF-8编码或解码。

```c#
public void SendMessage(MessageEntity me)
{
    ClientSocket.Send(new byte[] {0x00});
    ClientSocket. Send(Encoding.UTF8,GetBytes(JsonConvert.SerializeObject(me)));
    ClientSocket.Send(new byte[] { 0xff });
}
```

最后，调用DisConnection方法在通信结束时关闭连接。

## WebSocket性能分析

效率是实时数据传输的关键问题，也是衡量一个协议是否适合实时数据传输的重要标准。对异步传输中的Web套接字性能进行了测试。测试分为两部分，第一部分是基于HTTP请求对phpMyAdmin页面中包含五列三行数据的表进行排序，第二部分是基于websocket通信对单独页面中相同大小的表进行排序。在整个测试过程中，采用Google Chrome 5作为客户端浏览器，Wireshark网络协议分析仪作为监控工具，观察数据包和位流的变化。最后，得到以下结果（图1）：
从上面的数据（表1）可以看出，套接字连接的效率是HTTP连接的十倍。另一方面，参考Peter Lubbers和Frank Greco关于Ajax轮询和Web Socket效率比较的测试<sup>[7]</sup>，可以得出Web Socket在网络流量和延迟方面的性能要比HTTP好得多的结论，特别是在大量并发的情况下。

![图1. HTTP连接和套接字连接在流量和时间上的比较](C:\Users\Peanuts\AppData\Roaming\Typora\typora-user-images\image-20210316200512010.png)

<center>图1. HTTP连接和套接字连接在流量和时间上的比较</center>

<center>表1. 测试中使用的数据表</center>

<table style="text-align: center" rules=rows>
   <tr>
      <td></td>
      <td colspan="2">数据包数</td>
      <td colspan="2">字节数</td>
      <td colspan="2">时间（秒）</td>
   </tr>
   <tr>
      <td></td>
      <td>HTTP</td>
      <td>WebSocket</td>
      <td>HTTP</td>
      <td>WebSocket</td>
      <td>HTTP</td>
      <td>WebSocket</td>
   </tr>
   <tr>
      <td>客户端到服务端</td>
      <td>83</td>
      <td>5</td>
      <td>33,662</td>
      <td>372</td>
      <td></td>
      <td></td>
   </tr>
   <tr>
      <td>服务端到客户端</td>
      <td>77</td>
      <td>8</td>
      <td>45,600</td>
      <td>7456</td>
      <td></td>
      <td></td>
   </tr>
   <tr>
      <td>合计</td>
      <td>160</td>
      <td>13</td>
      <td>79,262</td>
      <td>7828</td>
      <td>~2.5</td>
      <td>~0.25</td>
   </tr>
</table>






## 总结

实时数据传输将是网络信息系统发展的必然趋势。websocket作为下一代Ajax将在Internet上得到广泛应用。目前，最流行的浏览器IE8及其更低版本仍然不支持websocket。然而，Kaazing公司一直在开发一种智能网关，它可以将低版本浏览器中使用的Ajax轮询和Comet转换为websocket即时通信。Web套接字协议和Web套接字API仍在更新中。在不久的将来，websocket可能会成为解决“C10K”问题的完美解决方案。

## 引用

[1]    D. G. Synodinos, “HTML 5 Web Sockets vs. Comet and Ajax,” 2008. http://www.infoq.com/news/2008/12/websockets-vs-com et-ajax 

[2]    Wikipedia, “WebSockets,” 2010. http://en.wikipedia.org/wiki/WebSockets 

[3]    Peter Lubbers, “Pro HTML 5 Programming,” Apress, Victoria, 2010. 

[4]    W3C, “The Web Sockets API,” 2009. http://www.w3.org/TR/2009/WD-websockets-20091222 

[5]    D. Sheiko, “Persistent Full Duplex Client-Server Connec-tion via Web Socket,” 2010. http://dsheiko.com/weblog/persistent-full-duplex-client-ser-ver-connection-via-web-socket

[6]    Makoto, “Living on the Edge of the WebSocket Proto-col,” 2010. http://blog.new-bamboo.co.uk/2010/6/7/living-on-the-edge of-the-websocket-protocol 

[7]    P. Lubbers and F. Greco, “HTML5 Web Sockets: A Quan-tum Leap in Scalability for the Web,” 2010. http://websocket.org/quantum.html, 2010. 