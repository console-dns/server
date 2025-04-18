package utils

import (
	"github.com/console-dns/spec/utils"
	spec_utils "github.com/console-dns/spec/utils"
)

var RegexIDName utils.Regex = "[_A-Za-z0-9\\-]{3,15}"

var DataRegex = map[string]utils.Regex{
	"IDName":  RegexIDName,
	"DnsType": spec_utils.RegexDnsType,
	"DnsName": spec_utils.RegexDnsName,
	"Ipv4":    spec_utils.RegexIPv4,
	"Ipv6":    spec_utils.RegexIPv6,
	"Host":    spec_utils.RegexHost,
	"Mail":    spec_utils.RegexMail,
}
