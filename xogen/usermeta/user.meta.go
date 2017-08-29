package usermeta

type Field int

const (
	ID Field = iota
	User1ID
	Name
	RefreshToken
)

type SQL int

const (
	SQL_REFORM = SQL(ID & User1ID & Name & RefreshToken)
	SQL_SELF   = SQL(ID & User1ID & Name)
	SQL_VIEW   = SQL(ID & Name)
)

type User_reform struct {
	ID           int    `json:"id"`            // id
	User1ID      int    `json:"user1_id"`      // user1_id
	Name         string `json:"name"`          // name
	RefreshToken []byte `json:"refresh_token"` // refresh_token
}

type User_self struct {
	ID      int    `json:"id"`       // id
	User1ID int    `json:"user1_id"` // user1_id
	Name    string `json:"name"`     // name
}

type User_view struct {
	ID   int    `json:"id"`   // id
	Name string `json:"name"` // name
}

type Query struct {
	Sql    SQL
	Fields string
	Ptrs   []interface{}
	Out    interface{}
}

type queries struct {
	Reform func() Query
	Self   func() Query
	View   func() Query
}

// table xo`query:"reform,self,viewer,admin,sign"`
// user1_id xo`query:"reform,self,-viewer,admin,sign"`
var Queries = queries{
	Reform: func() Query {
		return Query{
			Sql:    SQL_REFORM,
			Fields: "SELECT id, user1_id, name, refresh_token FROM user WHERE id = $1",
			Out:    interface{}(1),
		}
	},
}
