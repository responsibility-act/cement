package front

type ProductAttr struct {
	tableName struct{} `sql:"cc_product_attr"`
	ID        uint
	Value     string
	GroupID   uint
	Pos       int64
}

type ProductAttrGroup struct {
	tableName struct{} `sql:"cc_product_attr_group"`
	ID        uint
	Name      string
	Pos       int64
}
