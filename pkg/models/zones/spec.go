package zones

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"regexp"
	"strconv"

	"github.com/console-dns/server/pkg/utils"
	"github.com/pkg/errors"
	"gopkg.d7z.net/middleware/kv"
)

func parseTtl(ttl string) (uint32, error) {
	t, err := strconv.Atoi(ttl)
	if err != nil {
		return 0, err
	}
	if t < 0 {
		return 0, errors.New("invalid ttl")
	}
	return uint32(t), nil
}

func atoUint32(src string) (uint32, error) {
	r, err := strconv.Atoi(src)
	if err != nil {
		return 0, err
	}
	if r < 0 {
		return 0, fmt.Errorf("%s is not a valid uint32", src)
	}
	return uint32(r), nil
}

func validateHost(host string) error {
	if !regexp.MustCompile(string(utils.RegexHost)).MatchString(host) {
		return errors.New("域名格式错误")
	}
	return nil
}

func validateMail(mail string) error {
	if !regexp.MustCompile(string(utils.RegexMail)).MatchString(mail) {
		return errors.New("邮件格式错误")
	}
	return nil
}

type Validator interface {
	Validate() error
	Type() string
}

type Clone[T any] interface {
	Clone() T
}

type A struct {
	Ttl uint32 `json:"ttl" yaml:"ttl" toml:"ttl"`
	Ip  net.IP `json:"ip" yaml:"ip" toml:"ip"`
}

func (r *A) Type() string { return "A" }

func (r *A) Clone() *A {
	var ip net.IP
	if r.Ip != nil {
		ip = make(net.IP, len(r.Ip))
		copy(ip, r.Ip)
	}
	return &A{
		Ttl: r.Ttl,
		Ip:  ip,
	}
}

func (r *A) Validate() error {
	if r.Ip == nil {
		return errors.New("IP 不能为空")
	}
	if r.Ip.To4() == nil {
		return errors.New("不是合法的 IPv4 地址")
	}
	return nil
}

func (r *A) UnmarshalJSON(data []byte) error {
	type Alias A
	var aux Alias
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	*r = A(aux)
	return r.Validate()
}

func (r *A) MarshalJSON() ([]byte, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}
	type Alias A
	return json.Marshal((*Alias)(r))
}

func NewA(ip, ttl string) (*A, error) {
	t, err := parseTtl(ttl)
	if err != nil {
		return nil, err
	}
	res := &A{
		Ttl: t,
		Ip:  net.ParseIP(ip),
	}
	if err := res.Validate(); err != nil {
		return nil, err
	}
	return res, nil
}

func FromRecordA(f func(string) string) (*A, error) {
	return NewA(f("ip"), f("ttl"))
}

type AAAA struct {
	Ttl uint32 `json:"ttl" yaml:"ttl" toml:"ttl"`
	Ip  net.IP `json:"ip" yaml:"ip" toml:"ip"`
}

func (r *AAAA) Type() string { return "AAAA" }

func (r *AAAA) Clone() *AAAA {
	var ip net.IP
	if r.Ip != nil {
		ip = make(net.IP, len(r.Ip))
		copy(ip, r.Ip)
	}
	return &AAAA{
		Ttl: r.Ttl,
		Ip:  ip,
	}
}

func (r *AAAA) Validate() error {
	if r.Ip == nil {
		return errors.New("IP 不能为空")
	}
	if r.Ip.To4() != nil {
		return errors.New("不是合法的 IPv6 地址")
	}
	return nil
}

func (r *AAAA) UnmarshalJSON(data []byte) error {
	type Alias AAAA
	var aux Alias
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	*r = AAAA(aux)
	return r.Validate()
}

func (r *AAAA) MarshalJSON() ([]byte, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}
	type Alias AAAA
	return json.Marshal((*Alias)(r))
}

