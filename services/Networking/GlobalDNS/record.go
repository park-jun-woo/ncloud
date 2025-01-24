package GlobalDNS

import (
	"encoding/json"
	"fmt"

	"parkjunwoo.com/ncloud-sdk-go/services"
)

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
	if resp.StatusCode != 200 {
		return domain, nil, fmt.Errorf("Failed to HTTP GetRecord: %v", resp)
	}

	records := Records{}
	if err := json.NewDecoder(resp.Body).Decode(&records); err != nil {
		return domain, nil, fmt.Errorf("Failed to GetRecord JSON: %v", err)
	}
	defer resp.Body.Close()

	for _, record := range records.Content {
		if (host == "" || record.Host == host) &&
			(recordType == "" || record.Type == recordType) &&
			(recordContent == "" || record.Content == recordContent) &&
			record.DelYn == false {
			return domain, &record, nil
		}
	}

	return domain, nil, nil
}

// 레코드 설정, 있으면 수정하고 없으면 생성한다.
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

	_, err = ApplyDomain(access, domainName)
	if err != nil {
		return domain, nil, err
	}

	return domain, nil, nil
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
		return domain, nil, fmt.Errorf("Failed to HTTP PostDomain: %v", resp)
	}

	return GetRecord(access, domainName, recordType, recordContent, false)
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
		return domain, nil, fmt.Errorf("Failed to HTTP PostDomain: %v", resp)
	}

	return GetRecord(access, domainName, recordType, recordContent, false)
}

// 레코드 삭제
func DeleteRecord(access *services.Access, domainName string, recordType string, recordContent string) error {
	domain, record, err := GetRecord(access, domainName, recordType, recordContent, false)
	if err != nil {
		return err
	}
	if record == nil {
		return fmt.Errorf("record(%s %s) is not exists.", domainName, recordType)
	}

	endpoint := "https://globaldns.apigw.ntruss.com"
	url := fmt.Sprintf("/dns/v1/ncpdns/record/%d", domain.Id)
	body := []int{record.Id}
	resp, err := services.Request(access, "DELETE", endpoint, url, body)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Failed to HTTP DeleteRecord: %v", resp)
	}

	_, err = ApplyDomain(access, domainName)
	if err != nil {
		return err
	}

	return nil
}
