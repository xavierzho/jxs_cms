package global

import (
	"path/filepath"

	"data_backend/pkg/util"

	"github.com/oschwald/geoip2-golang"
)

var (
	IPDB util.IPDB
)

func SetupIpDB() (err error) {
	ipdb, err := geoip2.Open(filepath.Join(StoragePath, "data", "GeoLite2-City.mmdb"))
	if err != nil {
		return err
	}

	IPDB = util.IPDB{Reader: ipdb}

	return err
}
