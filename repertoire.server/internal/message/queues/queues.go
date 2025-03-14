package queues

type Queue string

const (
	MainQueue    Queue = "main_queue"
	SearchQueue  Queue = "search_queue"
	StorageQueue Queue = "storage_queue"
)
