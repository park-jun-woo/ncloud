package CertificateManager

import (
	"encoding/json"
	"fmt"
	"github.com/park-jun-woo/ncloud-sdk-go/services"

)

type Certificates struct {
	ReturnCode string `json:"returnCode"`
	ReturnMessage string `json:"returnMessage"`
	TotalRows int `json:"totalRows"`
	SslCertificateList []Certificate `json:"sslCertificateList"`
}

type Certificate struct {
	CertificateNo int `json:"certificateNo"`
	CertificateType string `json:"certificateType"`
	CertificateName string `json:"certificateName"`
	MemberNo string `json:"memberNo"`
	DnInfo string `json:"dnInfo"`
	DomainAddress string `json:"domainAddress"`
	RegDate string `json:"regDate"`
	UpdateDate string `json:"updateDate"`
	IssueDate string `json:"issueDate"`
	ValidStartDate string `json:"validStartDate"`
	ValidEndDate string `json:"validEndDate"`
	StatusCode string `json:"statusCode"`
	StatusName string `json:"statusName"`
	ExternalYn string `json:"externalYn"`
	DomainCode string `json:"domainCode"`
	CaInfo string `json:"caInfo"`
	CertSerialNumber string `json:"certSerialNumber"`
	CertPublicKeyInfo string `json:"certPublicKeyInfo"`
	CertSignAlgorithmName string `json:"certSignAlgorithmName"`
}

type CertificateReturn struct {
	ReturnCode string `json:"returnCode"`
	ReturnMessage string `json:"returnMessage"`
	TotalRows int `json:"totalRows"`
	SslCertificateList []Certificate `json:"sslCertificateList"`
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
		return nil, fmt.Errorf("Failed to GetCertificates JSON: %v", err)
	}
	defer resp.Body.Close()

	if certificates.ReturnCode != "0" {
		return nil, fmt.Errorf("Failed to GetCertificates: %v", certificates)
	}

	return &certificates, nil
}

func CreateExternalCertificate(access *services.Access, certificateName string, privateKey string, publicKeyCertificate string, certificateBody string, rootCA string) (*Certificate, error) {
	endpoint := "https://certificatemanager.apigw.ntruss.com"
	url := "/api/v1/certificate/withExternal"
	body := map[string]interface{}{
		"certificateName": certificateName,
		"privateKey": privateKey,
		"publicKeyCertificate": publicKeyCertificate,
		"certificateChain": certificateBody+rootCA,
	}
	resp, err := services.Request(access, "POST", endpoint, url, body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Failed to HTTP CreateExternalCertificate: %v", resp)
	}
	
	certificateReturn := CertificateReturn{}
	if err := json.NewDecoder(resp.Body).Decode(&certificateReturn); err != nil {
		return nil, fmt.Errorf("Failed to CreateExternalCertificate JSON: %v", err)
	}
	defer resp.Body.Close()

	if certificateReturn.ReturnCode != "0" {
		return nil, fmt.Errorf("Failed to CreateExternalCertificate: %v", certificateReturn)
	}

	for _, certificate := range certificateReturn.SslCertificateList {
		if certificate.CertificateName == certificateName {
			return &certificate, nil
		}
	}

	return nil, fmt.Errorf("Failed to CreateExternalCertificate: %v", certificateReturn)
}