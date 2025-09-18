package util

import (
	"fmt"
	"net"

	"github.com/oschwald/geoip2-golang"
)

type Address struct {
	City    string
	Country string
}

type IPDB struct {
	*geoip2.Reader
}

func (d IPDB) CheckLegalIP(ip string) string {
	record, err := d.City(net.ParseIP(ip))
	if err != nil {
		return fmt.Sprintf("根据IP查询区域报错: ip=%v, err=%v", ip, err)
	}

	return record.City.Names["en"]
}

func (d IPDB) ParseIPCity(ip string) Address {
	if string(ip) == "" {
		return Address{}
	}

	record, err := d.City(net.ParseIP(string(ip)))
	if err != nil {
		return Address{}
	}

	return Address{
		City:    record.City.Names["en"],
		Country: record.Country.Names["en"],
	}
}
