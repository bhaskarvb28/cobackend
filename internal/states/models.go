package states 

type State struct {
	ID	 string 
	Name string
}

type GetStatesQueryParams struct {
	Search string
	Order   string
}