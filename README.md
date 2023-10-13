<p align="center">
    <a href="https://netmoth.com" target="_blank" rel="noopener">
        <img src="https://github.com/netmoth/.github/raw/main/img/logo.svg" alt="A lightweight, fast, simple and complete solution for traffic analysis and intrusion detection" width="50%" />
    </a>
</p>

<p align="center">
    <a href="https://github.com/netmoth/netmoth/releases">
    <img src="https://img.shields.io/github/v/release/netmoth/netmoth?sort=semver&label=Release&color=651FFF" />
    </a>
    &nbsp;
    <a href="/LICENSE"><img src="https://img.shields.io/badge/license-Apache--2.0-green.svg"></a>
    &nbsp;
    <a href="https://goreportcard.com/report/github.com/netmoth/netmoth"><img src="https://goreportcard.com/badge/github.com/netmoth/netmoth"></a>
    &nbsp;
    <a href="https://www.codefactor.io/repository/github/netmoth/netmoth"><img src="https://www.codefactor.io/repository/github/netmoth/netmoth/badge" alt="CodeFactor" /></a>
    &nbsp;
    <a href="https://github.com/netmoth/netmoth"><img src="https://img.shields.io/badge/backend-go-orange.svg"></a>
    &nbsp;
    <a href="https://github.com/netmoth/netmoth/blob/main/go.mod"><img src="https://img.shields.io/github/go-mod/go-version/netmoth/netmoth?color=7fd5ea"></a>
</p>

---

## &nbsp;&nbsp;What is netmoth?

Netmoth is a lightweight, fast, simple and complete solution for traffic analysis and intrusion detection.

> ⚠️&nbsp;&nbsp;Current major version is zero (`v0.x.x`) to accommodate rapid development and fast iteration while getting early feedback from users. Please keep in mind that netmoth is still under active development and therefore full backward compatibility is not guaranteed before reaching v1.0.0.


## 🏆&nbsp;&nbsp;Features

- [x] Monitors traffic on all interfaces
- [x] Minimal configuration
- [x] PCAP
- [x] AF_PACKET 
- [x] PF_RING
- [ ] eBPF
- [x] Zero copy packet processing
- [x] Automatic TCP stream reassembly
- [x] Berkeley Packet Filter
- [x] Check IP on blocklist
- [x] Checking botnet on blocklist
- [ ] Checking certificate on blocklist
- [x] Checking tracker on blocklist
- [ ] Web-interface
- [ ] Rules
- [ ] Agents

## 📚&nbsp;&nbsp;Documentation

... coming soon ...

## 🏁&nbsp;&nbsp;Installation

Simple agent installation
```bash
mkdir netmoth
cd ./netmoth
curl -L https://raw.githubusercontent.com/netmoth/netmoth/main/config_example.yml > config.yml
curl -L https://github.com/netmoth/netmoth/releases/latest/download/netmoth_agent_Linux_x86_64 > netmoth_agent
sudo chmod u+x netmoth_agent
```

if necessary, make changes to the `config.yml` file, then run the agent
```bash
./netmoth_agent
```


## 👑&nbsp;&nbsp;Community

... coming soon ...

## 👍&nbsp;&nbsp;Contribute

We would for you to get involved with netmoth development! If you want to say **thank you** and/or support the active development of `netmoth`:

1. Add a [GitHub Star](https://github.com/netmoth/netmoth/stargazers) to the project.
2. Tweet about the project [on your Twitter](https%3A//twitter.com/intent/tweet?text%3D%F0%9F%9A%80%20netmoth%20-%20a%20lightweight%2C%20fast%2C%20simple%20and%20complete%20solution%20for%20traffic%20analysis%20and%20intrusion%20detection%20on%20%23Go%20https%3A//github.com/netmoth/netmoth).
3. Write a review or tutorial on [Medium](https://medium.com/), [Dev.to](https://dev.to/) or personal blog.

You can learn more about how you can contribute to this project in the [contribution guide](CONTRIBUTING.md).

## 🚨&nbsp;&nbsp;Security

... coming soon ...

