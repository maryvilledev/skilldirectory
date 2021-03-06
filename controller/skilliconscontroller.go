package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"skilldirectory/errors"
	"skilldirectory/model"
	"skilldirectory/util"
)

type SkillIconsController struct {
	*BaseController
}

func (c SkillIconsController) Base() *BaseController {
	return c.BaseController
}

func (c SkillIconsController) Get() error {
	return fmt.Errorf("Get Skills Icons not supported")
}

func (c SkillIconsController) Post() error {
	return c.addSkillIcon()
}

func (c SkillIconsController) Delete() error {
	return c.removeSkillIcon()
}

func (c SkillIconsController) Put() error {
	return c.addSkillIcon()
}

func (c SkillIconsController) Options() error {
	c.w.Header().Set("Access-Control-Allow-Headers", GetDefaultHeaders())
	c.w.Header().Set("Access-Control-Allow-Methods", "PUT, "+GetDefaultMethods())
	return nil
}

func (c *SkillIconsController) removeSkillIcon() error {
	// Get ID at end of request; return error if request contains no ID
	skillID := util.CheckForID(c.r.URL)
	if skillID == "" {
		return errors.MissingIDError(fmt.Errorf("no skill ID specified in request URL"))
	}

	// Attempt to delete image resource from S3
	err := c.fileSystem.Delete("dev/" + skillID)
	if err != nil {
		c.Warn(err)
		return err
	}
	skillIDInt, err := util.StringToID(skillID)
	if err != nil {
		return err
	}
	skill := model.QuerySkill(skillIDInt)

	updateMap := util.NewFilterMap("icon_url", "")
	// Attempt to delete record from database
	err = c.updates(skill, updateMap)
	if err != nil {
		c.Warnf("Failed to delete skill icon from database.")
		return errors.NoSuchIDError(fmt.Errorf(
			"unable to remove icon url form skill %s", skillID))
	}

	c.Printf("SkillIcon Deleted with ID: %s", skillID)
	return nil
}

// Creates new SkillIcon in database for POST requests to "/skillicons"
func (c *SkillIconsController) addSkillIcon() error {
	// Extract icon data from HTTP request
	iconFile, _, err := c.r.FormFile("icon")
	if err != nil {
		c.Warn("error getting icon form file: " + err.Error())
		return errors.ReadError(fmt.Errorf("Failed to parse icon field: %s", err))
	}
	defer iconFile.Close()

	skillValue := c.r.FormValue("skill_id")
	skillID, err := util.StringToID(skillValue)
	if err != nil {
		return err
	}
	skill := model.QuerySkill(skillID)

	// Capture data for later use before it is consumed by util.ValidateIcon
	iconFileBytes, _ := ioutil.ReadAll(iconFile)

	// Validity and error checking
	dataCopy := make([]byte, len(iconFileBytes))
	copy(dataCopy, iconFileBytes)
	_, err = util.ValidateIcon(bytes.NewReader(dataCopy))
	if err != nil {
		c.Warn("Invalid image data: ", err)
		return errors.InvalidPOSTBodyError(err)
	}
	err = c.first(&skill)
	if err != nil {
		c.Warn("ID does not exist: ", err.Error())
		return errors.InvalidPOSTBodyError(fmt.Errorf(
			"The %q field must contain ID of existing Skill in database", "skill_id"))
	}

	// Upload image to S3 cloud
	url, err := c.fileSystem.Write("dev/"+skillValue,
		bytes.NewReader(iconFileBytes))
	if err != nil {
		return fmt.Errorf("failed to save icon: %s", err)
	}

	updateMap := util.NewFilterMap("icon_url", url)
	err = c.updates(&skill, updateMap)
	if err != nil {
		c.Warnf("Update error: %v", err)
		return errors.SavingError(err)
	}

	b, err := json.Marshal(skill)
	if err != nil {
		return errors.MarshalingError(err)
	}
	c.w.Write(b)

	c.Printf("Saved icon: %s", url)
	return nil
}
