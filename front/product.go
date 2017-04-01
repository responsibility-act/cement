package front

type Product struct {
	tableName  struct{} `sql:"cc_product"`
	ID         uint
	Name       string
	Img        string
	ImgSpecial string
	Intro      string
	Detail     string
	Saled      uint
	CreatedAt  int64
	SaledAt    int64
	ShelfOffAt int64
	Color      string
}

type Sku struct {
	tableName   struct{} `sql:"cc_sku"`
	ID          uint
	Stock       uint
	Img         string
	Price       uint
	MarketPrice uint
	Freight     uint
	ProductID   uint
	Color       string
}

type SkuAttr struct {
	tableName struct{} `sql:"cc_sku_attr"`
	ID        uint
	SkuID     uint
	AttrID    uint
}
