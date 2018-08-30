package main

import (
	"testing"
	"time"
)

func testVersionOfCertIssue(domain string) (string, time.Time, error) {
	return "testVersion", time.Now().AddDate(0, 0, 1), nil
}
func TestGetCertShouldSucceed(t *testing.T) {

	var certs = NewCertContainer(testVersionOfCertIssue)

	_, err := certs.GetCert("test")

	if err != nil {
		t.Error("Getting cert should have succeed")
	}
}

func TestGetCertShouldAddToContainer(t *testing.T) {

	var certs = NewCertContainer(testVersionOfCertIssue)

	certs.GetCert("test")

	if len(certs.certs) != 1 {
		t.Errorf("Cert was not added for tracking after get, len is %d", len(certs.certs))
	}
}

func TestGetExpiredCertsShouldReturnAllExpired(t *testing.T) {

	var certs = NewCertContainer(testVersionOfCertIssue)

	var threshold = time.Now()

	var expired = len(certs.GetExpired(threshold))
	if expired != 0 {
		t.Errorf("There should be zero expired certs, was %d instead", expired)
	}

	certs.GetCert("domain")

	threshold = threshold.AddDate(0, 1, 0)

	expired = len(certs.GetExpired(threshold))
	if expired != 1 {
		t.Errorf("There should be one expired cert, was % d instead", expired)
	}
}

func TestRenewNonExistingShouldFail(t *testing.T) {
	var certs = NewCertContainer(testVersionOfCertIssue)

	_, err := certs.RenewCert("nonexistent")

	if err == nil {
		t.Errorf("Error should have been returned trying to renew non existing domain")
	}
}
