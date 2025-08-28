package ssl

import (
	"fmt"
	"log"
	"testing"
)

func TestGenerateCACertificatePEM(t *testing.T) {
	// 用于生成自签名CA（证书颁发机构）证书
	certPEM, keyPEM, err := GenerateCACertificatePEM()
	if err != nil {
		log.Fatalf("创建证书失败: %v", err)
	}

	fmt.Println("证书 PEM:")
	fmt.Println(string(certPEM))

	fmt.Println("私钥 PEM:")
	fmt.Println(string(keyPEM))

	// 证书 PEM:
	// -----BEGIN CERTIFICATE-----
	// MIIDBzCCAe+gAwIBAgICB+QwDQYJKoZIhvcNAQELBQAwFTETMBEGA1UEChMKRXhh
	// bXBsZSBDQTAeFw0yNTA2MTQxMzMyMDZaFw0zNTA2MTQxMzMyMDZaMBUxEzARBgNV
	// BAoTCkV4YW1wbGUgQ0EwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDC
	// 0oH62oSg1AI64mdxEXfqYnEWfM2GycUtj65eGRK3lwGwEKbJfk77OLguKW5joZ9c
	// FCFXQ8E/KwPKoY57NmlgjpJbhxfKxUqUs4dL7PNCLC0gdEDQoR0q822JgT2tO8sf
	// nG1LTClMqfOSenNLRQ6kH1JdcDh6wsdPuEX7VNHZhNA55+tY/HrQDD5S73pugT3C
	// ydCoaksGR8WvnZDQXBFKxWBqc31khDc2hzggBjQuO/ql5MxjCl0BdlJ8ZHo6mT2N
	// S2vVy0Ep9dVJyzo5cz7UISbF+Dlfpk3aSQgcw2V7Qs7M48f4exmoWqUlyeT4jULh
	// OPL9muT3/SKZbKK5L9fFAgMBAAGjYTBfMA4GA1UdDwEB/wQEAwIChDAdBgNVHSUE
	// FjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4E
	// FgQU1mUxfE9SkoY0WwtK6eW6VW8jNsIwDQYJKoZIhvcNAQELBQADggEBALVDb8IT
	// nA2JqoCSN1epItxwE4gLgpM5mhFs+k09HanqfV2aA6aw1p5iQ/akVoVaZFVAAJKw
	// majDya2s+DtGVN/bpn334B/vCioVwLw8gClX+1qTxDsYwm533P5JURebQYBfcoRY
	// U3F8ftDwCSZu4eRxbBACwy/mioRU1JJrs8U2r9RqWzZFueoxXRvy5WoLErkz09d9
	// ya0SlPYPCfzSFZgZguxdF5cN444/IW7+mnTLVgK1D4wI4ZJ8ISEC6xpGnGlhBHzV
	// MnlRU/biCix/bxonMJtwUsGAVtA0nSdgymJnvn08FVpjhaoeaBMRlwMP3wnZymqJ
	// BN3p3WsqQ4lnj9c=
	// -----END CERTIFICATE-----

	// 私钥 PEM:
	// -----BEGIN RSA PRIVATE KEY-----
	// MIIEogIBAAKCAQEAwtKB+tqEoNQCOuJncRF36mJxFnzNhsnFLY+uXhkSt5cBsBCm
	// yX5O+zi4LiluY6GfXBQhV0PBPysDyqGOezZpYI6SW4cXysVKlLOHS+zzQiwtIHRA
	// 0KEdKvNtiYE9rTvLH5xtS0wpTKnzknpzS0UOpB9SXXA4esLHT7hF+1TR2YTQOefr
	// WPx60Aw+Uu96boE9wsnQqGpLBkfFr52Q0FwRSsVganN9ZIQ3Noc4IAY0Ljv6peTM
	// YwpdAXZSfGR6Opk9jUtr1ctBKfXVScs6OXM+1CEmxfg5X6ZN2kkIHMNle0LOzOPH
	// +HsZqFqlJcnk+I1C4Tjy/Zrk9/0imWyiuS/XxQIDAQABAoIBACj0SRwjbwGDB9P6
	// j42ygx+LoaO9SRQ2WqOdAmXoBen/jbyGB5WwXmiLsBYOIhVCHsajua2HQfqmL3Yb
	// d6D5m3XPir9AWxVGW4r+YWjpzupAcJ0TqyNgVwoWIZboCv/dY7IJt2T+hekGifwn
	// DxEJ243PQsh/JHRT+UOOJHH0zudXWKR70EQYw1DD+IWxj+utCbcBMIT92WOzMQpR
	// NoxuK2EWsEuoKcKF4N8zD+fPr4WLAxmSnmoWLRi6P9gmVp9E4d7tgd8AaRxSzWBJ
	// tlLEuIm6Vk+x2Yuo2YJY0kF4XEtZlIk/OdHEnVt5Q9LHm32MbatpxOHmHT1dsXfy
	// ZBJ/DRECgYEAzq0ldsYqhxtfP7VFaQUMRzCiEs0Wnmfyuf7wk/dYLMy+PyvLiX/k
	// cLXfVIxy8Ov6NfT/CyquXfja7PTszGWJkZ6KNXLvyLTcc5X7XP315GzH0nQlM1Fe
	// tAPndgNlENS9mRIIpE7L7ghvkt7uOfMbGMfoCDqaTH+n61ZL8EZ/xxUCgYEA8VEj
	// 1HHMoA23LUW1mJZPXn15p8mHF82p7bZnUR1T90+kJUG2rQW65wiaakCI/D5LLW6D
	// hWaZXCeU2+ma58+UVd5sY6NMmIoiuFW0zoLZNl1kK9e8xGfBDfnioCo4YGizRH8d
	// U8f3A/gPUcVQkmbeT/2Xw8v2rMMSMFxs1Pm6+fECgYArb/ylv/SEPN6B90lFT0hL
	// Vg9aQDx2woYjTU+m6Z9gmw+JG11F4tlSTwdHL9WgiRgnavyHjkrjeUAZ+UgjlVua
	// fWWy4hs/ZbPHn0gbPU0G204MD1kaNgnfb8qf5QrCxNOsbjvevKjjuGYqyivrhgq1
	// 5J4BzL9NQK88KQEA2PBWGQKBgBWy6xreRL0bnp4Gh6a51Vc0xyysNWaRirciULX9
	// giBZ2/OxrgBu5HiD0Ia/WNH9s/rY1iC3shCUSpFftxsjEj6KaoqnE2sf+LFEm6Z6
	// I5f829YJZyLuBXEBSDyIr1sT7xK4r2VqNK75rj73FCCl+VWOAwiLHZo5TDhnBy47
	// anGxAoGAOZ+q+JP/a+yV3Oywc5liFgaqRRGk/6NfTnZ5/ZTQ/KMk5YqgstIA9wMM
	// AUT0pMtT52svH/Ki1nik/b9ZYlWsLpgrP8jQVmtFzkwYTO3k64py89PEJzO/fY+n
	// 6Rz7+aJ11l3q1niyu5jpfNhL7J2rGxOxtlxf4sohk3pWySTcPhw=
	// -----END RSA PRIVATE KEY-----
}
