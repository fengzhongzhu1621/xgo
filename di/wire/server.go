package wire

// Server 定义HTTP服务器结构体
type Server struct {
	DB *Database
}

// NewServer 创建一个新的服务器实例
func NewServer(db *Database) *Server {
	return &Server{DB: db}
}
