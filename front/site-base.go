package front

import (
	"time"

	"gopkg.in/pg.v5/orm"
)

type SiteBase struct {
	tableName struct{} `sql:"cc_site_base,alias:sb"`
	ID        int64
	UserID    int64
	Domain    string
	MainCdn   string
	CreatedAt int64
	Expires   int64
}

func (b *SiteBase) BeforeInsert(db orm.DB) error {
	b.CreatedAt = time.Now().Unix()
	return nil
}
