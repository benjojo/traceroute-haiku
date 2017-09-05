Traceroute Haiku’s
===

Sometimes I like to think that I do "serious" blog posts like ["The strange case of ICMP Type 69 on Linux"](https://blog.benjojo.co.uk/post/linux-icmp-type-69) or ["Anycast possibly done better"](https://blog.benjojo.co.uk/post/ipv6-backing-anycast-possibly-better), However I also do a lot of stupid ones like ["IP over AX.25 over 802.11 with ESP8266"](https://blog.benjojo.co.uk/post/AX25-over-wifi-with-ESP8266), ["I may be the only evil (bit) user on the internet"](https://blog.benjojo.co.uk/post/evil-bit-RFC3514-real-world-usage) or ["TOTP SSH port fluxing"](https://blog.benjojo.co.uk/post/ssh-port-fluxing-with-totp)

This one definitely falls into the latter category.

Traceroutes are a common network debugging diagnostic tool. They list (should) every router that your packet travels through to get to its final destination. Here is what it looks like from my flat to my website:

```
# traceroute6 benjojo.co.uk -q 1
traceroute to benjojo.co.uk (2400:cb00:2048:1::6814:6110), 30 hops max, 80 byte packets
 1  switch0.home-edge.bone.benjojo.co.uk (2a07:1500:4663::2)  0.245 ms
 2  home-ipv6.choopa-lhr.bone.benjojo.co.uk (2a07:1500:1111::2)  1.358 ms
 3  *
 4  2001:19f0:7400:8000::1 (2001:19f0:7400:8000::1)  1700.039 ms
 5  ldn-b3-link.telia.net (2001:2000:3080:dc1::1)  1.744 ms
 6  ldn-b5-v6.telia.net (2001:2000:3018:b::1)  1.968 ms
 7  cloudflare-ic-306325-ldn-b3.c.telia.net (2001:2000:3080:a4b::2)  2.307 ms
 8  2400:cb00:21:1024::a29e:99be (2400:cb00:21:1024::a29e:99be)  2.231 ms
```

Traceroutes work using the `Time To Live` ( for IPv4 ) or `Hop Limit` ( for IPv6 )

![IPv6 packet diagram](/blog-images/image1.png)

The idea of this section of IP packets is to stop packets from infinitely going in circles in the case of a fault in a network, it works by for every router a packet jumps though, this number is decreased. If the number reaches zero then the packet is dropped.

However to let the other side know that a packet was lost in the manner, the router is suppose to return the packet inside another packet to notify it:

![IPv6 ICMP packet in wireshark](/blog-images/image6.png)

The idea of traceroute is to purposely set these hop limit very low and incrementally increase it upwards to discover all routers in the path between you and the destination:

{{{GIF???}}}

In addition, traceroute tools helpfully lookup the "reverse DNS" of the IP address to find out more information about the router. However placing reverse DNS on these IP addresses is entirely optional but most operators do to help their clients debug things.

However people have used this common function of traceroute to build fun addresses to trace that often spell out funny things, one example being `bad.horse`:

```
ben@metropolis:~$ traceroute bad.horse --resolve-hostnames -q 1 -f 20
traceroute to bad.horse (162.252.205.157), 64 hops max
  1   162.252.205.3 (t01.nycmc1.ny.us.sn11.net)  91.633ms
  2   162.252.205.130 (bad.horse)  92.802ms
  3   162.252.205.131 (bad.horse)  97.476ms
  4   162.252.205.132 (bad.horse)  104.194ms
  5   162.252.205.133 (bad.horse)  107.089ms
  6   162.252.205.134 (he.rides.across.the.nation)  111.847ms
  7   162.252.205.135 (the.thoroughbred.of.sin)  117.324ms
  8   162.252.205.136 (he.got.the.application)  121.631ms
  9   162.252.205.137 (that.you.just.sent.in)  129.549ms
 10   162.252.205.138 (it.needs.evaluation)  132.131ms
 11   162.252.205.139 (so.let.the.games.begin)  138.446ms
 12   162.252.205.140 (a.heinous.crime)  145.218ms
 13   162.252.205.141 (a.show.of.force)  146.721ms
 14   162.252.205.142 (a.murder.would.be.nice.of.course)  151.821ms
 15   162.252.205.143 (bad.horse)  156.510ms
 16   162.252.205.144 (bad.horse)  161.704ms
 17   162.252.205.145 (bad.horse)  166.850ms
 18   162.252.205.146 (he-s.bad)  171.846ms
 19   162.252.205.147 (the.evil.league.of.evil)  180.018ms
 20   162.252.205.148 (is.watching.so.beware)  187.773ms
 21   162.252.205.149 (the.grade.that.you.receive)  188.834ms
 22   162.252.205.150 (will.be.your.last.we.swear)  191.852ms
 23   162.252.205.151 (so.make.the.bad.horse.gleeful)  196.578ms
 24   162.252.205.152 (or.he-ll.make.you.his.mare)  202.037ms
 25   162.252.205.153 (o_o)  207.124ms
 26   162.252.205.154 (you-re.saddled.up)  211.627ms
 27   162.252.205.155 (there-s.no.recourse)  219.691ms
 28   162.252.205.156 (it-s.hi-ho.silver)  222.107ms
 29   162.252.205.157 (signed.bad.horse)  225.336ms
```

Or Louis Poinsignon (https://www.mygb.eu), who made a version of his CV/Resume for traceroute!

```
# mtr -rwc 1  cv6.poinsignon.org
Start: Mon Sep  4 14:27:22 2017
HOST:                                              Loss%   Snt   Last   Avg  Best  Wrst StDev
  1.|-- switch0.home-edge.bone.benjojo.co.uk          0.0%     1    0.4   0.4   0.4   0.4   0.0
  2.|-- Benjojo-2.tunnel.tserv17.lon1.ipv6.he.net     0.0%     1    1.4   1.4   1.4   1.4   0.0
  3.|-- 10ge3-3.core1.lon2.he.net                     0.0%     1    8.4   8.4   8.4   8.4   0.0
  4.|-- 100ge6-2.core1.ams1.he.net                    0.0%     1   13.4  13.4  13.4  13.4   0.0
  5.|-- amsix.poneytelecom.eu                         0.0%     1   14.6  14.6  14.6  14.6   0.0
  6.|-- 2001:bc8:0:1::131                             0.0%     1   14.7  14.7  14.7  14.7   0.0
  7.|-- 2001:bc8:400:1::32                            0.0%     1   14.8  14.8  14.8  14.8   0.0
  8.|-- hello                                         0.0%     1   14.4  14.4  14.4  14.4   0.0
  9.|-- My.name.is.Louis.Poinsignon                   0.0%     1   14.6  14.6  14.6  14.6   0.0
 10.|-- I.am.a.network.and.systems.Engineer           0.0%     1   14.3  14.3  14.3  14.3   0.0
 11.|-- This.is.my.resume.over.traceroute             0.0%     1   14.7  14.7  14.7  14.7   0.0
 12.|-- o---Experience---o                            0.0%     1   14.5  14.5  14.5  14.5   0.0
 13.|-- 2017.Cloudflare.NetworkEngineer.London        0.0%     1   16.4  16.4  16.4  16.4   0.0
 14.|-- 2016.Cloudflare.NetworkEngineer.Intern.SF     0.0%     1   14.4  14.4  14.4  14.4   0.0
 15.|-- 2015.CEA.SoftwareEngineer.Intern.France       0.0%     1   14.4  14.4  14.4  14.4   0.0
 16.|-- 2014.Android.dev.Remote                       0.0%     1   14.6  14.6  14.6  14.6   0.0
 17.|-- o---Education---o                             0.0%     1   14.5  14.5  14.5  14.5   0.0
 18.|-- 2015-2016.DrexelUni.Exchange.CE.Philadelphia  0.0%     1   14.3  14.3  14.3  14.3   0.0
 19.|-- 2011-2016.UTT.Master.CE.France                0.0%     1   14.6  14.6  14.6  14.6   0.0
 20.|-- o---Skills---o                                0.0%     1   14.4  14.4  14.4  14.4   0.0
 21.|-- C.Java.Python.Maths                           0.0%     1   14.4  14.4  14.4  14.4   0.0
 22.|-- Net.Linux.Archicture                          0.0%     1   14.4  14.4  14.4  14.4   0.0
 23.|-- Statistics.Maths.Design.Photoshop             0.0%     1   14.2  14.2  14.2  14.2   0.0
 24.|-- o---Various---o                               0.0%     1   14.4  14.4  14.4  14.4   0.0
 25.|-- Swimming.and.karate                           0.0%     1   14.3  14.3  14.3  14.3   0.0
 26.|-- Piano                                         0.0%     1   14.4  14.4  14.4  14.4   0.0
 27.|-- o---Contact---o                               0.0%     1   14.4  14.4  14.4  14.4   0.0
 28.|-- mail.jobs.at.poinsignon.org                   0.0%     1   14.7  14.7  14.7  14.7   0.0
 29.|-- cv6.poinsignon.org                            0.0%     1   14.6  14.6  14.6  14.6   0.0
```

In my head I can think of two ways this can be done, either chain lot of fake interfaces inside a single system and use that as "router hops", or I can augment the whole thing with user space networking to generate the fake expiry messages to make it seem like there are routers in the path that are not there.

![stand back i'm going to do user space networking or science (xkcd)](/blog-images/image2.png)

To do this, I am going to use a built in system inside linux called TUN/TAP. This is the system that is used to make VPN’s work. The idea is that it can make a "network port" on your computer, that instead of going to a physical section of hardware, it instead goes to a program that handles the packet.

The idea is, to write a simple network adapter that would 4 fake router hops for any packet destined for an address ending in 4.

With these fake 4 router hops, I can fit 4 small sentences in the reverse DNS entries of the IP addresses. I opted to pick to show a [Haiku](https://en.wikipedia.org/wiki/Haiku) since they are small, typically 3 sentences:

>whitecaps on the bay
>
>the overhead cries
>
>of migrating birds
>
>Polona Oblak

I could turn this into

>Whitecaps.on.the.bay
>
>The.overhead.cries
>
>Of.migrating.birds
>
>author.Polona.Oblak

I then wrote the networking config to route a /48 of IPv6 space to a raspberry pi on my shelf, and got to writing the TUN adapter, the idea was to do something like this:

![networking stack diagram](/blog-images/image5.png)

Thankfully the process to generate TTL/Hop limit expired packets is fairly easy, and you can find the code to do this here: https://github.com/benjojo/traceroute-haiku/blob/master/haiku-tun/main.go

All wrapped up in a basic systemd service and we are good to go:

```
[Service]
Type=simple
ExecStart=/usr/bin/haiku-tun
ExecStartPost=/bin/sleep 2
ExecStartPost=/sbin/ifconfig haiku0 up
ExecStartPost=/bin/ip -6 addr add 2a07:1500:c::1 dev haiku0
ExecStartPost=/bin/ip -6 route add 2a07:1500:c::/64 dev haiku0
ExecStartPost=/sbin/sysctl -w net.ipv6.conf.all.forwarding=1

Restart=always
StandardOutput=syslog
StandardError=syslog
SyslogIdentifier=haiku
User=root

[Install]
WantedBy=multi-user.target
```

I then scraped all of [www.dailyhaiku.org](http://www.dailyhaiku.org) ( sorry ) using Lynx and a bash loop, and then wrote [a small go program](https://github.com/benjojo/traceroute-haiku/blob/master/haikus/generatezonefile.go) to generate the required BIND zone file entries:

```
ben@metropolis:~/traceroute-haiku/haikus$ ./haikus | more
2017/09/05 23:59:04 # We found 3078 haikus
0.0.0.0.0.2.0.0.0.0.0.0.0.0.0.0.0.0.0.0.c.0.0.0.0.0.5.1.7.0.a.2.ip6.arpa.        10        IN        PTR        haiku-trace.x.benjojo.co.uk.
1.0.0.0.0.2.0.0.0.0.0.0.0.0.0.0.0.0.0.0.c.0.0.0.0.0.5.1.7.0.a.2.ip6.arpa.        10        IN        PTR        balmy.breeze.
2.0.0.0.0.2.0.0.0.0.0.0.0.0.0.0.0.0.0.0.c.0.0.0.0.0.5.1.7.0.a.2.ip6.arpa.        10        IN        PTR        swarming.bees.circle.
3.0.0.0.0.2.0.0.0.0.0.0.0.0.0.0.0.0.0.0.c.0.0.0.0.0.5.1.7.0.a.2.ip6.arpa.        10        IN        PTR        the.river.bank.
4.0.0.0.0.2.0.0.0.0.0.0.0.0.0.0.0.0.0.0.c.0.0.0.0.0.5.1.7.0.a.2.ip6.arpa.        10        IN        PTR        author.olona.blak.
```

After that, [some basic bash script](https://github.com/benjojo/traceroute-haiku/blob/master/dns/x-template.sh) to rotate every minute the addresses that `haiku-trace.x.benjojo.co.uk` resolve to, to ensure that you have fresh ones every time, and we are away!

![haiku trace on windows](/blog-images/image4.png)

![haiku trace with linux and mtr](/blog-images/image3.gif)

All served from this precariously hanging pi in my living room :)