# FogMeta-Data-Rebuilder

FogMeta Data Rebuilder (Replication and Repair) is a guaranteed storage service based on the FEVM contract. It guarantees N replicas stored on filecoin network, 50% of the storage fund locked in FEVM contract will be used for the initial storage copies and the remaining 50% will be used for future replication when the data replicas loss occurs within the term. 

The reserved fee percentile is based on the current average storage provider failure rate. Assuming the current failure rate per SP is 30% in the term, with 99.99% SLA, we need to maintain 8 replicas all the time. When a replicas loss occurs, a 5% fee will be used for the fee of the replica deal, and unused funds in the contract will be refunded to the user after the term expired. The project won the Data Dao hackathon in 2022. The v1 has provided the following functions: Build (IPFS To Filecoin): upload the data to the IPFS gateway and keep 8 replicas to the Filecoin network. Reload (Filecoin To IPFS):


## Features

 - **Build** (IPFS To Filecoin): Use the [MCS SDK](https://docs.filswan.com/multi-chain-storage/developer-quickstart/sdk) to upload the data to the IPFS node.
 - **Rebuild** (Filecoin To IPFS): 
	- Find the storage providers' peerIds by the [Indexer node](https://github.com/filecoin-project/index-provider) 
	- Get the storage providers IDs from the [Filecoin Network](https://github.com/filecoin-project/lotus/blob/master/api/v0api/full.go)
	- Retrieve the data from the storage providers using [Lotus](https://github.com/filecoin-project/lotus)
	- **Rebuild** the data to the IPFS nodes
* **CID-Discover (IPFS)**: Rebuilder will check if the data content available on the IPFS gateway 
* **Auto-Rebuild**: once the lost data cid is found, the auto-rebuild process will trigger and reload data from the storage provider to the target IPFS gateway.

**Work to be done**

- Auto-replica: if the current count of replicas is below a threshold of 3/5, the auto-replica process will be triggered, and the cost will be deducted from the perpetual smart contract for the deal cost to store on replica storage providers. 

- Funding contract: auto payment for a replica, funds lock, refund

- Further work: used fund for group storage insurance, locked fund in defi yield farming

## Installation

> :bell:**go 1.17+** is required

 - Install and Compile
```shell
git clone https://github.com/Fogmeta/Rebuilder.git
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

