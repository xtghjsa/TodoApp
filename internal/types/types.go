package types

type Task struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Status      string `json:"status"`
}
type ShowTasks struct {
	Tasks []Task
}

type TaskAdd struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
}

type TaskId struct {
	Id string `json:"id"`
}
type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