func NewAAAA(ip, ttl string) (*AAAA, error) {
	u, err := parseTtl(ttl)
	if err != nil {
		return nil, err
	}
	res := &AAAA{
		Ttl: u,
		Ip:  net.ParseIP(ip),
	}
	if err := res.Validate(); err != nil {
		return nil, err
	}
	return res, nil
}

func FromRecordAAAA(f func(string) string) (*AAAA, error) {
	return NewAAAA(f("ip"), f("ttl"))
}

type CAA struct {
	Flag  uint8  `json:"flag"`
	Tag   string `json:"tag"`
	Value string `json:"value"`
}

func (r *CAA) Type() string { return "CAA" }

func (r *CAA) Clone() *CAA {
	return &CAA{
		Flag:  r.Flag,
		Tag:   r.Tag,
		Value: r.Value,
	}
}

func (r *CAA) Validate() error {
	if r.Flag != 0 && r.Flag != 128 {
		return errors.New("未知 flag 类型")
	}
	if len(r.Value) == 0 {
		return errors.New("value 不能为空")
	}
	switch r.Tag {
	case "issue", "issuewild", "iodef", "contactphone":
	case "contactemail":
		if err := validateMail(r.Value); err != nil {
			return err
		}
	default:
		return errors.New("未知 tag 类型")
	}
	return nil
}

func (r *CAA) UnmarshalJSON(data []byte) error {
	type Alias CAA
	var aux Alias
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	*r = CAA(aux)
	return r.Validate()
}

func (r *CAA) MarshalJSON() ([]byte, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}
	type Alias CAA
	return json.Marshal((*Alias)(r))
}

func NewCAA(flag, tag, value string) (*CAA, error) {
	f, err := strconv.Atoi(flag)
	if err != nil {
		return nil, err
	}
	res := &CAA{
		Flag:  uint8(f),
		Tag:   tag,
		Value: value,
	}
	if err := res.Validate(); err != nil {
		return nil, err
	}
	return res, nil
}

func FromRecordCAA(f func(string) string) (*CAA, error) {
	return NewCAA(f("flag"), f("tag"), f("value"))
}

type CNAME struct {
	Ttl  uint32 `json:"ttl"`
	Host string `json:"host"`
}

func (r *CNAME) Type() string { return "CNAME" }

func (r *CNAME) Clone() *CNAME {
	return &CNAME{
		Ttl:  r.Ttl,
		Host: r.Host,
	}
}

func (r *CNAME) Validate() error {
	return validateHost(r.Host)
}

func (r *CNAME) UnmarshalJSON(data []byte) error {
	type Alias CNAME
	var aux Alias
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	*r = CNAME(aux)
	return r.Validate()
}

func (r *CNAME) MarshalJSON() ([]byte, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}
	type Alias CNAME
	return json.Marshal((*Alias)(r))
}

func NewCNAME(host, ttl string) (*CNAME, error) {
	u, err := parseTtl(ttl)
	if err != nil {
		return nil, err
	}
	res := &CNAME{
		Ttl:  u,
		Host: host,
	}
	if err := res.Validate(); err != nil {
		return nil, err
	}
	return res, nil
}

func FromRecordCNAME(f func(string) string) (*CNAME, error) {
	return NewCNAME(f("host"), f("ttl"))
}

type MX struct {
	Ttl        uint32 `json:"ttl"`
	Host       string `json:"host"`
	Preference uint16 `json:"preference"`
}

func (r *MX) Type() string { return "MX" }

func (r *MX) Clone() *MX {
	return &MX{
		Ttl:        r.Ttl,
		Host:       r.Host,
		Preference: r.Preference,
	}
}

func (r *MX) Validate() error {
	return validateHost(r.Host)
}

func (r *MX) UnmarshalJSON(data []byte) error {
	type Alias MX
	var aux Alias
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	*r = MX(aux)
	return r.Validate()
}

func (r *MX) MarshalJSON() ([]byte, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}
	type Alias MX
	return json.Marshal((*Alias)(r))
}

