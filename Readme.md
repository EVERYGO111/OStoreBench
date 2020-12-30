OStoreBench(Object Store Benchmark) is a benchmark tool for distributed object storage systems. It supports Ceph、Openstack Swift and Seaweedfs now. OStoreBench can simulate application workloads in real world and replay log of applications.
OStoreBench has been accepted by 2020 BenchCouncil International Symposium on Benchmarking, Measuring and Optimizing(bench'20).  But it  has not yet been published. In this paper, we implement OStoreBench, a scenario benchmark suite for distributed object storage systems and evaluate three state-of-the-practice object storage systems with OStoreBench. In the future, we plan to solve the hotspots problem.

## Architecture
The following figure is the main components of OStoreBench. "Controller" receives user configuration and generates requests entries. The requests rate and requests body size satisfy specific distribution by using "Distribution Geneator". The request entries will be put intto "Request Queue" while the "Executor" fetches request entry continually to generate real requests to specific storage system.
![](http://7sbpmg.com1.z0.glb.clouddn.com/blog/images/cfsb_impl.png)

## Requirement
go version >= 1.13.x

## How to Use
1. Clone the project
```
git clone https://github.com/KGXarwen/COSB.git
```
2. cd COSB/

3. go mod init example.com/m

4. make deps

5. Complie
```
make [linux]
```
6. Run
```
./linux/main/cfsb -target weed -server 172.16.1.67:9333 -t OnlineService -reqnum 100000 -c 32
```
7. Help
```
./cosb -h
```
8. If error:
```
weed_driver.go:16:33: too many arguments in call to goseaweedfs.NewSeaweed
	have (string, string, nil, number, time.Duration)
	want (string, []string, int64, *http.Client)
weed_driver.go:16:33: multiple-value goseaweedfs.NewSeaweed() in single-value context
weed_driver.go:22:37: w.client.DownloadFile undefined (type *goseaweedfs.Seaweed has no field or method DownloadFile)
weed_driver.go:31:17: assignment mismatch: 3 variables but w.client.Upload returns 2 values
```
Please modify the version of goseaweedfs in COSB/go.mod to v0.1.1
