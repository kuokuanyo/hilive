package service

var services = make(Generators)

// Generator is func() (Service, error)
type Generator func() (Service, error)

// Generators is map[string]Generator, and Generator is func() (Service, error)
type Generators map[string]Generator

// List is map[string]Service, and Service is interface(method: Name)
type List map[string]Service

// Service is interface
// db.Connection(interface)也具有Service(interface)，都具有Name方法
type Service interface {
	Name() string
}

// Register 將參數將入services(map[string]Generator)中
func Register(k string, gen Generator) {
	if _, ok := services[k]; ok {
		panic("service has been registered")
	}
	services[k] = gen
}

// GetServices 初始化List(map[string]Service)
func GetServices() List {
	var (
		l   = make(List)
		err error
	)
	for k, gen := range services {
		if l[k], err = gen(); err != nil {
			panic("初始化Service失敗")
		}
	}
	return l
}

// Get 取得匹配參數的Service
func (g List) Get(k string) Service {
	if s, ok := g[k]; ok {
		return s
	}
	panic("找不到匹配的Service")
}

// Add 將參數加入至List(map[string]Service)
func (g List) Add(k string, service Service) {
	if _, ok := g[k]; ok {
		panic("service exist")
	}
	g[k] = service
}