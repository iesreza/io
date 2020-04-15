package auth

import (
	"fmt"
	"github.com/gofiber/fiber"
	"github.com/iesreza/io"
	"github.com/iesreza/io/lib/T"
	"github.com/iesreza/io/lib/jwt"
	"github.com/iesreza/io/user"
	"github.com/iesreza/validate"
	"gopkg.in/hlandau/passlib.v1"
	"net/http"
	"time"
)

type Controller struct {}
type AuthParams struct {
	Username   string `json:"username" form:"username" validate:"empty=false"`
	Password   string `json:"password" form:"password" validate:"empty=false"`
	Remember   bool   `json:"remember" form:"remember"`
	Return     string `json:"return" form:"return" validate:"empty=true | one_of=json,text,html"`
	Redirect   string `json:"redirect" form:"redirect"`
}

type Response struct {
	Success bool    `json:"success"`
	Message string  `json:"message"`
	Data    interface{} `json:"data"`
}

func GetUserByID(id interface{}) *user.User {
	user := user.User{}
	if db.Where("id = ?",id).Find(&user).RecordNotFound(){
		return nil
	}
	return &user
}

func GetUserByUsername(username interface{}) *user.User {
	user := user.User{}
	if db.Where("username = ?",username).Find(&user).RecordNotFound(){
		return nil
	}
	return &user
}

func GetUserByEmail(email interface{}) *user.User {
	user := user.User{}
	if db.Where("email = ?",email).Find(&user).RecordNotFound(){
		return nil
	}
	return &user
}


func GetGroup(v interface{}) *user.Group {
	group := user.Group{}
	if db.Where("id = ? OR code_name", v,v ).Find(&group).RecordNotFound(){
		return nil
	}
	return &group
}

func AuthUserByPassword(username,password string) (*user.User,error) {
	user := GetUserByUsername(username)
	if user == nil{
		return user,fmt.Errorf("username not found")
	}
	_,err := passlib.Verify(password,user.Password)
	if err != nil{
		return user,fmt.Errorf("password not match")
	}
	user.Seen = time.Now()
	user.Save()
	return user,nil
}

func (c Controller)Login(ctx *fiber.Ctx)  {
	var err error
	var user *user.User
	var token string

	r := io.Upgrade(ctx)
	r.Accepts("text/html","application/json")
	params := AuthParams{}
	r.BodyParser(&params)
	if params.Return == ""{
		params.Return = "text"
	}

	err = validate.Validate(&params)
	if err == nil{
		user,err = AuthUserByPassword(params.Username,params.Password)
		if err == nil {
			extend := io.GetConfig().JWT.Age
			if params.Remember {
				extend = 365 * 24 * time.Hour
			}
			token, err = jwt.Generate(map[string]interface{}{
				"id":       user.ID,
				"username": user.Username,
				"name":     user.Name,
				"seen":     user.Seen,
				"active":   user.Active,
				//"params":   user.Params,
			}, extend)
		}

	}
	switch params.Return {
	case "html":

		break
	case "json":
		if err != nil {
			r.Format(Response{
				Success:false,
				Message:err.Error(),
			})
		}else{
			r.Format(Response{
				Success:true,
				Data:token,
			})
		}
		break
	case "text":
		if err != nil {
			r.Status(http.StatusBadRequest)
			r.Write("nok ")
			r.Write(err)
		}else{
			r.Write("ok")
			r.Write(token)
		}
		break
	}

}

func (c Controller)CreateUser(ctx *fiber.Ctx)  {
	//TODO: Check permission
	var r    = io.Upgrade(ctx)
	var user = user.User{}
	err := r.BodyParser(&user)

	if err == nil {
		err := user.Save()
		if err == nil {
			ctx.Format(Response{
				Success:true,
				Data:map[string]interface{}{
					"id":user.ID,
					"username":user.Username,
					"name":user.Name,
					"param":user.Params,
				},
			})
		}else{
			r.Format(Response{
				Success: false,
				Message: err.Error(),
			})
		}
	}else{
		r.Format(Response{
			Success:false,
			Message: "unable read form",
		})
	}

}

func (c Controller)CreateRole(ctx *fiber.Ctx)  {
	//TODO: Check permission
	var r    = io.Upgrade(ctx)
	var role = user.Role{}
	err := r.BodyParser(&role)
	if err == nil {
		err := role.Save()
		if err == nil {
			ctx.Format(Response{
				Success:true,
				Data:map[string]interface{}{
					"id":role.ID,
					"name":role.Name,
					"code_name":role.CodeName,
					"parent":role.Parent,
				},
			})
		}else{
			r.Format(Response{
				Success: false,
				Message: err.Error(),
			})
		}
	}else{
		r.Format(Response{
			Success:false,
			Message: "unable read form",
		})
	}

}

func (c Controller)CreateGroup(ctx *fiber.Ctx)  {
	//TODO: Check permission
	var r    = io.Upgrade(ctx)
	var group = user.Group{}
	err := r.BodyParser(&group)
	if err == nil {
		err := group.Save()
		if err == nil {
			ctx.Format(Response{
				Success:true,
				Data:group,
			})
		}else{
			r.Format(Response{
				Success: false,
				Message: err.Error(),
			})
		}
	}else{
		r.Format(Response{
			Success:false,
			Message: "unable read form",
		})
	}

}