func NewMX(host, preference, ttl string) (*MX, error) {
	u, err := parseTtl(ttl)
	if err != nil {
		return nil, err
	}
	p, err := strconv.Atoi(preference)
	if err != nil {
		return nil, err
	}
	if p > math.MaxUint16 || p < 0 {
		return nil, errors.New("preference 区间不合法")
	}
	res := &MX{
		Ttl:        u,
		Host:       host,
		Preference: uint16(p),
	}
	if err := res.Validate(); err != nil {
		return nil, err
	}
	return res, nil
}

func FromRecordMX(f func(string) string) (*MX, error) {
	return NewMX(f("host"), f("preference"), f("ttl"))
}

type NS struct {
	Ttl  uint32 `json:"ttl"`
	Host string `json:"host"`
}

func (r *NS) Type() string { return "NS" }

func (r *NS) Clone() *NS {
	return &NS{
		Ttl:  r.Ttl,
		Host: r.Host,
	}
}

func (r *NS) Validate() error {
	return validateHost(r.Host)
}

func (r *NS) UnmarshalJSON(data []byte) error {
	type Alias NS
	var aux Alias
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	*r = NS(aux)
	return r.Validate()
}

func (r *NS) MarshalJSON() ([]byte, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}
	type Alias NS
	return json.Marshal((*Alias)(r))
}

func NewNS(host, ttl string) (*NS, error) {
	u, err := parseTtl(ttl)
	if err != nil {
		return nil, err
	}
	res := &NS{
		Ttl:  u,
		Host: host,
	}
	if err := res.Validate(); err != nil {
		return nil, err
	}
	return res, nil
}

func FromRecordNS(f func(string) string) (*NS, error) {
	return NewNS(f("host"), f("ttl"))
}

type SOA struct {
	Ttl     uint32 `json:"ttl"`
	MName   string `json:"mname"`
	RName   string `json:"rname"`
	Serial  uint32 `json:"serial"`
	Refresh uint32 `json:"refresh"`
	Retry   uint32 `json:"retry"`
	Expire  uint32 `json:"expire"`
	Minimum uint32 `json:"minimum"`
}

func (r *SOA) Type() string { return "SOA" }

func (r *SOA) Clone() *SOA {
	return &SOA{
		Ttl:     r.Ttl,
		MName:   r.MName,
		RName:   r.RName,
		Serial:  r.Serial,
		Refresh: r.Refresh,
		Retry:   r.Retry,
		Expire:  r.Expire,
		Minimum: r.Minimum,
	}
}

func (r *SOA) Validate() error {
	if err := validateHost(r.MName); err != nil {
		return err
	}
	if err := validateHost(r.RName); err != nil {
		return err
	}
	if r.Retry >= r.Refresh {
		return fmt.Errorf("retry (%d) must be less than refresh (%d)", r.Retry, r.Refresh)
	}
	if r.Expire <= r.Refresh+r.Retry {
		return fmt.Errorf("expire (%d) must be greater than refresh + retry (%d)", r.Expire, r.Refresh+r.Retry)
	}
	return nil
}

func (r *SOA) UnmarshalJSON(data []byte) error {
	type Alias SOA
	var aux Alias
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	*r = SOA(aux)
	return r.Validate()
}

func (r *SOA) MarshalJSON() ([]byte, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}
	type Alias SOA
	return json.Marshal((*Alias)(r))
}

func NewSOA(mName, rName, serial, refresh, retry, expire, minimum, ttl string) (*SOA, error) {
	u, err := parseTtl(ttl)
	if err != nil {
		return nil, err
	}
	serialInt, err := atoUint32(serial)
	if err != nil {
		return nil, err
	}
	refreshInt, err := atoUint32(refresh)
	if err != nil {
		return nil, err
	}
	retryInt, err := atoUint32(retry)
	if err != nil {
		return nil, err
	}
	expireInt, err := atoUint32(expire)
	if err != nil {
		return nil, err
	}
	minimumInt, err := atoUint32(minimum)
	if err != nil {
		return nil, err
	}
	res := &SOA{
		Ttl:     u,
		MName:   mName,
		RName:   rName,
		Serial:  serialInt,
		Refresh: refreshInt,
		Retry:   retryInt,
		Expire:  expireInt,
		Minimum: minimumInt,
	}
	if err := res.Validate(); err != nil {
		return nil, err
	}
	return res, nil
}

