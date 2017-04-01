package dbs

func (dbs *DBService) SiteExist(domain string) (ok bool) {
	dbs.sitesMu.RLock()
	defer dbs.sitesMu.RUnlock()
	_, ok = dbs.sites[domain]
	return
}
