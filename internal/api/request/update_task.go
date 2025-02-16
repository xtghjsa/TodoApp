package request

type UpdateTask struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
	ID          int    `json:"id"`
}
