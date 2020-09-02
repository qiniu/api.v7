package qvs

import (
	"testing"
)

func TestNamespaceCRUD(t *testing.T) {
	if skipTest() {
		t.SkipNow()
	}
	c := getTestManager()

	nsName := "testNamespace"
	nsAccessType := "rtmp"
	nsRTMPURLType := 1
	domain := []string{"qiniu1.com"}
	zone := "z0"
	hlslowlatency := false
	ns := &NameSpace{
		Name:          nsName,
		AccessType:    nsAccessType,
		RTMPURLType:   nsRTMPURLType,
		Domains:       domain,
		Zone:          zone,
		HLSLowLatency: hlslowlatency,
	}
	ns1, err := c.AddNamespace(ns)
	noError(t, err)

	ns2, err := c.QueryNamespace(ns1.ID)
	noError(t, err)
	shouldBeEqual(t, nsName, ns2.Name)
	shouldBeEqual(t, nsAccessType, ns2.AccessType)
	shouldBeEqual(t, nsRTMPURLType, ns2.RTMPURLType)
	shouldBeEqual(t, domain[0], ns2.Domains[0])
	shouldBeEqual(t, hlslowlatency, ns2.HLSLowLatency)

	nsName2 := "testNamespace2"
	nsAccessType2 := "gb28181"
	urlMode := 1
	domain2 := []string{"qiniuUM.com"}
	zone2 := "z0"
	hlslowlatency2 := false
	nsUM := &NameSpace{
		Name:          nsName2,
		AccessType:    nsAccessType2,
		UrlMode:       urlMode,
		Domains:       domain2,
		Zone:          zone2,
		HLSLowLatency: hlslowlatency2,
	}
	nsUMPost, err := c.AddNamespace(nsUM)
	noError(t, err)

	nsUMQuery, err := c.QueryNamespace(nsUMPost.ID)
	noError(t, err)
	shouldBeEqual(t, nsName2, nsUMQuery.Name)
	shouldBeEqual(t, nsAccessType2, nsUMQuery.AccessType)
	shouldBeEqual(t, urlMode, nsUMQuery.UrlMode)
	shouldBeEqual(t, domain2[0], nsUMQuery.Domains[0])
	shouldBeEqual(t, hlslowlatency2, nsUMQuery.HLSLowLatency)

	c.DeleteNamespace(nsUMPost.ID)

	ops := []PatchOperation{
		{
			Op:    "replace",
			Key:   "name",
			Value: "testNamespace2",
		},
		{
			Op:    "replace",
			Key:   "hlsLowLatency",
			Value: true,
		},
	}
	ns3, err := c.UpdateNamespace(ns1.ID, ops)
	noError(t, err)
	shouldBeEqual(t, "testNamespace2", ns3.Name)
	shouldBeEqual(t, nsAccessType, ns3.AccessType)
	shouldBeEqual(t, nsRTMPURLType, ns3.RTMPURLType)
	shouldBeEqual(t, domain[0], ns3.Domains[0])
	shouldBeEqual(t, true, ns3.HLSLowLatency)

	ns4 := &NameSpace{
		Name:        "testNamespace3",
		AccessType:  "rtmp",
		RTMPURLType: 1,
		Domains:     []string{"qiniu2.com"},
	}
	ns5, err := c.AddNamespace(ns4)
	noError(t, err)

	ns6 := &NameSpace{
		Name:        "testNamespace4",
		AccessType:  "gb28181",
		RTMPURLType: 2,
	}
	ns7, err := c.AddNamespace(ns6)
	noError(t, err)

	nss, _, err := c.ListNamespace(0, 2, "")
	noError(t, err)
	shouldBeEqual(t, 2, len(nss))

	nss, _, err = c.ListNamespace(1, 2, "")
	noError(t, err)
	shouldBeEqual(t, 2, len(nss))

	err = c.DisableNamespace(ns1.ID)
	noError(t, err)
	ret, err := c.QueryNamespace(ns1.ID)
	noError(t, err)
	shouldBeEqual(t, true, ret.Disabled)

	err = c.EnableNamespace(ns1.ID)
	noError(t, err)
	ret, err = c.QueryNamespace(ns1.ID)
	noError(t, err)
	shouldBeEqual(t, false, ret.Disabled)

	c.DeleteNamespace(ns1.ID)
	c.DeleteNamespace(ns5.ID)
	c.DeleteNamespace(ns7.ID)
}

func TestDomainCRUD(t *testing.T) {
	if skipTest() {
		t.SkipNow()
	}
	c := getTestManager()

	nsName := "testNamespace"
	nsAccessType := "rtmp"
	nsRTMPURLType := 1
	domain := []string{"qiniu1.com"}
	ns := &NameSpace{
		Name:        nsName,
		AccessType:  nsAccessType,
		RTMPURLType: nsRTMPURLType,
		Domains:     domain,
	}
	ns1, err := c.AddNamespace(ns)
	noError(t, err)

	err = c.AddDomain(ns1.ID, &DomainInfo{Domain: "qiniu2.com", Type: "publishRtmp"})
	noError(t, err)

	err = c.DeleteDomain(ns1.ID, "qiniu2.com")
	noError(t, err)

	err = c.AddDomain(ns1.ID, &DomainInfo{Domain: "qiniu2.com", Type: "publishRtmp"})
	noError(t, err)

	err = c.AddDomain(ns1.ID, &DomainInfo{Domain: "qiniu3.com", Type: "publishRtmp"})
	noError(t, err)

	err = c.AddDomain(ns1.ID, &DomainInfo{Domain: "qiniu4.com", Type: "publishRtmp"})
	noError(t, err)

	domains, err := c.ListDomain(ns1.ID)
	noError(t, err)
	shouldBeEqual(t, 7, len(domains))
	c.DeleteNamespace(ns1.ID)
}
