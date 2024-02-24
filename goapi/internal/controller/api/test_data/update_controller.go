package test_data

import (
	"clickhouse-api/internal/controller/api"
	testDataModel "clickhouse-api/internal/model"
	"clickhouse-api/internal/repository/test_data"
	"clickhouse-api/internal/service/queue"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

// swagger:parameters UpdateTestDataRequest
type UpdateTestDataRequest struct {
	// id of TestData
	// in: query
	Id string `json:"id"`
	//  in: body
	Body TestDataRequestParams `json:"body"`
}

type UpdateController struct {
	Repository *test_data.Repository
	*Writable
}

func NewUpdateController(queue *queue.FileQueue, repository *test_data.Repository, baseController *api.Controller) *UpdateController {
	return &UpdateController{
		Repository: repository,
		Writable: &Writable{
			Controller: baseController,
			Queue:      queue,
		},
	}
}

// swagger:route PUT /test_data/{id} TestData UpdateTestDataRequest
// Create TestData
//
// responses:
//
//	500: Response
//	400: Response
//	201: Response
func (controller *UpdateController) Main(writer http.ResponseWriter, request *http.Request) {
	startTime := time.Now()
	url := request.URL.Path
	controller.Counter.WithLabelValues(url).Inc()
	updateRequest, err := testDataModel.CreateTestDataFromRequestBody(request)
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

	testData.Sign = 1
	testData.Version++

	testData.Data = updateRequest.Data

	if err = controller.AddToQueue(testData); err != nil {
		controller.HandleError(url, err, writer, http.StatusInternalServerError)

		return
	}

	controller.MessageResponse(writer, "success", http.StatusOK)
}
