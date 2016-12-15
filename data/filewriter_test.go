package data

import (
	"skilldirectory/model"
	"testing"
)

var dataConnector = NewAccessor(NewFileWriter(""))
var testFileName = "testwrite"
var badFileName = "testbadread"

func TestWriteFile(t *testing.T) {
	t.Log("Testing Write File")
	newSkill := model.NewSkill("01234", "test", "language", model.Link{}, []model.Link{}, []model.Link{})
	dataConnector.Save(testFileName, newSkill)
}

func TestBadWriteFile(t *testing.T) {
	t.Log("Testing Bad Write File")
	badMarshal := func() {}
	err := dataConnector.Save(badFileName, badMarshal)
	if err == nil {
		t.Errorf("Expecting Error")
	}
}

func TestReadFile(t *testing.T) {
	t.Log("Testing Read File")
	var skill = model.Skill{}
	err := dataConnector.Read(testFileName, &skill)
	if err != nil {
		t.Errorf("Read File Failed: %v", err)
	}
}

func TestBadReadFile(t *testing.T) {
	t.Log("Testing Bad Read File")
	var skill = model.Skill{}
	err := dataConnector.Read(badFileName, &skill)
	if err == nil {
		t.Errorf("Bad Read File Failed: %v", err)
	}

}

func TestDelete(t *testing.T) {
	t.Log("Testing Delete File")
	err := dataConnector.Delete(testFileName)
	if err != nil {
		t.Errorf("Delete File Failed: %v", err)
	}
}
