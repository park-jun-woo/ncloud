package GlobalDNS

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/park-jun-woo/ncloud/services"
	"golang.org/x/net/publicsuffix"
)

// 루트 도메인과 서브도메인을 분리하는 함수수
func GetDomainParts(domainName string) (string, string, error) {
	eTLDPlusOne, err := publicsuffix.EffectiveTLDPlusOne(domainName)
	if err != nil {
		return "", "", fmt.Errorf("Invalid domain: %v", err)
	}

	// 서브도메인 구분
	subdomain := ""
	if domainName != eTLDPlusOne {
		subdomain = domainName[:len(domainName)-len(eTLDPlusOne)-1]
	}

	return eTLDPlusOne, subdomain, nil
}

// 도메인 조회
func GetDomain(access *services.Access, domainName string, postDomain bool) (*Domain, error) {
	rootDomain, _, err := GetDomainParts(domainName)
	if err != nil {
		return nil, fmt.Errorf("Invalid Domain Name: %v", err)
	}
	endpoint := "https://globaldns.apigw.ntruss.com"
	url := fmt.Sprintf("/dns/v1/ncpdns/domain?page=0&size=10&domainName=%s", rootDomain)
	resp, err := services.Request(access, "GET", endpoint, url, nil)
	if err != nil {
		return nil, err
	}

	domains := Domains{}
	if err := json.NewDecoder(resp.Body).Decode(&domains); err != nil {
		return nil, fmt.Errorf("Failed to GetDomain JSON: %v", err)
	}
	defer resp.Body.Close()

	for _, domain := range domains.Content {
		if domain.Name == rootDomain {
			return &domain, nil
		}
	}

	if len(domains.Content) == 0 && postDomain == true {
		log.Printf("도메인이 존재하지 않아 신규등록. %s", domainName)
		domain, err := PostDomain(access, rootDomain, "")
		return domain, err
	}

	return nil, nil
}

// 도메인 등록
func PostDomain(access *services.Access, domainName string, comments string) (*Domain, error) {
	rootDomain, _, err := GetDomainParts(domainName)
	if err != nil {
		return nil, fmt.Errorf("Invalid Domain Name: %v", err)
	}

	endpoint := "https://globaldns.apigw.ntruss.com"
	url := "/dns/v1/ncpdns/domain"
	body := DomainCreateRequest{Name: rootDomain, Comments: comments}
	resp, err := services.Request(access, "POST", endpoint, url, body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Failed to PostDomain: %v", resp)
	} else {
		return GetDomain(access, domainName, false)
	}
}

// 레코드 조회
func GetRecord(access *services.Access, domainName string, recordType string, recordContent string, postDomain bool) (*Domain, *Record, error) {
	_, host, err := GetDomainParts(domainName)
	if err != nil {
		return nil, nil, fmt.Errorf("Invalid Domain Name: %v", err)
	}
	if host == "" {
		host = "@"
	}

	domain, err := GetDomain(access, domainName, postDomain)
	if err != nil {
		return nil, nil, err
	}

	endpoint := "https://globaldns.apigw.ntruss.com"
	url := fmt.Sprintf("/dns/v1/ncpdns/record/%d?page=0&size=1000&recordType=%s", domain.Id, recordType)
	resp, err := services.Request(access, "GET", endpoint, url, nil)
	if err != nil {
		return domain, nil, err
	}

	records := Records{}
	if err := json.NewDecoder(resp.Body).Decode(&records); err != nil {
		return domain, nil, fmt.Errorf("Failed to GetRecord JSON: %v", err)
	}
	defer resp.Body.Close()

	for _, record := range records.Content {
		if host != "" && record.Host == host &&
			recordType != "" && record.Type == recordType &&
			recordContent != "" && record.Content == recordContent &&
			record.DelYn == false {
			return domain, &record, nil
		}
	}

	return domain, nil, nil
}

func SetRecord(access *services.Access, domainName string, recordType string, recordContent string, recordTtl int, postDomain bool) (*Domain, *Record, error) {
	domain, record, err := GetRecord(access, domainName, recordType, recordContent, postDomain)
	if err != nil {
		return domain, nil, err
	}

	if record == nil {
		return PostRecord(access, domain, domainName, recordType, recordContent, recordTtl)
	} else {
		return putRecord(access, domain, domainName, record.Id, recordType, recordContent, recordTtl)
	}
}

// 레코드 등록
func PostRecord(access *services.Access, domain *Domain, domainName string, recordType string, recordContent string, recordTtl int) (*Domain, *Record, error) {
	if domain == nil {
		return nil, nil, fmt.Errorf("Domain is nil")
	}

	_, host, err := GetDomainParts(domainName)
	if err != nil {
		return nil, nil, fmt.Errorf("Invalid Domain Name: %v", err)
	}

	endpoint := "https://globaldns.apigw.ntruss.com"
	url := fmt.Sprintf("/dns/v1/ncpdns/record/%d", domain.Id)
	body := []RecordCreateRequest{
		{
			Host:    host,
			Type:    recordType,
			Content: recordContent,
			Ttl:     recordTtl,
		},
	}
	resp, err := services.Request(access, "POST", endpoint, url, body)
	if err != nil {
		return domain, nil, err
	}

	if resp.StatusCode != 200 {
		return domain, nil, fmt.Errorf("Failed to PostDomain: %v", resp)
	} else {
		return GetRecord(access, domainName, recordType, recordContent, false)
	}
}

// 레코드 수정
func putRecord(access *services.Access, domain *Domain, domainName string, recordId int, recordType string, recordContent string, recordTtl int) (*Domain, *Record, error) {
	if domain == nil {
		return nil, nil, fmt.Errorf("Domain is nil")
	}

	_, host, err := GetDomainParts(domainName)
	if err != nil {
		return nil, nil, fmt.Errorf("Invalid Domain Name: %v", err)
	}

	endpoint := "https://globaldns.apigw.ntruss.com"
	url := fmt.Sprintf("/dns/v1/ncpdns/record/%d", domain.Id)
	body := []RecordUpdateRequest{
		{
			Id:      recordId,
			Host:    host,
			Type:    recordType,
			Content: recordContent,
			Ttl:     recordTtl,
		},
	}
	resp, err := services.Request(access, "PUT", endpoint, url, body)
	if err != nil {
		return domain, nil, err
	}

	if resp.StatusCode != 200 {
		return domain, nil, fmt.Errorf("Failed to PostDomain: %v", resp)
	} else {
		return GetRecord(access, domainName, recordType, recordContent, false)
	}
}
