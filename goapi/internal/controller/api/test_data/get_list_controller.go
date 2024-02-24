package test_data

import (
	"clickhouse-api/internal/controller/api"
	"clickhouse-api/internal/helper"
	testDataModel "clickhouse-api/internal/model"
	"clickhouse-api/internal/repository/test_data"
	"clickhouse-api/pkg/paginator"
	"net/http"
	"time"
)

// swagger:response GetListResponse
type GetListResponse struct {
	// in: body
	Body struct {
		Data       []testDataModel.TestData `json:"data"`
		Message    string                   `json:"message"`
		Pagination struct {
			Page   int `json:"page"`
			Offset int `json:"offset"`
			Limit  int `json:"limit"`
			Total  int `json:"total"`
			Pages  int `json:"pages"`
		} `json:"pagination"`
		Status bool `json:"status"`
	}
}

type GetListController struct {
	*Readable
}

func NewGetAllController(repository *test_data.Repository, baseController *api.Controller) *GetListController {
	return &GetListController{
		&Readable{
			Repository: repository,
			Controller: baseController,
		},
	}
}

// swagger:route GET /test_data TestData GetTestDataList
// Get list of TestData
//
// responses:
//
//	500: Response
//	404: Response
//	200: GetListResponse
func (controller *GetListController) Main(writer http.ResponseWriter, request *http.Request) {
	startTime := time.Now()
	url := request.URL.Path
	controller.Counter.WithLabelValues(url).Inc()
	defer controller.Defer(request, url, startTime)

	queryParams, err := helper.ParseQueryValues(request.URL.Query())
	if err != nil {
		controller.HandleError(url, err, writer, http.StatusBadRequest)

		return
	}

	modelMap, err := helper.ModelToMap(testDataModel.TestData{})
	if err != nil {
		controller.HandleError(url, err, writer, http.StatusBadRequest)

		return
	}

	modelQueryParams, err := helper.GetModelQueryValues(queryParams, modelMap)
	if err != nil {
		controller.HandleError(url, err, writer, http.StatusBadRequest)

		return
	}

	orderQueryParams, err := helper.GetOrderQueryValues(queryParams, modelMap)
	if err != nil {
		controller.HandleError(url, err, writer, http.StatusBadRequest)

		return
	}

	total, err := controller.Repository.CountBy(modelQueryParams)
	if err != nil {
		controller.HandleError(url, err, writer, http.StatusBadRequest)

		return
	}

	insPaginator := paginator.NewPaginator(request.URL.Query(), total)
	if total == 0 {
		controller.PaginationResponse(writer, insPaginator, nil)

		return
	}

	ids, err := controller.Repository.FindIdsBy(insPaginator.Offset, insPaginator.Limit, modelQueryParams)
	if err != nil {
		controller.HandleError(url, err, writer, http.StatusBadRequest)

		return
	}

	list, err := controller.Repository.FindByIds(ids, orderQueryParams)
	if err != nil {
		controller.HandleError(url, err, writer, http.StatusBadRequest)

		return
	}

	controller.PaginationResponse(writer, insPaginator, list)
}
