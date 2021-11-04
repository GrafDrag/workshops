package model

type SearchEvent struct {
	UserID   int
	Title    string
	Timezone string
	DateFrom string
	DateTo   string
	TimeFrom string
	TimeTo   string
}