func FromRecordSOA(f func(string) string) (*SOA, error) {
	return NewSOA(f("mname"), f("rname"), f("serial"), f("refresh"), f("retry"), f("expire"), f("minimum"), f("ttl"))
}

type SRV struct {
	Ttl      uint32 `json:"ttl"`
	Priority uint16 `json:"priority"`
	Weight   uint16 `json:"weight"`
	Port     uint16 `json:"port"`
	Target   string `json:"target"`
}

func (r *SRV) Type() string { return "SRV" }

func (r *SRV) Clone() *SRV {
	return &SRV{
		Ttl:      r.Ttl,
		Priority: r.Priority,
		Weight:   r.Weight,
		Port:     r.Port,
		Target:   r.Target,
	}
}

func (r *SRV) Validate() error {
	return validateHost(r.Target)
}

func (r *SRV) UnmarshalJSON(data []byte) error {
	type Alias SRV
	var aux Alias
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	*r = SRV(aux)
	return r.Validate()
}

func (r *SRV) MarshalJSON() ([]byte, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}
	type Alias SRV
	return json.Marshal((*Alias)(r))
}

func NewSRV(priority, weight, port, target, ttl string) (*SRV, error) {
	u, err := parseTtl(ttl)
	if err != nil {
		return nil, err
	}
	p, err := strconv.Atoi(priority)
	if err != nil {
		return nil, err
	}
	w, err := strconv.Atoi(weight)
	if err != nil {
		return nil, err
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}
	if p > math.MaxUint16 || p < 0 {
		return nil, errors.New("priority 区间不合法")
	}
	if w > math.MaxUint16 || w < 0 {
		return nil, errors.New("weight 区间不合法")
	}
	if portInt > math.MaxUint16 || portInt < 0 {
		return nil, errors.New("port 区间不合法")
	}
	res := &SRV{
		Ttl:      u,
		Priority: uint16(p),
		Weight:   uint16(w),
		Port:     uint16(portInt),
		Target:   target,
	}
	if err := res.Validate(); err != nil {
		return nil, err
	}
	return res, nil
}

func FromRecordSRV(f func(string) string) (*SRV, error) {
	return NewSRV(f("priority"), f("weight"), f("port"), f("target"), f("ttl"))
}

type TXT struct {
	Ttl  uint32 `json:"ttl"`
	Text string `json:"text"`
}

func (r *TXT) Type() string { return "TXT" }

func (r *TXT) Clone() *TXT {
	return &TXT{
		Ttl:  r.Ttl,
		Text: r.Text,
	}
}

func (r *TXT) Validate() error {
	return nil
}

func (r *TXT) UnmarshalJSON(data []byte) error {
	type Alias TXT
	var aux Alias
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	*r = TXT(aux)
	return r.Validate()
}

func (r *TXT) MarshalJSON() ([]byte, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}
	type Alias TXT
	return json.Marshal((*Alias)(r))
}

func NewTXT(txt, ttl string) (*TXT, error) {
	u, err := parseTtl(ttl)
	if err != nil {
		return nil, err
	}
	res := &TXT{
		Ttl:  u,
		Text: txt,
	}
	if err := res.Validate(); err != nil {
		return nil, err
	}
	return res, nil
}

func FromRecordTXT(f func(string) string) (*TXT, error) {
	return NewTXT(f("text"), f("ttl"))
}

type Record struct {
	A     []*A     `json:"A,omitempty"`
	AAAA  []*AAAA  `json:"AAAA,omitempty"`
	TXT   []*TXT   `json:"TXT,omitempty"`
	CNAME []*CNAME `json:"CNAME,omitempty"`
	NS    []*NS    `json:"NS,omitempty"`
	MX    []*MX    `json:"MX,omitempty"`
	SRV   []*SRV   `json:"SRV,omitempty"`
	CAA   []*CAA   `json:"CAA,omitempty"`
	SOA   *SOA     `json:"SOA,omitempty"`
}

