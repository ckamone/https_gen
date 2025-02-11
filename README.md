# base info
http custom generator
tested with go1.18.1 linux/amd64

## set up src ips
```bash
ip a add 192.168.0.101/24 dev ens19; \
ip a add 192.168.0.102/24 dev ens19; \
ip a add 192.168.0.103/24 dev ens19; \
ip a add 192.168.0.104/24 dev ens19; \
ip a add 192.168.0.105/24 dev ens19;
```

## some os optimisations
```bash
sysctl net.ipv4.tcp_max_orphans=65535; \
sysctl net.ipv4.tcp_max_tw_buckets=65536; \
sysctl net.netfilter.nf_conntrack_max=10000000; \
ulimit -n 1000000;
```

## testing
`go run main.go`
or
`./https_gen -cps=10 -uri=1kb.html -ips=192.168.0.100,192.168.0.101 -log=tst.log`

## building
`go build -o https_gen main.go`