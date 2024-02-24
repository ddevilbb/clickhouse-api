package api

import (
	"net/http"
	"time"
)

type PingController struct {
	*Controller
}

func NewPingController(baseController *Controller) *PingController {
	return &PingController{
		Controller: baseController,
	}
}

func (controller *PingController) Main(writer http.ResponseWriter, request *http.Request) {
	startTime := time.Now()
	url := request.URL.Path
	controller.Counter.WithLabelValues(url).Inc()
	defer controller.Defer(request, url, startTime)

	controller.MessageResponse(writer, "pong", http.StatusOK)
}
