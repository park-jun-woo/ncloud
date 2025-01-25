package GlobalDNS

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"golang.org/x/net/publicsuffix"
	"parkjunwoo.com/ncloud-sdk-go/services"
)

// 루트 도메인과 서브도메인을 분리하는 함수
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

// 루트 도메인 이름과 서브도메인 이름을 분리하는 함수, .com, .co.kr 등을 제외한 도메인 이름만 추출
func GetDomainNames(domainName string) (string, string, error) {
	// Get the effective top-level domain plus one (eTLD+1), which is the root domain
	eTLDPlusOne, err := publicsuffix.EffectiveTLDPlusOne(domainName)
	if err != nil {
		return domainName, "", nil
	}

	// Extract the root domain name without the TLD
	splitDomain := strings.SplitN(eTLDPlusOne, ".", 2)
	if len(splitDomain) < 2 {
		return "", "", fmt.Errorf("unable to extract root domain from: %v", eTLDPlusOne)
	}
	rootDomain := splitDomain[0]

	// Extract subdomain if it exists
	subdomain := ""
	if domainName != eTLDPlusOne {
		subdomain = strings.TrimSuffix(domainName, "."+eTLDPlusOne)
	}

	return rootDomain, subdomain, nil
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
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to HTTP GetDomain: %v", resp)
	}

	domains := Domains{}
	if err := json.NewDecoder(resp.Body).Decode(&domains); err != nil {
		return nil, fmt.Errorf("failed to GetDomain JSON: %v", err)
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
		return nil, fmt.Errorf("failed to HTTP PostDomain: %v", resp)
	}

	return GetDomain(access, domainName, false)
}

func ApplyDomain(access *services.Access, domainName string) (*Domain, error) {
	domain, err := GetDomain(access, domainName, true)
	if err != nil {
		return nil, err
	}

	endpoint := "https://globaldns.apigw.ntruss.com"
	url := fmt.Sprintf("/dns/v1/ncpdns/record/apply/%d", domain.Id)
	resp, err := services.Request(access, "PUT", endpoint, url, nil)
	if err != nil {
		return domain, err
	}
	if resp.StatusCode != 200 {
		return domain, fmt.Errorf("failed to HTTP PostDomain: %v", resp)
	}

	return domain, nil
}
