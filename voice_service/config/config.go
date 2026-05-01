package config

var (
	QueuesNames = map[string]QueueConf{
		"A_QUEUE": QueueConf{
			Name:     "ocr_response_q",
			Consumer: true,
			Dlq:      QueueInternals{Name: "ocr_response_dlq_q", Ttl: 0},
			Retry:    QueueInternals{Name: "ocr_response_retry_q", Ttl: 5000},
		},
		"B_QUEUE": QueueConf{Name: "voice_output_q", Consumer: false},
	}

	RabbitMqUrl = "amqp://guest:guest@rabbitmq:5672/"
	RetryCount  = 5
)

type QueueConf struct {
	Name     string
	Consumer bool
	Dlq      QueueInternals
	Retry    QueueInternals
}

type QueueInternals struct {
	Name string
	Ttl  int32 // if 0 then none
}
