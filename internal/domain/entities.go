package domain

// Project структура под list
type Project struct {
	Id    string `json:"uuid"`
	Title string `json:"title"`
}

// ProjectDetail структура под detail
type ProjectDetail struct {
	Id    string `json:"uuid"`
	Title string `json:"title"`
	// TODO нужно будет добавить версионирование
}

// ProjectCreateOrUpdate структура под create и update
type ProjectCreateOrUpdate struct {
	Title string `json:"title"`
}
