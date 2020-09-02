package storage

func getUpHost(config *Config, ak, bucket string) (upHost string, err error) {
	var zone *Zone
	if config.Zone != nil {
		zone = config.Zone
	} else if zone, err = GetZone(ak, bucket); err != nil {
		return
	}

	scheme := "http://"
	if config.UseHTTPS {
		scheme = "https://"
	}

	host := zone.SrcUpHosts[0]
	if config.UseCdnDomains {
		host = zone.CdnUpHosts[0]
	}

	upHost = scheme + host
	return
}
