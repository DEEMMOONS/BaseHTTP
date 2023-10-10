package server


type Server struct {
  cache map[string]models.Order
  db *sql.DB
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
  db, err := sql.Open("postgres", connStr)
  if err != nil {
    return nil, err
  }
  log.Printf("Database is up\n")
  return &Server {
    db: db,
    cache: make(map[string]model.Order),
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
  s.router.HandleFunc("/order/{order_uid}", s.GetByID).Methods("GET")
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
  return fmt.Sprintf("%s:%s", s.config.Address, s.config.Port)
}

func (s *Server) connectToStream() error {
  sc, err := stan.Connect("test-cluster", "subscriber", stan.NatsURL("nats://localhost:4222"))
  if err != nil {
		return err
	}
	sub, err := sc.Subscribe(s.config.SubscribeSubject, s.handleRequest)
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
		log.Printf("Cache updated\n")
		database.AddOrder(s.db, data)
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
  orders := make([]models.Order, 0)
}