func (c Controller)EditUser(ctx *fiber.Ctx)  {
	//TODO: Check permission
	var r    = io.Upgrade(ctx)
	var user user.User
	var id = T.Must(r.Params("id")).UInt()
	if db.Where("id = ?",id).Find(&user).RecordNotFound(){
		r.Format(Response{
			Success: false,
			Message: "invalid id",
		})
		return
	}
	err := r.BodyParser(&user)
	user.ID = id
	if err == nil {
		err := user.Save()
		if err == nil {
			ctx.Format(Response{
				Success:true,
				Data:map[string]interface{}{
					"id":user.ID,
					"username":user.Username,
					"name":user.Name,
					"param":user.Params,
				},
			})
		}else{
			r.Format(Response{
				Success: false,
				Message: err.Error(),
			})
		}
	}else{
		r.Format(Response{
			Success:false,
			Message: "unable read form",
		})
	}

}

func (c Controller)EditRole(ctx *fiber.Ctx)  {
	//TODO: Check permission
	var r    = io.Upgrade(ctx)
	var role user.Role
	var id = r.Params("id")
	if db.Where("id = ? OR code_name = ?",id,id).Find(&role).RecordNotFound(){
		r.Format(Response{
			Success: false,
			Message: "invalid id",
		})
		return
	}
	rid := role.ID
	err := r.BodyParser(&role)
	role.ID = rid
	if err == nil {
		err := role.Save()
		if err == nil {
			ctx.Format(Response{
				Success:true,
				Data:role,
			})
		}else{
			r.Format(Response{
				Success: false,
				Message: err.Error(),
			})
		}
	}else{
		r.Format(Response{
			Success:false,
			Message: "unable read form",
		})
	}

}

func (c Controller)EditGroup(ctx *fiber.Ctx)  {
	//TODO: Check permission
	var r    = io.Upgrade(ctx)
	var group user.Group
	var id = r.Params("id")
	if db.Where("id = ? OR code_name = ?",id,id).Find(&group).RecordNotFound(){
		r.Format(Response{
			Success: false,
			Message: "invalid id",
		})
		return
	}
	var gid = group.ID
	err := r.BodyParser(&group)
	group.ID = gid
	if err == nil {
		err := group.Save()
		if err == nil {
			ctx.Format(Response{
				Success:true,
				Data:group,
			})
		}else{
			r.Format(Response{
				Success: false,
				Message: err.Error(),
			})
		}
	}else{
		r.Format(Response{
			Success:false,
			Message: "unable read form",
		})
	}

}

func (c Controller)GetGroups(ctx *fiber.Ctx) {
	var groups []user.Group
	err := db.Find(&groups).Error
	if err != nil {
		ctx.Format(Response{
			Success: false,
			Message: err.Error(),
		})
	}else{
		ctx.Format(Response{
			Success: true,
			Message: "",
			Data:groups,
		})
	}
}

func (c Controller)GetGroup(ctx *fiber.Ctx) {
	var group user.Group
	var id = ctx.Params("id")
	if db.Where("id = ? OR code_name = ?",id,id).Find(&group).RecordNotFound(){
		ctx.Format(Response{
			Success: false,
			Message: "invalid id",
		})
		return
	}else{
		ctx.Format(Response{
			Success: true,
			Message: "",
			Data:group,
		})
	}
}


func (c Controller)GetRoles(ctx *fiber.Ctx) {
	var roles []user.Role
	err := db.Find(&roles).Error
	if err != nil {
		ctx.Format(Response{
			Success: false,
			Message: err.Error(),
		})
	}else{
		ctx.Format(Response{
			Success: true,
			Message: "",
			Data:roles,
		})
	}
}

func (c Controller)GetRole(ctx *fiber.Ctx) {
	var role user.Role
	var id = ctx.Params("id")
	if db.Where("id = ? OR code_name = ?",id,id).Find(&role).RecordNotFound(){
		ctx.Format(Response{
			Success: false,
			Message: "invalid id",
		})
		return
	}else{
		ctx.Format(Response{
			Success: true,
			Message: "",
			Data:role,
		})
	}
}

func (c Controller)GetRoleGroups(ctx *fiber.Ctx) {
	var role user.Role
	var id = ctx.Params("id")
	if db.Where("id = ? OR code_name = ?",id,id).Find(&role).RecordNotFound(){
		ctx.Format(Response{
			Success: false,
			Message: "invalid id",
		})
		return
	}else{
		var groups []user.Group
		db.Joins(`INNER JOIN group_roles ON "group".id = group_roles.group_id`).Where("group_roles.role_id = ?",role.ID).Find(&groups)
		ctx.Format(Response{
			Success: true,
			Message: "",
			Data:groups,
		})
	}
}

func (c Controller)GetUser(ctx *fiber.Ctx) {
	var user user.User
	var id = ctx.Params("id")
	if db.Where("id = ? OR username = ? OR email = ?",id,id,id).Find(&user).RecordNotFound(){
		ctx.Format(Response{
			Success: false,
			Message: "invalid id",
		})
		return
	}else{
		ctx.Format(Response{
			Success: true,
			Message: "",
			Data:user,
		})
	}
}

func (c Controller)GetMe(ctx *fiber.Ctx) {
	r := io.Upgrade(ctx)
	ctx.Format(Response{
		Success: true,
		Message: "",
		Data:r.User,
	})
}

func (c Controller)GetAllUsers(ctx *fiber.Ctx) {
	var users []user.User
	err := db.Offset(ctx.Params("offset")).Limit(ctx.Params("limit")).Find(&users).Error
	if err != nil {
		ctx.Format(Response{
			Success: false,
			Message: err.Error(),
		})
	}else{
		ctx.Format(Response{
			Success: true,
			Message: "",
			Data:users,
		})
	}
}

func (c Controller)GetAllPermissions(ctx *fiber.Ctx)  {
	var perms []user.Permission
	err := db.Find(&perms).Error
	if err != nil {
		ctx.Format(Response{
			Success: false,
			Message: err.Error(),
		})
	}else{
		ctx.Format(Response{
			Success: true,
			Message: "",
			Data:perms,
		})
	}
}