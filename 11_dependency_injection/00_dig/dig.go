package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/dig"
)

/*
场景：
	审查一个HTTP服务器的代码，当客户端发出GET请求时，它会提供JSON响应/people

*/

// 1.返回的Person
type Person struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// 2. 服务器配置
//Enabled告诉我们我们的应用程序是否应该返回实际数 DatabasePath告诉我们数据库在哪里（我们正在使用sqlite）。Port告诉我们将运行我们的服务器的端口
type Config struct {
	Enabled      bool
	DatabasePath string
	Port         string
}

func NewConfig() *Config {
	return &Config{
		Enabled:      true,
		DatabasePath: "11_dependency_injection/00_dig/example.db",
		Port:         "8000",
	}
}

// 3. 打开数据库链接
func ConnectDatabase(config *Config) (*sql.DB, error) {
	return sql.Open("sqlite3", config.DatabasePath)
}

// 4. PersonRepository仓库：负责从我们的数据库中提取人员并将这些数据库结果反序列化为合适的Person结构
type PersonRepository struct {
	database *sql.DB //管理数据库连接
}

// PersonRepository需要建立数据库连接。它公开了一个单独的函数FindAll，它使用我们的数据库连接返回一个Person表示数据库中数据的结构列表
func (repository *PersonRepository) FindAll() []*Person {
	rows, _ := repository.database.Query(
		`SELECT id, name, age FROM people;`,
	)
	defer rows.Close()

	people := []*Person{}

	for rows.Next() {
		var (
			id   int
			name string
			age  int
		)

		rows.Scan(&id, &name, &age)

		people = append(people, &Person{
			Id:   id,
			Name: name,
			Age:  age,
		})
	}

	return people
}

// 创建PersonRepository仓库
func NewPersonRepository(database *sql.DB) *PersonRepository {
	return &PersonRepository{database: database}
}

// 5. 为了在我们的HTTP服务器和PersonRepository我们之间提供一个图层，我们将创建一个PersonService
// 我们PersonService依赖于Config和PersonRepository。它公开了一个被称为“ FindAll有条件地调用PersonRepository应用程序是否被启用” 的函数
type PersonService struct {
	config     *Config
	repository *PersonRepository
}

func (service *PersonService) FindAll() []*Person {
	if service.config.Enabled {
		return service.repository.FindAll()
	}

	return []*Person{}
}

func NewPersonService(config *Config, repository *PersonRepository) *PersonService {
	return &PersonService{config: config, repository: repository}
}

//6. 这是负责运行一个HTTP服务器并委托给我们的合适的请求PersonService
type Server struct {
	config        *Config
	personService *PersonService
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/people", s.peopleHandler)

	return mux
}

func (s *Server) Run() {
	httpServer := &http.Server{
		Addr:    ":" + s.config.Port,
		Handler: s.Handler(),
	}

	httpServer.ListenAndServe()
}

func (s *Server) peopleHandler(w http.ResponseWriter, r *http.Request) {
	people := s.personService.FindAll()
	bytes, _ := json.Marshal(people)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func NewServer(config *Config, service *PersonService) *Server {
	return &Server{
		config:        config,
		personService: service,
	}
}

// 修改前：可怕的初始化
//func main() {
//	config := NewConfig()
//
//	db, err := ConnectDatabase(config)
//
//	if err != nil {
//		panic(err)
//	}
//
//	personRepository := NewPersonRepository(db)
//
//	personService := NewPersonService(config, personRepository)
//
//	server := NewServer(config, personService)
//
//	server.Run()
//}

//2。修改后
func BuildContainer() *dig.Container {
	container := dig.New()

	container.Provide(NewConfig)
	container.Provide(ConnectDatabase)
	container.Provide(NewPersonRepository)
	container.Provide(NewPersonService)
	container.Provide(NewServer)

	return container
}

func main() {
	container := BuildContainer()

	err := container.Invoke(func(server *Server) {
		server.Run()
	})

	if err != nil {
		panic(err)
	}
}

/*
好处
	创建我们的组件与创建它们的依赖关系的解耦。比如说，我们PersonRepository现在需要访问Config。
	我们所要做的就是改变我们的NewPersonRepository构造函数以包含Config作为参数。我们的代码中没有其他的变化。
	缺乏全局状态，缺少调用init（需要时懒惰地创建依赖关系，只创建一次，无需易错的init设置），并且易于对各个组件进行测试。
	想象一下，在测试中创建容器并要求完全构建的对象进行测试。或者，使用所有依赖关系的模拟实现来创建一个对象。所有这些在DI方法中都更容易。

*/
