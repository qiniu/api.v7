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
	domain := []string{"qtest.com"}
	ns := &NameSpace{
		Name:        nsName,
		AccessType:  nsAccessType,
		RTMPURLType: nsRTMPURLType,
		Domains:     domain,
	}
	ns1, err := c.AddNamespace(ns)
	noError(t, err)

	ns2, err := c.QueryNamespace(ns1.ID)
	noError(t, err)
	shouldBeEqual(t, nsName, ns2.Name)
	shouldBeEqual(t, nsAccessType, ns2.AccessType)
	shouldBeEqual(t, nsRTMPURLType, ns2.RTMPURLType)
	shouldBeEqual(t, domain[0], ns2.Domains[0])

	ops := []PatchOperation{
		{
			Op:    "replace",
			Key:   "name",
			Value: "testNamespace2",
		},
	}
	ns3, err := c.UpdateNamespace(ns1.ID, ops)
	noError(t, err)
	shouldBeEqual(t, "testNamespace2", ns3.Name)
	shouldBeEqual(t, nsAccessType, ns3.AccessType)
	shouldBeEqual(t, nsRTMPURLType, ns3.RTMPURLType)
	shouldBeEqual(t, domain[0], ns3.Domains[0])

	ns4 := &NameSpace{
		Name:        "testNamespace3",
		AccessType:  "rtmp",
		RTMPURLType: 1,
		Domains:     []string{"atest.com"},
	}
	ns5, err := c.AddNamespace(ns4)
	noError(t, err)

	ns6 := &NameSpace{
		Name:       "testNamespace4",
		AccessType: "gb28181",
	}
	ns7, err := c.AddNamespace(ns6)
	noError(t, err)

	nss, total, err := c.ListNamespace(0, 2, "")
	noError(t, err)
	shouldBeEqual(t, int64(3), total)
	shouldBeEqual(t, 2, len(nss))

	nss, total, err = c.ListNamespace(2, 2, "")
	noError(t, err)
	shouldBeEqual(t, int64(3), total)
	shouldBeEqual(t, 1, len(nss))

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
	domain := []string{"qtest.com"}
	ns := &NameSpace{
		Name:        nsName,
		AccessType:  nsAccessType,
		RTMPURLType: nsRTMPURLType,
		Domains:     domain,
	}
	ns1, err := c.AddNamespace(ns)
	noError(t, err)

	err = c.AddDomain(ns1.ID, &DomainInfo{Domain: "test001.com", Type: "publishRtmp"})
	noError(t, err)

	err = c.DeleteDomain(ns1.ID, "test001.com")
	noError(t, err)

	err = c.AddDomain(ns1.ID, &DomainInfo{Domain: "test001.com", Type: "publishRtmp"})
	noError(t, err)

	err = c.AddDomain(ns1.ID, &DomainInfo{Domain: "test002.com", Type: "publishRtmp"})
	noError(t, err)

	err = c.AddDomain(ns1.ID, &DomainInfo{Domain: "test003.com", Type: "publishRtmp"})
	noError(t, err)

	domains, err := c.ListDomain(ns1.ID)
	noError(t, err)
	shouldBeEqual(t, 7, len(domains))
	c.DeleteNamespace(ns1.ID)
}
