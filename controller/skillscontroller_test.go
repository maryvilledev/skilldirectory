package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"skilldirectory/data"
	"skilldirectory/model"
	"testing"
)

func TestSkillsControllerBase(t *testing.T) {
	base := BaseController{}
	sc := SkillsController{BaseController: &base}

	if base != *sc.Base() {
		t.Error("Expecting Base() to return base pointer")
	}
}

func TestGetAllSkills(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/skills", nil)
	sc := getSkillsController(request, &data.MockDataAccessor{})

	err := sc.Get()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetAllSkills_Error(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/skills", nil)
	sc := getSkillsController(request, &data.MockErrorDataAccessor{})

	err := sc.Get()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestGetSkill(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/skills/1234", nil)
	sc := getSkillsController(request, &data.MockDataAccessor{})

	err := sc.Get()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetSkill_Error(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/skills/1234", nil)
	sc := getSkillsController(request, &data.MockErrorDataAccessor{})

	err := sc.Get()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

// func TestGetSkillsFiltered(t *testing.T) {
// 	base := BaseController{}
// 	skillsConnector := SkillsController{BaseController: &base}
//
// 	responseRecorder := httptest.NewRecorder()
// 	request := httptest.NewRequest(http.MethodGet, "/skills?skilltype=scripted", nil)
// 	accessor := Newdata.MockInMemoryDataAccessor()
// 	base.Init(responseRecorder, request, &accessor)
//
// 	newScriptedSkill := model.NewSkill("1234", "TestSkillName", model.ScriptedSkillType)
// 	skillsConnector.session.Save(newScriptedSkill.Id, newScriptedSkill)
// 	newCompiledSkill := model.NewSkill("2136", "TestSkillName", model.CompiledSkillType)
// 	skillsConnector.session.Save(newCompiledSkill.Id, newCompiledSkill)
//
// 	err := skillsConnector.Get()
// 	if err != nil {
// 		t.Errorf("Did not expect error when getting skills with filter")
// 	}
//
// 	correctResponseBody := "[{\"Blogs\":[],\"id\":\"1234\",\"name\":\"TestSkillName\",\"skilltype\":\"scripted\"}]"
// 	if responseRecorder.Body.String() != correctResponseBody {
// 		t.Errorf("Failed to properly filter based on skilltype. "+
// 			"Expected Response body to be \n\t %s\n But got\n\t %s\\n",
// 			correctResponseBody, responseRecorder.Body.String())
// 	}
// }
//
// func TestGetSkillsFilteredBadSkillType(t *testing.T) {
// 	base := BaseController{}
// 	skillsConnector := SkillsController{BaseController: &base}
//
// 	responseRecorder := httptest.NewRecorder()
// 	request := httptest.NewRequest(http.MethodGet, "/skills?skilltype=invalid-skill-type", nil)
// 	accessor := Newdata.MockInMemoryDataAccessor()
// 	base.Init(responseRecorder, request, &accessor)
//
// 	err := skillsConnector.Get()
// 	if err == nil {
// 		t.Errorf("Expected error due to invalid skill type.")
// 	}
// }
//
//func TestGetSkillsFilteredError(t *testing.T) {
//	base := BaseController{}
//	skillsConnector := SkillsController{BaseController: &base}
//
//	responseRecorder := httptest.NewRecorder()
//	request := httptest.NewRequest(http.MethodGet, "/skills?skilltype=scripted", nil)
//	accessor := data.MockErrorDataAccessor{}
//	base.Init(responseRecorder, request, &accessor)
//
//	err := skillsConnector.Get()
//	if err == nil {
//		t.Errorf("Expecting error for TestGetSkillsFilteredError")
//	}
//}

func TestDeleteSkill(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/skills/1234", nil)
	sc := getSkillsController(request, &data.MockDataAccessor{})

	err := sc.Delete()
	if err != nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestDeleteSkill_Error(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/skills/1234", nil)
	sc := getSkillsController(request, &data.MockErrorDataAccessor{})

	err := sc.Delete()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestDeleteSkill_NoKey(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/skills/", nil)
	sc := getSkillsController(request, &data.MockDataAccessor{})

	err := sc.Delete()
	if err == nil {
		t.Errorf("Expected error when no key: %s", err.Error())
	}
}

func TestPostSkill(t *testing.T) {
	body := getReaderForNewSkill("", "BestSkillNameEver", model.ScriptedSkillType)
	sc := getSkillsController(httptest.NewRequest(http.MethodPost, "/skills", body), &data.MockDataAccessor{})

	err := sc.Post()
	if err != nil {
		t.Errorf("Post failed: %s", err.Error())
	}
}

func TestPostSkill_NoName(t *testing.T) {
	body := getReaderForNewSkill("1234", "", model.ScriptedSkillType)
	sc := getSkillsController(httptest.NewRequest(http.MethodPost, "/skills", body), &data.MockDataAccessor{})

	err := sc.Post()
	if err == nil {
		t.Errorf("Expected error due to empty %q field in Skill POST request.", "name")
	}
}

func TestPostSkill_NoSkillType(t *testing.T) {
	body := getReaderForNewSkill("1234", "SomeName", "")
	sc := getSkillsController(httptest.NewRequest(http.MethodPost, "/skills", body), &data.MockDataAccessor{})

	err := sc.Post()
	if err == nil {
		t.Errorf("Expected error due to empty %q field in Skill POST request.", "skill_type")
	}
}

func TestPostSkill_InvalidType(t *testing.T) {
	body := getReaderForNewSkill("", "", "badtype")
	sc := getSkillsController(httptest.NewRequest(http.MethodPost, "/skills", body), &data.MockDataAccessor{})

	err := sc.Post()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestPostSkill_NoSkill(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/skills/", nil)
	sc := getSkillsController(request, &data.MockDataAccessor{})

	err := sc.Post()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

func TestPostSkill_Error(t *testing.T) {
	body := getReaderForNewSkill("", "", model.ScriptedSkillType)
	request := httptest.NewRequest(http.MethodPost, "/skills", body)
	sc := getSkillsController(request, &data.MockErrorDataAccessor{})

	err := sc.Post()
	if err == nil {
		t.Errorf("Expected error: %s", err.Error())
	}
}

/*
getSkillsController is a helper function for creating and initializing a new BaseController with
the given HTTP request and DataAccessor. Returns a new SkillsController created with that BaseController.
*/
func getSkillsController(request *http.Request, dataAccessor data.DataAccess) SkillsController {
	base := BaseController{}
	base.Init(httptest.NewRecorder(), request, dataAccessor)
	return SkillsController{BaseController: &base}
}

/*
getReaderForNewSkill is a helper function for a new Skill with the given id, name, and skillType.
This Skill is then marshaled into JSON. A new Reader is created and returned for the resulting []byte.
*/
func getReaderForNewSkill(id, name, skillType string) *bytes.Reader {
	newSkill := model.NewSkill(id, name, skillType)
	b, _ := json.Marshal(newSkill)
	return bytes.NewReader(b)
}