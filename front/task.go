package front

type Task struct {
	tableName struct{} `sql:"cc_task,alias:t"`
	ID        uint
	Name      string
	Cash      uint
	Repeat    uint
	NotBefore int64
	ExpiresAt int64
	Active    bool
}
