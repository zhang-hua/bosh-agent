package monit

type Status interface {
	GetIncarnation() (int, error)
	ServicesInGroup(name string) (services []Service)
	ServiceStatsInGroup(name string) (serviceStats map[string]interface{})
}

type Service struct {
	Monitored bool
	Status    string
}
