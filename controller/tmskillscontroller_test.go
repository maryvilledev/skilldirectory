package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"skilldirectory/data"
	"skilldirectory/model"
	"testing"

	"github.com/Sirupsen/logrus"
)

func TestTMSkillsController_Base(t *testing.T) {
	base := BaseController{}
	tc := TMSkillsController{BaseController: &base}

	if base != *tc.Base() {
		t.Error("Expected Base() to return base pointer")
	}
}

func TestGetAllTMSkills(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/tmskills", nil)
	tc := getTMSkillsController(request, &data.MockDataAccessor{})

	err := tc.Get()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetAllTMSkills_Error(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/tmskills", nil)
	tc := getTMSkillsController(request, &data.MockErrorDataAccessor{})

	err := tc.Get()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestGetTMSkill(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/tmskills/1234", nil)
	tc := getTMSkillsController(request, &data.MockDataAccessor{})

	err := tc.Get()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetTMSkill_Error(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/tmskills/1234", nil)
	tc := getTMSkillsController(request, &data.MockErrorDataAccessor{})

	err := tc.Get()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestDeleteTMSkill(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/tmskills/1234", nil)
	tc := getTMSkillsController(request, &data.MockDataAccessor{})

	err := tc.Delete()
	if err != nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestDeleteTMSkill_Error(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/tmskills/1234", nil)
	tc := getTMSkillsController(request, &data.MockErrorDataAccessor{})

	err := tc.Delete()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestDeleteTMSkill_NoKey(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/tmskills/", nil)
	tc := getTMSkillsController(request, &data.MockDataAccessor{})

	err := tc.Delete()
	if err == nil {
		t.Errorf("Expected error when no key: %s", err.Error())
	}
}

func TestPostTMSkill(t *testing.T) {
	body := getReaderForNewTMSkill("1234", "2345", "3456")
	request := httptest.NewRequest(http.MethodPost, "/tmskills", body)
	tc := getTMSkillsController(request, &data.MockDataAccessor{})

	err := tc.Post()
	if err != nil {
		t.Errorf("Post failed: %s", err.Error())
	}
}

func TestPostTMSkill_NoSkillID(t *testing.T) {
	body := getReaderForNewTMSkill("1234", "", "3456")
	request := httptest.NewRequest(http.MethodPost, "/tmskills", body)
	tc := getTMSkillsController(request, &data.MockDataAccessor{})

	err := tc.Post()
	if err == nil {
		t.Errorf("Expected error due to empty %q field in TMSkill POST"+
			" request.", "skill_id")
	}
}

func TestPostTMSkill_NoTeamMemberID(t *testing.T) {
	body := getReaderForNewTMSkill("1234", "2345", "")
	request := httptest.NewRequest(http.MethodPost, "/tmskills", body)
	tc := getTMSkillsController(request, &data.MockDataAccessor{})

	err := tc.Post()
	if err == nil {
		t.Errorf("Expected error due to empty %q field in TMSkill POST"+
			" request.", "team_member_id")
	}
}

func TestPostTMSkill_NoTMSkill(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/tmskills", nil)
	tc := getTMSkillsController(request, &data.MockDataAccessor{})

	err := tc.Post()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestPostTMSkill_Error(t *testing.T) {
	body := getReaderForNewTMSkill("1234", "2345", "3456")
	request := httptest.NewRequest(http.MethodPost, "/tmskills", body)
	tc := getTMSkillsController(request, &data.MockErrorDataAccessor{})

	err := tc.Post()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func Test_validateTMSkillFields(t *testing.T) {
	tc := getTMSkillsController(nil, &data.MockErrorDataAccessor{})
	tmSkill := model.TMSkill{
		SkillID: "1234",
	}
	err := tc.validateTMSkillFields(&tmSkill)
	if err == nil {
		t.Errorf("validateTMSkillFields() failed to detect empty " +
			"TMSkill.TeamMemberID field.")
	}

	tmSkill = model.TMSkill{
		TeamMemberID: "1234",
	}
	err = tc.validateTMSkillFields(&tmSkill)
	if err == nil {
		t.Errorf("validateTMSkillFields() failed to detect empty " +
			"TMSkill.SkillID field.")
	}

	tmSkill = model.TMSkill{
		SkillID:      "1234",
		TeamMemberID: "1234",
	}
	err = tc.validateTMSkillFields(&tmSkill)
	if err == nil {
		t.Errorf("validateTMSkillFields() failed to detect invalid " +
			"ID field.")
	}

	tmSkill = model.TMSkill{
		SkillID:      "1234",
		TeamMemberID: "1234",
		Proficiency:  9000,
	}
	err = tc.validateTMSkillFields(&tmSkill)
	if err == nil {
		t.Errorf("validateTMSkillFields() failed to detect invalid " +
			"TMSkill.Proficiency field.")
	}
}

func Test_validateTMSkillID(t *testing.T) {
	tc := getTMSkillsController(nil, &data.MockErrorDataAccessor{})
	tmSkill := model.TMSkill{
		ID: "1234",
	}
	err := tc.validateTMSkillID(&tmSkill)
	if err == nil {
		t.Errorf("validateTMSkillID() failed to detect invalid TMSkill.ID field.")
	}
}

func TestConvertToStruct(t *testing.T) {
	preTMSkills := []interface{}{model.TMSkill{SkillID: "1234"}, model.TMSkill{SkillID: "5678"}}
	pretTMSkillsStruct := []model.TMSkill{model.TMSkill{SkillID: "1234"}, model.TMSkill{SkillID: "5678"}}
	tmskill, err := convertToStruct(preTMSkills)
	if err != nil {
		t.Errorf("Convert to struct failed: %v", err)
	}

	if !reflect.DeepEqual(pretTMSkillsStruct, tmskill) {
		t.Error("Deep equal failed fro Convert to struct")
	}
}

func TestConvertSkillsToDTOs(t *testing.T) {
	preTMSkills := []model.TMSkill{model.TMSkill{SkillID: "1234"}, model.TMSkill{SkillID: "5678"}}
	preTMSkillsDTO := []model.TMSkillDTO{preTMSkills[0].NewTMSkillDTO("", ""), preTMSkills[1].NewTMSkillDTO("", "")}
	tmController := getTMSkillsController(nil, data.MockDataAccessor{})
	tmSkillsDTO := tmController.converSkillsToDTOs(preTMSkills)
	if !reflect.DeepEqual(preTMSkillsDTO, tmSkillsDTO) {
		t.Error("Expecting a match of tmskills -> tmskillsDTO")
	}
}

func TestConvertSkillsToDTOsError(t *testing.T) {
	preTMSkills := []model.TMSkill{model.TMSkill{SkillID: "1234"}, model.TMSkill{SkillID: "5678"}}
	preTMSkillsDTO := []model.TMSkillDTO{}
	tmController := getTMSkillsController(nil, data.MockErrorDataAccessor{})
	tmSkillsDTO := tmController.converSkillsToDTOs(preTMSkills)
	if !reflect.DeepEqual(preTMSkillsDTO, tmSkillsDTO) {
		t.Error("Expecting a match of tmskills -> tmskillsDTO")
	}

}

/*
getTMSkillsController is a helper function for creating and initializing a new
BaseController with the given HTTP request and DataAccessor. Returns a new
TMSkillsController created with that BaseController.
*/
func getTMSkillsController(request *http.Request, dataAccessor data.DataAccess) TMSkillsController {
	base := BaseController{}
	base.Init(httptest.NewRecorder(), request, dataAccessor, nil, logrus.New())
	return TMSkillsController{BaseController: &base}
}

/*
getReaderForNewTMSkill is a helper function for a new TMSkill with the given id,
skillID, and teamMemberID. This TMSkill is then marshaled into JSON. A new Reader
is created and returned for the resulting []byte.
*/
func getReaderForNewTMSkill(id, skillID, teamMemberID string) *bytes.Reader {
	newTMSkill := model.NewTMSkillDefaults(id, skillID, teamMemberID)
	b, _ := json.Marshal(newTMSkill)
	return bytes.NewReader(b)
}

type WrapInterface struct {
	array []interface{}
}
