[![Build Status](https://github.com/PacketEngine/packetengine/workflows/Go/badge.svg?branch=main)](https://github.com/PacketEngine/packetengine/actions?query=branch%3Amain)
[![Release](https://img.shields.io/github/release/PacketEngine/packetengine.svg)](https://github.com/PacketEngine/packetengine/releases)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/PacketEngine/packetengine)
[![Go Report Card](https://goreportcard.com/badge/github.com/PacketEngine/packetengine)](https://goreportcard.com/report/github.com/PacketEngine/packetengine)
![GitHub](https://img.shields.io/github/license/PacketEngine/packetengine)
![GitHub issues](https://img.shields.io/github/issues/PacketEngine/packetengine)

---

PacketEngine is a service which makes subdomain enumeration easy. It handles passive and active enumeration, wordlists, bruteforcing, alterations and permutations, wildcards, alerts, and much more. It's completely automated, leaving you to focus on research and hunting rather than setting up infrastructure and monitoring.

To get started, you'll need a [PacketEngine account](https://packetengine.co.uk) and an API key, and have at least one domain added.

# Features

Currently, the PacketEngine API provides read-only endpoints. The main goal of the CLI application is to be able to be integrated with other tools, such as [httpx](https://github.com/projectdiscovery/httpx), [nuclei](https://github.com/projectdiscovery/nuclei), [bbrf](https://github.com/honoki/bbrf-client), and others.

The CLI application doesn't perform any scanning itself -- it gives you instant subdomains. The PacketEngine service performs the scans automatically throughout the day.

# Usage

Standalone:

```console
$ packetengine init your-api-key
API key set!
$ packetengine subdomains uberinternal.com
devpod-us-or.uberinternal.com
devpod-us.uberinternal.com
crl.pki.uberinternal.com
stack.uberinternal.com
metal-api-preprod.uberinternal.com
creativeproduction.uberinternal.com
ussh.uberinternal.com
team-dev.uberinternal.com
flyte-poc.uberinternal.com
upt-staging.uberinternal.com
lzc-crane.uberinternal.com
stash.uberinternal.com
productops.uberinternal.com
[...]
```

Docker:

```console
$ docker run -v packetengine-config:/home/packetengine/.config/packetengine packetengine/packetengine init your-api-key
API key set!
$ docker run -v packetengine-config:/home/packetengine/.config/packetengine packetengine/packetengine subdomains uberinternal.com
devpod-us-or.uberinternal.com
devpod-us.uberinternal.com
crl.pki.uberinternal.com
stack.uberinternal.com
metal-api-preprod.uberinternal.com
creativeproduction.uberinternal.com
ussh.uberinternal.com
team-dev.uberinternal.com
flyte-poc.uberinternal.com
upt-staging.uberinternal.com
lzc-crane.uberinternal.com
stash.uberinternal.com
productops.uberinternal.com
[...]

```

httpx:

```console
$ packetengine subdomains uberinternal.com | httpx -silent -status-code -title -mc 200
https://arize.uberinternal.com [200] [Arize AI]
https://emea-vpn-lab.uberinternal.com [200] []
https://chef-staging.uberinternal.com [200] [Chef Automate]
https://chef.uberinternal.com [200] [Chef Automate]
https://metal-api-preprod.uberinternal.com [200] []
https://metal-api-staging.uberinternal.com [200] []
https://metal-api.uberinternal.com [200] []
https://rfa.uberinternal.com [200] [Nuix ECC Server]
https://vpn-emea-any.uberinternal.com [200] []
https://vpn-world-any.awscorp.uberinternal.com [200] []
https://vpn-amere-any.uberinternal.com [200] []
https://vpn-amerw-any.uberinternal.com [200] []
https://vpn-apac-any.uberinternal.com [200] []
[...]
```

nuclei:

```console
packetengine subdomains uberinternal.com | httpx -silent | nuclei
```

bbrf:

```console
packetengine subdomains uberinternal.com | bbrf domain add -
```

# Without Tags

PacketEngine automatically tags subdomains if their DNS records point to private IP space or IPv6 IP addresses. If you want to exclude any tags you can use the `--without-tags` flag.

```console
packetengine subdomains --without-tags=ipv6,private-ip uberinternal.com | httpx -silent -status-code -title -mc 200
```

# Installation

Using `go install`:

```console
go install -v github.com/PacketEngine/packetengine/cmd/packetengine@latest
```

Using Docker:

```console
docker pull packetengine/packetengine:latest
```

Using Snap:

TODO

Using Brew:

```console
brew tap PacketEngine/packetengine
brew install packetengine
```

# License

The PacketEngine CLI is available under the MIT license. See the LICENSE file for more info.
