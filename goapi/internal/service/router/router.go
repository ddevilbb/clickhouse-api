package router

import (
	"clickhouse-api/internal/controller/api"
	"clickhouse-api/internal/controller/api/test_data"
	"clickhouse-api/internal/middleware/auth"
	testDataRepository "clickhouse-api/internal/repository/test_data"
	"clickhouse-api/internal/service/config"
	"clickhouse-api/internal/service/queue"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
)

type Router struct {
	PingController            *api.PingController
	TestDataCreateController  *test_data.CreateController
	TestDataGetAllController  *test_data.GetListController
	TestDataGetItemController *test_data.GetItemController
	TestDataDeleteController  *test_data.DeleteController
	TestDataUpdateController  *test_data.UpdateController
	MetricsHandler            http.Handler
	authMiddleware            *auth.Middleware
}

func New(
	conf *config.Config,
	createQueue *queue.FileQueue,
	updateQueue *queue.FileQueue,
	deleteQueue *queue.FileQueue,
	repository *testDataRepository.Repository,
	reg *prometheus.Registry,
	metricsHandler http.Handler,
) *Router {
	baseController := api.NewController(reg)
	PingController := api.NewPingController(baseController)

	TestDataCreateController := test_data.NewCreateController(createQueue, baseController)
	TestDataGetAllController := test_data.NewGetAllController(repository, baseController)
	TestDataGetItemController := test_data.NewGetItemController(repository, baseController)
	TestDataDeleteController := test_data.NewDeleteController(deleteQueue, repository, baseController)
	TestDataUpdateController := test_data.NewUpdateController(updateQueue, repository, baseController)

	return &Router{
		PingController:            PingController,
		TestDataCreateController:  TestDataCreateController,
		TestDataGetAllController:  TestDataGetAllController,
		TestDataGetItemController: TestDataGetItemController,
		TestDataDeleteController:  TestDataDeleteController,
		TestDataUpdateController:  TestDataUpdateController,
		MetricsHandler:            metricsHandler,
		authMiddleware:            auth.New(conf, reg),
	}
}

func (router *Router) GetRouter() *mux.Router {
	muxRouter := mux.NewRouter()

	muxRouter.Path("/metrics").Handler(router.MetricsHandler)

	muxRouter.Handle("/docs/swagger.yaml", http.FileServer(http.Dir("./")))
	opts := middleware.SwaggerUIOpts{SpecURL: "/docs/swagger.yaml"}
	sh := middleware.SwaggerUI(opts, nil)
	muxRouter.Handle("/docs", sh)

	apiRouter := muxRouter.PathPrefix("/api").Subrouter()
	apiRouter.Use(router.authMiddleware.Basic)

	apiV1Router := apiRouter.PathPrefix("/v1").Subrouter()
	apiV1Router.HandleFunc("/ping", router.PingController.Main).Methods(http.MethodGet)
	apiV1Router.HandleFunc("/test_data", router.TestDataCreateController.Main).Methods(http.MethodPost)
	apiV1Router.HandleFunc("/test_data", router.TestDataGetAllController.Main).Methods(http.MethodGet)
	apiV1Router.HandleFunc("/test_data/{id:[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$}", router.TestDataGetItemController.Main).Methods(http.MethodGet)
	apiV1Router.HandleFunc("/test_data/{id:[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$}", router.TestDataDeleteController.Main).Methods(http.MethodDelete)
	apiV1Router.HandleFunc("/test_data/{id:[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$}", router.TestDataUpdateController.Main).Methods(http.MethodPut)

	return muxRouter
}
