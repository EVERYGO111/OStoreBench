COSB(Cloud Object Storage Benchmark) is a benchmark tool for distributed object storage systems. It supports Ceph、Openstack Swift and Seaweedfs now. COSB can simulate application workloads in real world and replay log of applications.

## Architecture
The following figure is the main components of COSB. "Controller" receives user configuration and generates requests entries. The requests rate and requests body size satisfy specific distribution by using "Distribution Geneator". The request entries will be put intto "Request Queue" while the "Executor" fetches request entry continually to generate real requests to specific storage system.
![](http://7sbpmg.com1.z0.glb.clouddn.com/blog/images/cfsb_impl.png)

## How to Use
1. Clone the project
```
git clone git@github.com:KDF5000/COSB.git
```
2. Complie
```
make [linux]
```
3. Run
```
./cosb -target weed -server 172.16.1.67:9333 -t OnlineService -reqnum 100000 -c 32
```
4. Help
```
./cosb -h
```