func NewRecord() *Record {
	return &Record{
		A:     make([]*A, 0),
		AAAA:  make([]*AAAA, 0),
		TXT:   make([]*TXT, 0),
		CNAME: make([]*CNAME, 0),
		NS:    make([]*NS, 0),
		MX:    make([]*MX, 0),
		SRV:   make([]*SRV, 0),
		CAA:   make([]*CAA, 0),
	}
}

func (r *Record) Save(ctx context.Context, storage kv.KV) error {
	_ = utils.AutoKVMarshal(ctx, storage, "A", r.A)
	_ = utils.AutoKVMarshal(ctx, storage, "AAAA", r.AAAA)
	_ = utils.AutoKVMarshal(ctx, storage, "TXT", r.TXT)
	_ = utils.AutoKVMarshal(ctx, storage, "CNAME", r.CNAME)
	_ = utils.AutoKVMarshal(ctx, storage, "NS", r.NS)
	_ = utils.AutoKVMarshal(ctx, storage, "MX", r.MX)
	_ = utils.AutoKVMarshal(ctx, storage, "SRV", r.SRV)
	_ = utils.AutoKVMarshal(ctx, storage, "CAA", r.CAA)
	_ = utils.AutoKVMarshal(ctx, storage, "SOA", r.SOA)
	return nil
}

func (r *Record) Load(ctx context.Context, storage kv.KV) error {
	_ = utils.AutoKVUnmarshal(ctx, storage, "A", &r.A)
	_ = utils.AutoKVUnmarshal(ctx, storage, "AAAA", &r.AAAA)
	_ = utils.AutoKVUnmarshal(ctx, storage, "TXT", &r.TXT)
	_ = utils.AutoKVUnmarshal(ctx, storage, "CNAME", &r.CNAME)
	_ = utils.AutoKVUnmarshal(ctx, storage, "NS", &r.NS)
	_ = utils.AutoKVUnmarshal(ctx, storage, "MX", &r.MX)
	_ = utils.AutoKVUnmarshal(ctx, storage, "SRV", &r.SRV)
	_ = utils.AutoKVUnmarshal(ctx, storage, "CAA", &r.CAA)
	_ = utils.AutoKVUnmarshal(ctx, storage, "SOA", &r.SOA)
	return nil
}

func (r *Record) ModA(old, new *A) error {
	if old == nil {
		r.A = append(r.A, new)
		return nil
	}
	for i, a := range r.A {
		if a.Ip.Equal(old.Ip) {
			r.A[i] = new
			return nil
		}
	}
	return errors.New("record not found")
}

func (r *Record) RemoveA(old *A) error {
	for i, a := range r.A {
		if a.Ip.Equal(old.Ip) {
			r.A = append(r.A[:i], r.A[i+1:]...)
			return nil
		}
	}
	return errors.New("record not found")
}

func (r *Record) ModAAAA(old, new *AAAA) error {
	if old == nil {
		r.AAAA = append(r.AAAA, new)
		return nil
	}
	for i, a := range r.AAAA {
		if a.Ip.Equal(old.Ip) {
			r.AAAA[i] = new
			return nil
		}
	}
	return errors.New("record not found")
}

func (r *Record) RemoveAAAA(old *AAAA) error {
	for i, a := range r.AAAA {
		if a.Ip.Equal(old.Ip) {
			r.AAAA = append(r.AAAA[:i], r.AAAA[i+1:]...)
			return nil
		}
	}
	return errors.New("record not found")
}

func (r *Record) ModTXT(old, new *TXT) error {
	if old == nil {
		r.TXT = append(r.TXT, new)
		return nil
	}
	for i, a := range r.TXT {
		if a.Text == old.Text {
			r.TXT[i] = new
			return nil
		}
	}
	return errors.New("record not found")
}

