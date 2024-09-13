package ssl

import "crypto/tls"

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
