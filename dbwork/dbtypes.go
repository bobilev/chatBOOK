package dbwork

type StateUser struct {
	LastStore int
	LastStep string
}
type Store struct {
	Storeid int
	Text string
	Media int
}
type Step struct {
	StoreId int
	StepID string
	Text string
	Media int
	Answers []Answer
	TypeDoc string
	AccessKey string
}
type Answer struct {
	NextStep string
	Text string
}