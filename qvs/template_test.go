package qvs

import (
	"testing"
)

func TestTemplateCRUD(t *testing.T) {
	if skipTest() {
		t.SkipNow()
	}
	c := getTestManager()

	tmpl := &Template{
		Name:               "test001",
		Bucket:             "xxxx",
		TemplateType:       1,
		JpgOverwriteStatus: true,
		RecordType:         2,
	}
	tmpl1, err := c.AddTemplate(tmpl)
	noError(t, err)

	tmpl2, err := c.QueryTemplate(tmpl1.ID)
	noError(t, err)
	shouldBeEqual(t, tmpl.Name, tmpl2.Name)
	shouldBeEqual(t, tmpl.Bucket, tmpl2.Bucket)
	shouldBeEqual(t, tmpl.TemplateType, tmpl2.TemplateType)
	shouldBeEqual(t, tmpl.JpgOverwriteStatus, tmpl2.JpgOverwriteStatus)
	shouldBeEqual(t, tmpl.RecordType, tmpl2.RecordType)

	ops := []PatchOperation{
		{
			Op:    "replace",
			Key:   "name",
			Value: "test002",
		},
	}
	tmpl3, err := c.UpdateTemplate(tmpl2.ID, ops)
	noError(t, err)
	shouldBeEqual(t, "test002", tmpl3.Name)
	shouldBeEqual(t, tmpl.Bucket, tmpl3.Bucket)
	shouldBeEqual(t, tmpl.TemplateType, tmpl3.TemplateType)
	shouldBeEqual(t, tmpl.JpgOverwriteStatus, tmpl3.JpgOverwriteStatus)
	shouldBeEqual(t, tmpl.RecordType, tmpl3.RecordType)

	tmpl4 := &Template{
		Name:               "test003",
		Bucket:             "xxxx",
		TemplateType:       1,
		JpgOverwriteStatus: true,
		RecordType:         2,
	}
	tmpl5, err := c.AddTemplate(tmpl4)
	noError(t, err)

	tmpl6 := &Template{
		Name:               "test004",
		Bucket:             "xxxx",
		TemplateType:       1,
		JpgOverwriteStatus: true,
		RecordType:         2,
	}
	tmpl7, err := c.AddTemplate(tmpl6)
	noError(t, err)

	nss, total, err := c.ListTemplate(0, 2, 0, "", 1)
	noError(t, err)
	shouldBeEqual(t, int64(3), total)
	shouldBeEqual(t, 2, len(nss))

	nss, total, err = c.ListTemplate(2, 2, 0, "", 1)
	noError(t, err)
	shouldBeEqual(t, int64(3), total)
	shouldBeEqual(t, 1, len(nss))

	c.DeleteTemplate(tmpl1.ID)
	c.DeleteTemplate(tmpl5.ID)
	c.DeleteTemplate(tmpl7.ID)
}
