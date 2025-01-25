package CertificateManager

import (
	"encoding/json"
	"fmt"
	"strings"

	"parkjunwoo.com/ncloud-sdk-go/services"
	"parkjunwoo.com/ncloud-sdk-go/services/Networking/GlobalDNS"
)

type Certificates struct {
	ReturnCode         string        `json:"returnCode"`
	ReturnMessage      string        `json:"returnMessage"`
	TotalRows          int           `json:"totalRows"`
	SslCertificateList []Certificate `json:"sslCertificateList"`
}

type Certificate struct {
	CertificateNo         int    `json:"certificateNo"`
	CertificateType       string `json:"certificateType"`
	CertificateName       string `json:"certificateName"`
	MemberNo              string `json:"memberNo"`
	DnInfo                string `json:"dnInfo"`
	DomainAddress         string `json:"domainAddress"`
	RegDate               string `json:"regDate"`
	UpdateDate            string `json:"updateDate"`
	IssueDate             string `json:"issueDate"`
	ValidStartDate        string `json:"validStartDate"`
	ValidEndDate          string `json:"validEndDate"`
	StatusCode            string `json:"statusCode"`
	StatusName            string `json:"statusName"`
	ExternalYn            string `json:"externalYn"`
	DomainCode            string `json:"domainCode"`
	CaInfo                string `json:"caInfo"`
	CertSerialNumber      string `json:"certSerialNumber"`
	CertPublicKeyInfo     string `json:"certPublicKeyInfo"`
	CertSignAlgorithmName string `json:"certSignAlgorithmName"`
}

type CertificateReturn struct {
	ReturnCode         string        `json:"returnCode"`
	ReturnMessage      string        `json:"returnMessage"`
	TotalRows          int           `json:"totalRows"`
	SslCertificateList []Certificate `json:"sslCertificateList"`
}

type ExternalCertificate struct {
	CertificateName      string `json:"certificateName"`
	PrivateKey           string `json:"privateKey"`
	PublicKeyCertificate string `json:"publicKeyCertificate"`
	CertificateChain     string `json:"certificateChain"`
}

func cleanPEM(pem string) string {
	return strings.TrimSpace(strings.ReplaceAll(pem, "\r\n", "\n"))
}

func GetCertificates(access *services.Access) (*Certificates, error) {
	endpoint := "https://certificatemanager.apigw.ntruss.com"
	url := "/api/v1/certificates"
	resp, err := services.Request(access, "GET", endpoint, url, nil)
	if err != nil {
		return nil, err
	}

	certificates := Certificates{}
	if err := json.NewDecoder(resp.Body).Decode(&certificates); err != nil {
		return nil, fmt.Errorf("failed to GetCertificates JSON: %v", err)
	}
	defer resp.Body.Close()

	if certificates.ReturnCode != "0" {
		return nil, fmt.Errorf("failed to GetCertificates: %v", certificates)
	}

	return &certificates, nil
}

func CreateExternalCertificate(access *services.Access, certificateName string, privateKey string, certificateBody string, certificateChain string, rootCA string) (*Certificate, error) {
	endpoint := "https://certificatemanager.apigw.ntruss.com"
	url := "/api/v1/certificate/withExternal"

	rootDomain, subdomain, err := GlobalDNS.GetDomainNames(certificateName)
	if err != nil {
		return nil, err
	}
	if subdomain != "" {
		certificateName = "c-" + rootDomain + "-" + strings.ReplaceAll(subdomain, ".", "-")
	} else {
		certificateName = "c-" + rootDomain
	}
	fmt.Printf("certificateName: %v\n", certificateName)

	body := ExternalCertificate{
		CertificateName:      certificateName,
		PrivateKey:           cleanPEM(privateKey),
		PublicKeyCertificate: cleanPEM(certificateBody),
		CertificateChain:     cleanPEM(certificateChain) + "\n" + cleanPEM(rootCA),
	}
	resp, err := services.Request(access, "POST", endpoint, url, body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to HTTP CreateExternalCertificate: %v", resp)
	}

	certificateReturn := CertificateReturn{}
	if err := json.NewDecoder(resp.Body).Decode(&certificateReturn); err != nil {
		return nil, fmt.Errorf("failed to CreateExternalCertificate JSON: %v", err)
	}
	defer resp.Body.Close()

	if certificateReturn.ReturnCode != "0" {
		return nil, fmt.Errorf("failed to CreateExternalCertificate: %v", certificateReturn)
	}

	for _, certificate := range certificateReturn.SslCertificateList {
		if certificate.CertificateName == certificateName {
			return &certificate, nil
		}
	}

	return nil, fmt.Errorf("failed to CreateExternalCertificate?: %v", certificateReturn)
}

func DeleteCertificate() {

}
