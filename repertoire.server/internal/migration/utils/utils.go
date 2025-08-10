package utils

import (
	"fmt"
	"repertoire/server/data/search"
	"time"

	"github.com/meilisearch/meilisearch-go"
)

type MigrationStatus struct {
	Id        string    `json:"id"`
	IsApplied bool      `json:"is_applied"`
	Timestamp time.Time `json:"timestamp"`
}

func HasMigrationAlreadyBeenApplied(client search.MeiliClient, uid string) bool {
	_, err := client.GetIndex("migration_version")
	if err != nil {
		taskInfo, err := client.CreateIndex(&meilisearch.IndexConfig{
			Uid:        "migration_version",
			PrimaryKey: "id",
		})
		if err != nil {
			panic(err)
		}
		for {
			task, _ := client.GetTask(taskInfo.TaskUID)
			if task.Status != meilisearch.TaskStatusEnqueued && task.Status != meilisearch.TaskStatusProcessing {
				break
			}
		}
	}

	var documentResults meilisearch.DocumentsResult
	err = client.Index("migration_version").GetDocuments(&meilisearch.DocumentsQuery{}, &documentResults)
	if err != nil {
		panic(err)
	}

	for _, result := range documentResults.Results {
		if result["id"] == uid && result["is_applied"].(bool) {
			return true
		}
	}

	return false
}

func SaveMigrationStatus(client search.MeiliClient, uid string, name string) {
	status := &MigrationStatus{
		Id:        uid,
		IsApplied: true,
		Timestamp: time.Now(),
	}
	addTask, err := client.Index("migration_version").AddDocuments([]any{status})
	if err != nil {
		panic(err)
	}
	for {
		task, err := client.GetTask(addTask.TaskUID)
		if err != nil {
			panic(err)
		}
		if task.Status == meilisearch.TaskStatusEnqueued || task.Status == meilisearch.TaskStatusProcessing {
			break
		}
	}
	fmt.Println(time.Now().Format("2006/01/02 15:01:05") + " OK " + uid + "_" + name + " imported!")
}
