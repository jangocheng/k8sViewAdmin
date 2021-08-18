package helpers

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"github.com/shenyisyn/goft-gin/goft"
	"io/ioutil"
	"k8sapi/src/models"
	"math/big"
	rd "math/rand"
	"os"
	"time"
)

//解析 k8s 的ca证书，返回对象
func parseK8sCA(CAFile,CAKey string ) (*x509.Certificate, *rsa.PrivateKey){
	caFile, err := ioutil.ReadFile(CAFile)
	goft.Error(err)
	caBlock, _ := pem.Decode(caFile)

	caCert, err := x509.ParseCertificate(caBlock.Bytes) //ca 证书对象
	goft.Error(err)
	//解析私钥
	keyFile, err := ioutil.ReadFile(CAKey)
	goft.Error(err)
	keyBlock, _ := pem.Decode(keyFile)
	caPriKey, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes) //是要对象
	goft.Error(err)
	return caCert,caPriKey
}

const CAFILE="./userConfig/ca.crt"
const CAKEY="./userConfig/ca.key"
func DeleteK8sUser(cn string ){
	err:=os.Remove(fmt.Sprintf("./k8susers/%s.pem",cn))
	goft.Error(err)
	err=os.Remove(fmt.Sprintf("./k8susers/%s_key.pem",cn))
	goft.Error(err)
}
//签发用户证书，其中cn是必须填的
func GenK8sUser(cn,o string) {
	caCert,caPriKey:=parseK8sCA(CAFILE,CAKEY)
	if cn==""{
		goft.Error(fmt.Errorf("CN is required"))
	}
	certTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(rd.Int63()), //证书序列号
		Subject: pkix.Name{
			Country:            []string{"CN"},
			Organization:       []string{o},
			//OrganizationalUnit: []string{"可填课不填"},
			Province:           []string{cn},
			CommonName:         cn,//CN
			Locality:           []string{"beijing"},
		},
		NotBefore:             time.Now(),//证书有效期开始时间
		NotAfter:              time.Now().AddDate(1, 0, 0),//证书有效期
		BasicConstraintsValid: true, //基本的有效性约束
		IsCA:                  false,   //是否是根证书
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}, //证书用途(客户端认证，数据加密)
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageDataEncipherment,
		EmailAddresses:        []string{"UserAccount@jtthink.com"},
	}

	//生成公私钥--秘钥对
	priKey, err := rsa.GenerateKey(rand.Reader, 2048)
	goft.Error(err)
	//创建证书 对象
	clientCert, err := x509.CreateCertificate(rand.Reader, certTemplate, caCert, &priKey.PublicKey, caPriKey)
	goft.Error(err)

	//编码证书文件和私钥文件
	clientCertPem := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: clientCert,
	}
	clientCertFile, err := os.OpenFile(fmt.Sprintf("./k8susers/%s.pem",cn), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	goft.Error(err)
	defer clientCertFile.Close()
	err = pem.Encode(clientCertFile, clientCertPem)
	goft.Error(err)

	buf := x509.MarshalPKCS1PrivateKey(priKey)
	keyPem := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: buf,
	}
	clientKeyFile, _ := os.OpenFile( fmt.Sprintf("./k8susers/%s_key.pem",cn), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	defer clientKeyFile.Close()
	err = pem.Encode(clientKeyFile, keyPem)
	goft.Error(err)
}


func getCertType( alg x509.PublicKeyAlgorithm) string {
	switch alg{
	case x509.RSA:
			return "RSA"
	case x509.DSA:
		return "DSA"
	case x509.ECDSA:
		return "ECDSA"
	case x509.UnknownPublicKeyAlgorithm:
		return "Unknow"
	}
	return "Unknow"
}

//解析证书
func ParseCert(crt []byte) *models.CertModel{
	var cert tls.Certificate
	//解码证书
	certBlock, restPEMBlock := pem.Decode(crt)
	if certBlock == nil {
		return nil
	}
	cert.Certificate = append(cert.Certificate, certBlock.Bytes)
	//处理证书链
	certBlockChain, _ := pem.Decode(restPEMBlock)
	if certBlockChain != nil {
		cert.Certificate = append(cert.Certificate, certBlockChain.Bytes)
	}

	//解析证书
	x509Cert, err := x509.ParseCertificate(certBlock.Bytes)

	if err != nil {
		return nil
	} else {
		return &models.CertModel{
			CN:x509Cert.Subject.CommonName,
			Issuer:x509Cert.Issuer.CommonName,
			Algorithm:getCertType(x509Cert.PublicKeyAlgorithm),
			BeginTime:x509Cert.NotBefore.Format("2006-01-02 15:03:04"),
			EndTime:x509Cert.NotAfter.Format("2006-01-02 15:03:04"),
		}
	}
}