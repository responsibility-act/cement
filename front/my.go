package front

type MyFan struct {
	tableName    struct{} `sql:"cc_user,alias:u"`
	ID           uint
	CreatedAt    int64
	Nickname     string
	HeadImageURL string
	//	User1        uint
	//	Fans         []MyFan
}
