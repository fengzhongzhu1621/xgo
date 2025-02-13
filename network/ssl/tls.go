package ssl

import "crypto/tls"

func ConfigureTLS() *tls.Config {
	return &tls.Config{
		// 设置 TLS 的最低版本为 TLS 1.2。这确保了使用较新的加密协议，增强了安全性，
		// 避免了较旧且不安全的协议（如 SSLv3、TLS 1.0 和 TLS 1.1）。
		MinVersion: tls.VersionTLS12,
		// 推荐启用 TLS 1.3 以获得更好的安全性和性能
		MaxVersion: tls.VersionTLS13,
		// 指示服务器优先选择其支持的加密套件列表中的套件。这有助于服务器强制使用特定的加密算法，提升安全性和性能。
		// 已被弃用，取而代之的是通过明确指定 CipherSuites 列表来控制服务器使用的加密套件。
		// 需要手动指定希望服务器支持的加密套件，而不是依赖于优先级设置。
		// PreferServerCipherSuites: true,
		// 明确指定服务器支持的加密套件
		// 定义了一组服务器支持的加密套件。这些套件结合了密钥交换机制、对称加密算法和消息认证码（MAC），用于保护数据传输的安全性。
		// 这些套件使用了 ECDHE（椭圆曲线 Diffie-Hellman 密钥交换）进行密钥协商，提供了前向保密性（Forward Secrecy），
		// 并且结合了 AES-GCM 或 ChaCha20-Poly1305 进行对称加密，确保数据的机密性和完整性。
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			// 根据需要添加更多加密套件
		},

		// 要求客户端提供证书。
		// 如果需要双向认证，可以配置 ClientAuth
		// ClientAuth: tls.RequireAndVerifyClientCert,

		// 其他 TLS 配置选项
		// SessionTicketsDisabled: false,
		// SessionTicketKey:       sessionTicketKey,
	}
}

// ClientTslConfNoVerity InsecureSkipVerify是一个布尔类型的字段，用于控制是否跳过TLS服务器证书验证。
// 当设置为true时，客户端将不会验证服务器的TLS证书是否有效，这可能会导致安全风险，因为中间人攻击者可能会伪造证书。
func ClientTslConfNoVerity() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: true,
	}
}

// ClientTslConfVerityServer 用于配置 TLS 客户端的 SSL/TLS 设置，以便验证服务器的证书。
func ClientTslConfVerityServer(caFile string) (*tls.Config, error) {
	// 读取 CA 证书文件，返回证书池
	certPool, err := LoadX509Cert(caFile)
	if err != nil {
		return nil, err
	}

	// 配置 TLS 客户端
	tlsConfig := &tls.Config{
		RootCAs: certPool,
	}

	return tlsConfig, nil
}

// ClientTslConfVerity 用于创建一个 TLS 客户端的 SSL/TLS 配置，包括加载 CA 证书、客户端证书和私钥
func ClientTslConfVerity(caFile, certFile, keyFile, passwd string) (*tls.Config, error) {
	// 读取 CA 证书文件
	certPool, err := LoadX509Cert(caFile)
	if err != nil {
		return nil, err
	}
	// 从指定的客户端证书和私钥文件中加载 TLS 证书，并在需要时解密私钥，返回一个 *tls.Certificate 对象
	cert, err := LoadX509Certificates(certFile, keyFile, passwd)
	if err != nil {
		return nil, err
	}

	// 配置 TLS 客户端
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,                     // 跳过服务器证书的验证
		RootCAs:            certPool,                 // CA 证书池，用于验证服务器证书
		Certificates:       []tls.Certificate{*cert}, // 客户端证书和私钥
	}

	return tlsConfig, nil
}

// ServerTslConf 服务端Tsl配置
func ServerTslConf(caFile, certFile, keyFile, passwd string) (*tls.Config, error) {
	if "" == caFile {
		return ServerTslConfVerity(certFile, keyFile, passwd)
	}

	return ServerTslConfVerityClient(caFile, certFile, keyFile, passwd)
}

// ServerTslConfVerity 创建一个 TLS 服务器的 SSL/TLS 配置，包括加载服务器证书和私钥。
func ServerTslConfVerity(certFile, keyFile, passwd string) (*tls.Config, error) {
	cert, err := LoadX509Certificates(certFile, keyFile, passwd)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{*cert}, // 服务器证书和私钥
	}

	return tlsConfig, nil
}

// ServerTslConfVerityClient 创建一个 TLS 服务器的 SSL/TLS 配置，包括加载 CA 证书、服务器证书和私钥，并要求客户端提供并验证其证书。
// 这有助于确保只有拥有有效客户端证书的用户才能访问服务器。
func ServerTslConfVerityClient(caFile, certFile, keyFile, passwd string) (*tls.Config, error) {
	// 读取 CA 证书文件
	caPool, err := LoadX509Cert(caFile)
	if err != nil {
		return nil, err
	}

	// 从指定的证书和私钥文件中加载 TLS 证书，并在需要时解密私钥
	cert, err := LoadX509Certificates(certFile, keyFile, passwd)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		ClientCAs:    caPool,                         // CA 证书池，用于验证客户端证书
		Certificates: []tls.Certificate{*cert},       // 服务器证书和私钥
		ClientAuth:   tls.RequireAndVerifyClientCert, // 要求客户端提供证书并验证其有效性
	}

	return tlsConfig, nil
}
