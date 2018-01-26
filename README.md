# MCPingServer - Golang minecraft stand-in server

MCPingServer is a minimalistic Minecraft server which only responds to ping
requests and kicks any clients which attempt to join the server. All facets of
the server are configurable: favicon location, MOTD, player count, player cap,
as well as kick message. All messages are colored and formatted with the
correct formatting (>=1.7 json style, <= 1.6 section color character).

# Installation

MCPingServer is a Golang based program, which means it can (and is!) distributed
in a variety of compiled varieties. There should be a build available for your
operating system/architecture. However, if there is not, you are still able to
build MCPingServer from source (outlined below).

# Building From Source

Building from source requires `go` to be installed on the machine designated
for building MCPingServer. Once `go` is installed, building MCPingServer can
be achieved by running the following set of commands:

    git clone https://github.com/Ichbinjoe/MCPingServer.git
    cd MCPingServer
    go build .

This results in an executable to be created within the current directory.

# Using MCPingServer

MCPingServer is a simple command line executable. The behavior of MCPingServer
is mutated by a set of command line flags that all can be optionally supplied
to MCPingServer.

Flags:

+ `-bindAddr Ip:Port` Ip:Port combination for MCPingServer to bind to (default "0.0.0.0:25565")
+ `-cap players` Reported max player count for the server (default 20)
+ `-players players` Reported current player count for the server (default 0)
+ `-favicon faviconLocation` Location to read the favicon from. Favicon should
be a 64px by 64px png (default "favicon.png")
+ `-kickMsg msg` Message to send to players who attempt to join the server.
Supports `&` color codes (default: "&4This is not a joinable server!")
+ `-motd motd` Message to set the MOTD to. Supports `&` colors as well as `/n`
as the newline seperator (default "A Golang inplace server")
+ `-serverVersionName name` Sets the name of the version of Minecraft to respond
as (default "1.12.2")
+ `-serverVersionNumber number` Sets the protocol version of Minecraft to 
respond as (default 340)