func (r *Record) RemoveTXT(old *TXT) error {
	for i, a := range r.TXT {
		if a.Text == old.Text {
			r.TXT = append(r.TXT[:i], r.TXT[i+1:]...)
			return nil
		}
	}
	return errors.New("record not found")
}

func (r *Record) ModCNAME(old, new *CNAME) error {
	if old == nil {
		r.CNAME = append(r.CNAME, new)
		return nil
	}
	for i, a := range r.CNAME {
		if a.Host == old.Host {
			r.CNAME[i] = new
			return nil
		}
	}
	return errors.New("record not found")
}

func (r *Record) RemoveCNAME(old *CNAME) error {
	for i, a := range r.CNAME {
		if a.Host == old.Host {
			r.CNAME = append(r.CNAME[:i], r.CNAME[i+1:]...)
			return nil
		}
	}
	return errors.New("record not found")
}

func (r *Record) ModNS(old, new *NS) error {
	if old == nil {
		r.NS = append(r.NS, new)
		return nil
	}
	for i, a := range r.NS {
		if a.Host == old.Host {
			r.NS[i] = new
			return nil
		}
	}
	return errors.New("record not found")
}

func (r *Record) RemoveNS(old *NS) error {
	for i, a := range r.NS {
		if a.Host == old.Host {
			r.NS = append(r.NS[:i], r.NS[i+1:]...)
			return nil
		}
	}
	return errors.New("record not found")
}

func (r *Record) ModMX(old, new *MX) error {
	if old == nil {
		r.MX = append(r.MX, new)
		return nil
	}
	for i, a := range r.MX {
		if a.Host == old.Host && a.Preference == old.Preference {
			r.MX[i] = new
			return nil
		}
	}
	return errors.New("record not found")
}

func (r *Record) RemoveMX(old *MX) error {
	for i, a := range r.MX {
		if a.Host == old.Host && a.Preference == old.Preference {
			r.MX = append(r.MX[:i], r.MX[i+1:]...)
			return nil
		}
	}
	return errors.New("record not found")
}

func (r *Record) ModSRV(old, new *SRV) error {
	if old == nil {
		r.SRV = append(r.SRV, new)
		return nil
	}
	for i, a := range r.SRV {
		if a.Target == old.Target && a.Port == old.Port {
			r.SRV[i] = new
			return nil
		}
	}
	return errors.New("record not found")
}

func (r *Record) RemoveSRV(old *SRV) error {
	for i, a := range r.SRV {
		if a.Target == old.Target && a.Port == old.Port {
			r.SRV = append(r.SRV[:i], r.SRV[i+1:]...)
			return nil
		}
	}
	return errors.New("record not found")
}

func (r *Record) ModCAA(old, new *CAA) error {
	if old == nil {
		r.CAA = append(r.CAA, new)
		return nil
	}
	for i, a := range r.CAA {
		if a.Tag == old.Tag && a.Value == old.Value && a.Flag == old.Flag {
			r.CAA[i] = new
			return nil
		}
	}
	return errors.New("record not found")
}

func (r *Record) RemoveCAA(old *CAA) error {
	for i, a := range r.CAA {
		if a.Tag == old.Tag && a.Value == old.Value && a.Flag == old.Flag {
			r.CAA = append(r.CAA[:i], r.CAA[i+1:]...)
			return nil
		}
	}
	return errors.New("record not found")
}

func (r *Record) ModSOA(old, new *SOA) error {
	r.SOA = new
	return nil
}

func (r *Record) RemoveSOA(old *SOA) error {
	r.SOA = nil
	return nil
}

type Zone struct {
	Records map[string]*Record `json:"records"`
}

func NewZone() *Zone {
	return &Zone{
		Records: make(map[string]*Record),
	}
}

