package test_data

import (
	"clickhouse-api/internal/controller/api"
	"clickhouse-api/internal/repository/test_data"
	"clickhouse-api/internal/service/queue"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type DeleteController struct {
	Repository *test_data.Repository
	*Writable
}

func NewDeleteController(queue *queue.FileQueue, repository *test_data.Repository, baseController *api.Controller) *DeleteController {
	return &DeleteController{
		Repository: repository,
		Writable: &Writable{
			Controller: baseController,
			Queue:      queue,
		},
	}
}

// swagger:route DELETE /test_data/{id} TestData DeleteTestData
// Delete TestData by Id
//
// responses:
//
//	500: Response
//	404: Response
//	204:
func (controller *DeleteController) Main(writer http.ResponseWriter, request *http.Request) {
	startTime := time.Now()
	url := request.URL.Path
	controller.Counter.WithLabelValues(url).Inc()
	defer controller.Defer(request, url, startTime)

	id := mux.Vars(request)["id"]
	testData, err := controller.Repository.Find(id)
	if err != nil {
		controller.HandleError(url, err, writer, http.StatusNotFound)

		return
	}

	testData.Sign = -1

	if err = controller.AddToQueue(testData); err != nil {
		controller.HandleError(url, err, writer, http.StatusInternalServerError)

		return
	}

	controller.MessageResponse(writer, "success", http.StatusNoContent)
}
