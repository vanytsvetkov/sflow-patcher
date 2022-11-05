# sFlow proxy

## Introduction

`sflow-patcher` is a lightweight sFlow proxy that strips VXLAN headers from the sFlow raw packet records, keeping all other records and counters intact. All UDP packets that `sflow-patcher` failed to parse as sFlow datagrams are relayed as is. UDP source addresses and ports are preserved in order to provide compatibility with sFlow analyzers thar rely on source IP instead of agendId value.

The VxLAN header detection is dead-simple: `sflow-patcher` goes through the raw packet record headers and checks if the packet happens to be a UDP packet with destination port 4789.

## Building

Install Go >=1.13, run `make`, and grab the binary from `out` directory in the repo. `sflow-patcher` uses pcap for emitting raw frames, so make sure you have `libpcap` headers installed.

## Usage

```
Usage:
  sflow-patcher [flags]

Flags:
  -b, --bind string        address and port to bind on (default "0.0.0.0:5000")
  -s, --buffer-size int    input buffer size in bytes (default 1500)
  -d, --debug              enable debug logging
  -m, --dst-mac string     destination MAC address
  -h, --help               help for sflow-patcher
  -i, --out-if string      outgoing interface
  -r, --route-map string   path to the collector route map file
  -v, --vlan-map string    path to the vlan map file
  -w, --workers int        number of workers (default 10)
```

`dst-mac`, `out-if`, `route-map` and `vlan-map` parameters are mandatory. `buffer-size` should fit any received sFlow datagram. `workers` indicates how many packets will be processed in parallel (could be increased in case of packet drops).

### Example

`sflow-patcher -b 0.0.0.0:16789 -i eth0 -m 00:11:22:33:44:55 -r /etc/sflow-patcher-routes.yaml -v /etc/sflow-patcher-vlans.yaml` will listen on UDP port 16789 and route processed packets according to the map specified in `/etc/sflow-patcher-routes.yaml`, sending ethernet frames to MAC address 00:11:22:33:44:55 from interface eth0.

## Route & VLAN maps

Route and VLAN maps specifier the routing & filtering of sFlow datagrams in YAML format. It is a simple dictionaries: with agent addresses as keys and collector addresses and ports as values for route-map case; and with VIDs as keys and boolean-like rules as values for vlan-map case.

> You can send a SIGHUP to `sflow-patcher` in order to re-read map files without restarting the process.

### Example

#### route-map file (e.g. `/etc/sflow-patcher-routes.yaml`)

```yaml
---
default:    172.16.20.100:5000
172.18.3.1: 172.16.20.101:5000
172.18.3.2: 172.16.20.102:5000
172.18.3.3: 172.16.20.102:6000
```

A special agent address `default` could be used to specify a catch-all collector. 

#### vlan-map file (e.g. `/etc/sflow-patcher-vlans.yaml`)
```yaml
---
0:    false
450:  true
1000: true
4095: false
```

Specify a value of false for some VLAN ID in order to drop sFlow samples learned in the specified VLAN.<br>
Note: Packets with an empty data field **will not be retransmitted**.


