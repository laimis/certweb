package main

import (
	"sync"
	"time"
)

type refreshCertService func(domain string) (string, time.Time, error)

// CertContainer holds all the certs
type CertContainer struct {
	sync.Mutex
	certs       map[string]Cert
	refreshFunc refreshCertService
}

// Cert describes the cert
type Cert struct {
	Domain  string
	CertVal string
	Expires time.Time
}

// NewCertContainer inits certs
func NewCertContainer(refreshFunc refreshCertService) CertContainer {
	return CertContainer{
		certs:       make(map[string]Cert),
		refreshFunc: refreshFunc,
	}
}

func getCert(domain string, refreshFunc refreshCertService) (Cert, error) {
	var certValue, expiration, err = refreshFunc(domain)
	if err != nil {
		return Cert{}, err
	}

	var cert = Cert{
		CertVal: certValue,
		Domain:  domain,
		Expires: expiration,
	}

	return cert, nil
}

// GetCert returns cert for a given id
func (container *CertContainer) GetCert(domain string) (Cert, error) {
	container.Lock()
	defer container.Unlock()

	cert, ok := container.certs[domain]
	if !ok {

		issuedCert, err := getCert(domain, container.refreshFunc)
		if err != nil {
			return Cert{}, err
		}

		container.certs[domain] = issuedCert
		cert = issuedCert
	}

	return cert, nil
}

// RenewCert returns cert for a given id
func (container *CertContainer) RenewCert(domain string) (Cert, error) {
	container.Lock()
	defer container.Unlock()

	cert, err := getCert(domain, container.refreshFunc)
	if err != nil {
		return Cert{}, err
	}

	container.certs[domain] = cert

	return cert, nil
}

// GetExpired returns expired certs
func (container *CertContainer) GetExpired(expiration time.Time) []Cert {
	container.Lock()
	defer container.Unlock()

	expiredCerts := []Cert{}

	for _, v := range container.certs {
		if v.Expires.Before(expiration) {
			expiredCerts = append(expiredCerts, v)
		}
	}

	return expiredCerts
}
