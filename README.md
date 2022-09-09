# filecoin-ipfs-data-rebuilder

## Features

`filecoin-ipfs-data-rebuilder` is a data build-and-rebuild tool between the IPFS network and the Filecoin network. It provides the following functions:

 - **Build** (IPFS To Filecoin): Use the [MCS SDK](https://docs.filswan.com/multi-chain-storage/developer-quickstart/sdk) upload the data to the IPFS node, and **build** at leat 5 cold backups to the  Filecoin network.
 - **Rebuild** (Filecoin To IPFS): 
	- Find the storage providers' peerIds by the [Indexer node](https://github.com/filecoin-project/index-provider) 
	 - Get the storage providers IDs from the [Filecoin Network](https://github.com/filecoin-project/lotus/blob/master/api/v0api/full.go)
	 - Retrieve the data from the storage providers using [Lotus](https://github.com/filecoin-project/lotus)
	 - **Rebuild** the data to the IPFS nodes
 - **Auto-Discover** (IPFS): Rebuilder will check automatically if the IPFS node is healthy, the bad IPFS node can be **discovered** in time
 - **Auto-Rebuild**: The bad IPFS nodes will be removed, and if the hot-backup is 0, Rebuilder will auto-rebuild it from Filecoin network. 

## Installation

> :bell:**go 1.17+** is required

 - Install and Compile
```shell
git clone https://github.com/Fogmeta/filecoin-ipfs-data-rebuilder.git
cd filecoin-ipfs-data-rebuilder
make
```

 - Config

	After compiled, you should update **./config.toml**:
```
[server]
RunMode="debug" # debug or release, default debug
HttpPort=8000 # Default `8000`, web api port for extension in future

[database]
User="root" # mysql database username
Password="password" # mysql database password
Host="127.0.0.1:3306" # mysql database uri(ip:port).
Name="rebuild" # mysql database name

[indexer]
Urls = ["https://cid.contact","https://index-finder.kencloud.com","https://index-finder-piknik.kencloud.com","https://index-finder-sxx.kencloud.com"] # filecoin indexer node urls

[lotus]
FullNodeApi = "https://api.node.glif.io" # lotus full node api url.
DownloadDir="/opt/data/" # The data will retrieved to the path
Address="f3xxxxxxxx" # used to propose the retrieve deals


[uploader]
IpfsUrls = ["https://upload.example.com;https://download.example.com","https://upload.example2.com;https://download.example2.com"] # The IPFS nodes that you hope to rebuild to 

```
 - Run
	You can run it by the follwing command:
```shell
nohup ./rebuilder >> rebuilder.log &
```
## License

[Apache](https://github.com/filswan/go-swan-provider/blob/main/LICENSE)

