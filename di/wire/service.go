package wire

// ////////////////////////////////////////////////////////
type Service struct {
	DB IDatabase
}

func (s *Service) DoSomething() string {
	return s.DB.Query()
}

// NewService 是 Service 的提供者函数
func NewService(db IDatabase) *Service {
	return &Service{
		DB: db,
	}
}
