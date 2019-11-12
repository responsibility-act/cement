package dbs

import (
	"sync"

	"github.com/empirefox/cement/config"
	"github.com/empirefox/cement/front"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"gopkg.in/pg.v5"
)

type DBS interface {
	SiteExist(domain string) bool
}

type DBService struct {
	config config.Postgres
	db     *pg.DB
	logger zap.Logger

	sites   map[string]front.SiteBase
	sitesMu sync.RWMutex
}

func NewDBService(config *config.Config) (*DBService, error) {
	opts := new(pg.Options)
	copier.Copy(opts, &config.Postgres)
	dbs := &DBService{
		config: config.Postgres,
		db:     pg.Connect(opts),
		logger: config.Logger,
		sites:  make(map[string]front.SiteBase),
	}

	var sites []front.SiteBase
	if err := dbs.db.Model(&sites).Select(); err != nil {
		return nil, err
	}

	for _, site := range sites {
		dbs.sites[site.Domain] = site
	}
	return dbs, nil
}

func (dbs *DBService) Ok() {}
