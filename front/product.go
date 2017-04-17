package front

type Product struct {
	tableName struct{} `sql:"cc_product"`
	ID        uint
	Name      string
	Price     uint
	Contact   bool
	Unit      string
	Detail    string
	Default   bool
}
