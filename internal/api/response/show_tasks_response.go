package response

type ShowTasksItem struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Status      string `json:"status"`
}

type ShowTasks struct {
	Tasks []ShowTasksItem
}
