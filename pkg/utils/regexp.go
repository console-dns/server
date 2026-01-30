package utils

import (
	"fmt"
	"regexp"
)

type Regex string

func (r Regex) Valid(src string) error {
	if regexp.MustCompile(string(r)).MatchString(src) {
		return nil
	}
	return fmt.Errorf("格式错误: %s", src)
}

var RegexIDName Regex = `[_A-Za-z0-9\-]{3,15}`

var (
	RegexHost    = Regex(`^([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])(\.([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]{0,61}[a-zA-Z0-9]))+$`)
	RegexMail    = Regex(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`)
	RegexDnsName = Regex(`^([a-zA-Z0-9_]{1,61})|(\*)$`)
	RegexDnsType = Regex(`^(A|AAAA|TXT|CNAME|NS|MX|SRV|CAA|SOA)$`)
	RegexIPv4    = Regex(`^((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.?\b){4}$`)
	RegexIPv6    = Regex(`^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$`)
)

var DataRegex = map[string]Regex{
	"IDName":  RegexIDName,
	"DnsType": RegexDnsType,
	"DnsName": RegexDnsName,
	"Ipv4":    RegexIPv4,
	"Ipv6":    RegexIPv6,
	"Host":    RegexHost,
	"Mail":    RegexMail,
}
