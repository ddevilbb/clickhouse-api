package test_data

import (
	"clickhouse-api/internal/controller/api"
	testDataModel "clickhouse-api/internal/model"
	"clickhouse-api/internal/repository/test_data"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

// swagger:response GetItemResponse
type GetItemResponse struct {
	// in: body
	TestData testDataModel.TestData
}

type GetItemController struct {
	*Readable
}

func NewGetItemController(repository *test_data.Repository, baseController *api.Controller) *GetItemController {
	return &GetItemController{
		&Readable{
			Repository: repository,
			Controller: baseController,
		},
	}
}

// @formatter:off

// swagger:route GET /test_data/{id} TestData GetTestDataById
// Get TestData by id
//
// Parameters:
//   - name: id
//     in: query
//     required: true
//     type: string
//     format: string
//
// responses:
//
//	500: Response
//	404: Response
//	200: GetItemResponse
func (controller *GetItemController) Main(writer http.ResponseWriter, request *http.Request) {
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

	controller.DataResponse(writer, "data", testData, http.StatusOK)
}
