## L2TP libnetwork driver

#### Usage

Start the driver:

```
reset && go build -o l2tp cmd/l2tp/main.go && sudo ./l2tp
```

Example usage:

```
docker network create --ipam-driver=ipam --ipam-opt config="/home/machine/config.yaml" --driver=l2tp --opt config="/home/machine/config.yaml" pear
```

Where `/home/machine/config.yaml` is:

```
lns_addr: "172.16.1.57"
subnet: "10.0.5.0/32"
ppp_username: "any"
ppp_password: "any"
```

An `xl2pd` LAC server is expected to listen at `lns_addr`; it has to assign addresses from `subnet`.

Possible options and default values are defined by:

```
type config struct {
	LNSAddr             string `required:"true" yaml:"lns_addr"`
	Subnet              string `required:"true" yaml:"subnet"`
	PPPUsername         string `required:"false" yaml:"ppp_username"`
	PPPPassword         string `required:"false" yaml:"ppp_password"`
	PPPMTU              string `required:"false" yaml:"ppp_mtu" default:"1410"`
	PPPMRU              string `required:"false" yaml:"ppp_mru" default:"1410"`
	PPPIdle             string `required:"false" yaml:"ppp_idle" default:"1800"`
	PPPConnectDelay     string `required:"false" yaml:"ppp_connect_delay" default:"5000"`
	PPPDebug            bool   `required:"false" yaml:"ppp_debug" default:"true"`
	PPPNoauth           bool   `required:"false" yaml:"ppp_noauth" default:"true"`
	PPPNoccp            bool   `required:"false" yaml:"ppp_noccp" default:"true"`
	PPPDefaultRoute     bool   `required:"false" yaml:"ppp_default_route" default:"true"`
	PPPUsepeerdns       bool   `required:"false" yaml:"ppp_use_peer_dns" default:"true"`
	PPPLock             bool   `required:"false" yaml:"ppp_lock" default:"true"`
	PPPIPCPAcceptLocal  bool   `required:"false" yaml:"ppp_ipcp_accept_local" default:"true"`
	PPPIPCPAcceptRemote bool   `required:"false" yaml:"ppp_ipcp_accept_remote" default:"true"`
	PPPRefuseEAP        bool   `required:"false" yaml:"ppp_refuse_eap" default:"true"`
	PPPRequireMSChapV2  bool   `required:"false" yaml:"ppp_require_mschap_v2" default:"true"`
}
```

Run container like this:

```
docker run -itd --network="pear" httpd
```

A `ppp`-interface will be present inside the container. You can start talking to hosts from the private network right away.