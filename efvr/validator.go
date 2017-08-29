package efvr

import (
	"net"
)

var Funcs = map[string]func(string) bool{
	"cidrv4": IsCIDRv4,
	"cidrv6": IsCIDRv6,
	"cidr":   IsCIDR,
	"mac":    IsMAC,
}

// IsCIDR check if the string is an valid CIDR notiation (IPV4 & IPV6)
func IsCIDR(str string) bool {
	_, _, err := net.ParseCIDR(str)
	return err == nil
}

func IsCIDRv4(str string) bool {
	ip, _, err := net.ParseCIDR(str)
	return err == nil && ip.To4() != nil
}

func IsCIDRv6(str string) bool {
	ip, _, err := net.ParseCIDR(str)
	return err == nil && ip.To4() == nil
}

func IsMAC(val string) bool {
	_, err := net.ParseMAC(val)
	return err == nil
}