func (z *Zone) ModRecord(name string, f func(r *Record) error) error {
	r := z.Records[name]
	if r == nil {
		r = NewRecord()
		z.Records[name] = r
	}
	if err := f(r); err != nil {
		return err
	}
	if r.IsEmpty() {
		delete(z.Records, name)
	}
	return nil
}

func (z *Zone) Save(ctx context.Context, storage kv.KV) error {
	for name, record := range z.Records {
		err := record.Save(ctx, storage.Child(name))
		if err != nil {
			return err
		}
	}
	return nil
}

func (z *Zone) Load(ctx context.Context, storage kv.KV) error {
	list, err := storage.List(ctx, "")
	if err != nil {
		return err
	}
	for key := range list {
		// key is like "recordName/A"
		parts := regexp.MustCompile(`/`).Split(key, 2)
		if len(parts) > 0 {
			recordName := parts[0]
			if _, ok := z.Records[recordName]; !ok {
				record := NewRecord()
				err := record.Load(ctx, storage.Child(recordName))
				if err != nil {
					return err
				}
				z.Records[recordName] = record
			}
		}
	}
	return nil
}

type ZoneData struct {
	Data map[string]*Zone `json:"zones"`
}

func NewZones() *ZoneData {
	return &ZoneData{
		Data: make(map[string]*Zone),
	}
}

func (z *ZoneData) GetZone(name string) *Zone {
	return z.Data[name]
}

func (z *ZoneData) GetRecords(name string) map[string]*Record {
	if zone, ok := z.Data[name]; ok {
		return zone.Records
	}
	return nil
}

func (z *ZoneData) ListRecords() map[string]map[string]Record {
	res := make(map[string]map[string]Record)
	for zoneName, zone := range z.Data {
		res[zoneName] = make(map[string]Record)
		for recordName, record := range zone.Records {
			res[zoneName][recordName] = *record
		}
	}
	return res
}

func (z *ZoneData) CopyFrom(r *Record, zoneName, recordName, dnsType string) {
	zone := z.Data[zoneName]
	if zone == nil {
		zone = NewZone()
		z.Data[zoneName] = zone
	}
	record := zone.Records[recordName]
	if record == nil {
		record = NewRecord()
		zone.Records[recordName] = record
	}
	switch dnsType {
	case "A":
		record.A = append(record.A, r.A...)
	case "AAAA":
		record.AAAA = append(record.AAAA, r.AAAA...)
	case "TXT":
		record.TXT = append(record.TXT, r.TXT...)
	case "CNAME":
		record.CNAME = append(record.CNAME, r.CNAME...)
	case "NS":
		record.NS = append(record.NS, r.NS...)
	case "MX":
		record.MX = append(record.MX, r.MX...)
	case "SRV":
		record.SRV = append(record.SRV, r.SRV...)
	case "CAA":
		record.CAA = append(record.CAA, r.CAA...)
	case "SOA":
		record.SOA = r.SOA
	}
}

func (z *ZoneData) Clean(all bool) {
	for zoneName, zone := range z.Data {
		for recordName, record := range zone.Records {
			if record.IsEmpty() {
				delete(zone.Records, recordName)
			}
		}
		if len(zone.Records) == 0 && all {
			delete(z.Data, zoneName)
		}
	}
}

func (z *ZoneData) AddZone(name string) error {
	if err := utils.RegexHost.Valid(name); err != nil {
		return err
	}
	if _, ok := z.Data[name]; ok {
		return errors.New("zone already exists")
	}
	z.Data[name] = NewZone()
	return nil
}

func (z *ZoneData) RemoveZone(name string) error {
	if _, ok := z.Data[name]; !ok {
		return errors.New("zone does not exist")
	}
	delete(z.Data, name)
	return nil
}

func (r *Record) IsEmpty() bool {
	return len(r.A) == 0 && len(r.AAAA) == 0 && len(r.TXT) == 0 && len(r.CNAME) == 0 &&
		len(r.NS) == 0 && len(r.MX) == 0 && len(r.SRV) == 0 && len(r.CAA) == 0 && r.SOA == nil
}
