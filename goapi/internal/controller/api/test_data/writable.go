package test_data

import (
	"clickhouse-api/internal/controller/api"
	"clickhouse-api/internal/service/queue"
	"encoding/json"
)

type Writable struct {
	*api.Controller
	Queue *queue.FileQueue
}

func (controller *Writable) AddToQueue(model interface{}) error {
	modelByte, err := json.Marshal(model)
	if err != nil {
		return err
	}

	if err = controller.Queue.Add(string(modelByte)); err != nil {
		return err
	}

	return nil
}
