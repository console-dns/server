package zones

import (
	"github.com/console-dns/spec/models"
)

var RecordTypes = []string{"A", "AAAA", "TXT", "CNAME", "NS", "MX", "SRV", "CAA", "SOA"}

type DNSView struct {
	Name  string
	Index int
	Type  string
	Value string
	TTL   uint32
}

func newPrintDNSGroup(name, t, value string, index int, ttl uint32) *DNSView {
	return &DNSView{
		Name:  name,
		Type:  t,
		Index: index,
		Value: value,
		TTL:   ttl,
	}
}

func ToList(r *models.Record, group string) []*DNSView {
	result := make([]*DNSView, 0)
	for i, a := range r.A {
		result = append(result, newPrintDNSGroup(group, "A", a.Ip.String(), i, a.Ttl))
	}
	for i, a := range r.AAAA {
		result = append(result, newPrintDNSGroup(group, "AAAA", a.Ip.String(), i, a.Ttl))
	}
	for i, a := range r.TXT {
		result = append(result, newPrintDNSGroup(group, "TXT", a.Text, i, a.Ttl))
	}
	for i, a := range r.CNAME {
		result = append(result, newPrintDNSGroup(group, "CNAME", a.Host, i, a.Ttl))
	}
	for i, a := range r.NS {
		result = append(result, newPrintDNSGroup(group, "NS", a.Host, i, a.Ttl))
	}
	for i, a := range r.MX {
		result = append(result, newPrintDNSGroup(group, "MX", a.Host, i, a.Ttl))
	}
	for i, a := range r.SRV {
		result = append(result, newPrintDNSGroup(group, "SRV", a.Target, i, a.Ttl))
	}
	for i, a := range r.CAA {
		result = append(result, newPrintDNSGroup(group, "CAA", a.Value, i, 0))
	}
	if r.SOA != nil {
		result = append(result, newPrintDNSGroup(group, "SOA", r.SOA.MName, -1, r.SOA.Ttl))
	}
	return result
}
