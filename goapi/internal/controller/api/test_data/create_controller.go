package test_data

import (
	"clickhouse-api/internal/controller/api"
	"clickhouse-api/internal/model"
	"clickhouse-api/internal/service/queue"
	"log"
	"net/http"
	"time"
)

// swagger:parameters CreateTestDataRequest
type ReqTestDataBody struct {
	//  in: body
	Body TestDataRequestParams `json:"body"`
}

type TestDataRequestParams struct {
	// Data of product
	// in: string
	Data string `json:"data"`
}

type CreateController struct {
	*Writable
}

func NewCreateController(queue *queue.FileQueue, baseController *api.Controller) *CreateController {
	return &CreateController{
		&Writable{
			Controller: baseController,
			Queue:      queue,
		},
	}
}

// swagger:route POST /test_data TestData CreateTestDataRequest
// Create TestData
//
// responses:
//
//	500: Response
//	400: Response
//	201: Response
func (controller *CreateController) Main(writer http.ResponseWriter, request *http.Request) {
	startTime := time.Now()
	url := request.URL.Path
	controller.Counter.WithLabelValues(url).Inc()
	defer controller.Defer(request, url, startTime)

	testData, err := model.CreateTestDataFromRequestBody(request)
	if err != nil {
		controller.HandleError(url, err, writer, http.StatusBadRequest)

		return
	}

	validateMessage, status := testData.IsValid()
	if !status {
		controller.HandleValidationError(url, validateMessage, writer, http.StatusBadRequest)

		return
	}

	testData.Sign = 1
	testData.Version = 1
	testData.CreatedAt = time.Now()

	if err = controller.AddToQueue(testData); err != nil {
		controller.HandleError(url, err, writer, http.StatusInternalServerError)
		log.Fatal()

		return
	}

	controller.MessageResponse(writer, "success", http.StatusCreated)
}
