package v1

import (
	"ginLearn/models"
	"ginLearn/pkg/e"
	"ginLearn/pkg/setting"
	"ginLearn/pkg/util"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
)

//@Summary 根据名字获取标签
//@Produce json
//@Param name query int true "name"
//@Success 200 {string} gin.H "{"code":200,"data":{},"msg":"ok"}"
//@Router /api/v1/tags [get]
func GetTags(c *gin.Context) {
	name := c.Query("name")
	maps := make(map[string]interface{})
	data := make(map[string]interface{})
	if name != "" {
		maps["name"] = name
	}
	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
	}
	code := e.SUCCESS
	data["lists"] = models.GetTags(util.GetPage(c), setting.PageSize, maps)
	data["total"] = models.GetTagTotal(maps)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})

}

//@Summary 新增标签
//@Produce json
//@Param name query string true "name"
//@Param state query int true "state"
//@Param created_by query string true "createdBy"
//@Success 200 {string} gin.H "{"code":200,"data":{},"msg":"ok"}"
//@Router /api/v1/tags [post]
func AddTags(c *gin.Context) {
	name := c.Query("name")
	state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()
	createdBy := c.Query("created_by")
	valid := validation.Validation{}
	valid.Required(name, "name").Message("名称不能为空")
	valid.MaxSize(name, 100, "name").Message("名字最长100字符")
	valid.Required(createdBy, "create_by").Message("创建人不能为空")
	valid.MaxSize(createdBy, 100, "created_by").Message("创建人最长100字符")
	valid.Range(state, 0, 1, "state").Message("状态只有0,1")
	code := e.INVALID_PARAMS
	if ! valid.HasErrors() {
		if !models.ExistTagByName(name) {
			code = e.SUCCESS
			models.AddTag(name, state, createdBy)
		} else {
			code = e.ERROR_EXIST_TAG
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}

//@Summary 编辑标签
//@Produce json
//@Param id query int true "id"
//@Param name query string true "name"
//@Param modified_by query string true "modifiedBy"
//@Success 200 {string} gin.H "{"code":200,"data":{},"msg":"ok"}"
//@Router /api/v1/tags [put]
func EditTags(c *gin.Context) {
	name := c.Query("name")
	id := com.StrTo(c.Query("id")).MustInt()
	modifiedBy := c.Query("modified_by")
	valid := validation.Validation{}
	state := -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}
	valid.Required(name, "name").Message("名称不能为空")
	valid.MaxSize(name, 100, "name").Message("名字最长100字符")
	valid.Required(modifiedBy, "create_by").Message("创建人不能为空")
	valid.MaxSize(modifiedBy, 100, "created_by").Message("创建人最长100字符")
	code := e.INVALID_PARAMS
	if ! valid.HasErrors() {
		if models.ExistTagByID(id) {
			code = e.SUCCESS
			data := make(map[string]interface{})
			data["modified_by"] = modifiedBy
			if name != "" {
				data["name"] = name
			}
			if state != -1 {
				data["state"] = state
			}
			models.EditTag(id, data)
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})

}

//@Summary 删除标签
//@Produce json
//@Param id query int true "id"
//@Success 200 {string} gin.H "{"code":200,"data":{},"msg":"ok"}"
//@Router /api/v1/tags [delete]
func DeleteTags(c *gin.Context) {
	id := com.StrTo(c.Query("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")
	code := e.INVALID_PARAMS
	if ! valid.HasErrors() {
		if models.ExistTagByID(id) {
			code = e.SUCCESS
			models.DeleteTag(id)
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}
