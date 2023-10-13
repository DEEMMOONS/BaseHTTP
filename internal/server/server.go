package server

import (
  "fmt"
  "net/http"
  "encoding/json"
  "log"
  "github.com/go-pg/pg"
  "github.com/nats-io/stan.go"
  "github.com/DEEMMOONS/BaseHTTP/internal/database"
  "github.com/gorilla/mux"
)

type Server struct {
  cache map[string]database.Order
  db *pg.DB
  config *config
  router *mux.Router
	sc     stan.Conn
	sub    stan.Subscription
}

func NewServer(cfgPath string) (*Server, error) {
  config, err := CreateConfig(cfgPath)
  if err != nil {
      return nil, err
  }
  db:= pg.Connect(&pg.Options{
		User:     config.DB.User,
		Password: config.DB.Password,
		Database: config.DB.Database,
	})

  if err != nil {
    return nil, err
  }
  log.Printf("Database is up\n")
  return &Server {
    db: db,
    cache: make(map[string]database.Order),
    config: config,
    router: mux.NewRouter(),
  }, nil
}

func (s *Server) Up() error{
  addr := s.getAddr()
	if err := s.createCache(); err != nil {
		return err
	}
	if err := s.connectToStream(); err != nil {
		return err
	}
  s.router.HandleFunc("/order/{order_uid}", s.getOrder).Methods("GET")
  log.Printf("Server is up on %s\n", addr)
  log.Fatal(http.ListenAndServe(addr, s.router))
  return nil
}

func (s *Server) Down() {
  log.Printf("Server is down\n")
  s.db.Close()
  s.sub.Unsubscribe()
	s.sc.Close()
}

func (s *Server) getAddr() string {
  return fmt.Sprintf("%s:%s", s.config.Host.Address, s.config.Host.Port)
}

func (s *Server) connectToStream() error {
  sc, err := stan.Connect("test-cluster", "subscriber", stan.NatsURL("nats://localhost:4222"))
  if err != nil {
		return err
	}
	sub, err := sc.Subscribe(s.config.Host.SubscribeSubject, s.handleRequest)
	if err != nil {
		return err
	}
	s.sc, s.sub = sc, sub
	return nil
}

func (s *Server) handleRequest(m *stan.Msg) {
	data := database.Order{}
	err := json.Unmarshal(m.Data, &data)
	if err != nil {
		return
	}
	if ok := s.addToCache(data); ok {
		if err := s.addOrder(data); err != nil {
      log.Printf("Order adding error: %w\n", err)
    }
		log.Printf("Data are updated\n")
  }
}

func (s *Server) addToCache(data database.Order) bool {
	_, ok := s.cache[data.OrderUid]
	if ok {
		return false
	}
	s.cache[data.OrderUid] = data
	return true
}

func (s *Server) createCache() error {
	orders := make([]database.Order, 0)
	err := s.db.Model(&orders).Select()
	if err != nil {
		return err
	}
	for _, order := range orders {
		s.cache[order.OrderUid] = order
	}
	return nil
}

func (s *Server) addOrder(data database.Order) error {
  if err := database.AddOrder(s.db, data); err != nil {
    return err
	}
	return nil
}

func (s *Server) getOrder(w http.ResponseWriter, r *http.Request) { 
  vars := mux.Vars(r)
  id := vars["order_uid"]
  data, ok := s.cache[id]
	if !ok {
		http.Error(w, `ID not found`, 400)
    return
	}
	ans, err := json.Marshal(data)
	if err != nil {
    http.Error(w, `Internal server Error`, 500)
		return
	}
  w.Header().Add("Content-Type", "application/json")
  w.WriteHeader(200)
	_, err2 := w.Write(ans)
	if err2 != nil {
		log.Println(err2)
	}
}
