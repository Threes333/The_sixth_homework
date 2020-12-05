package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)
type LoginUser struct {
	User           string `form:"user" json:"user"`//用户名
	Username       string `form:"username" json:"username" binding:"required"`//账号
	Password       string `form:"password" json:"password" binding:"required"`
}
type RegisterUser struct {
	LoginUser
	RepeatPassword string `form:"password2" json:"password2" binding:"required"`
}
var db *sql.DB
func InitDb()(err error){
	dsn := "root:qazpl.123456@tcp(127.0.0.1:3306)/threes"
	db,err = sql.Open("mysql",dsn)
	if err != nil{
		return err
	}
	if err = db.Ping();err != nil{
		return err
	}
	return nil
}
func QueryOne(username interface{},user *LoginUser)  {
	sqlStr := "Select username, user, password from userdata where username = ?"
	err := db.QueryRow(sqlStr,username).Scan(&user.Username, &user.User, &user.Password)
	if err != nil{
		log.Println(err.Error())
		return
	}
}
func Insert(user RegisterUser) {
	str := "Insert into userdata (username, user, password) value(?,?,?)"
	stmt,err := db.Prepare(str)
	if err != nil{
		log.Println(err.Error())
		return
	}
	_,err = stmt.Exec(user.Username,user.User, user.Password)
	if err != nil{
		log.Println(err.Error())
		return
	}
}
func Update(newValue,username interface{})  {
	str := "Update userdata set password = ? where username = ?"
	stmt,err := db.Prepare(str)
	if err != nil{
		log.Println(err.Error())
		return
	}
	fmt.Println(newValue,username)
	_,err = stmt.Exec(newValue,username)
	if err != nil{
		log.Println(err.Error())
		return
	}
}
func check(user RegisterUser)(problem int){
	var val interface{}
	problem = 0
	str1 := "select user from userdata where user = ?"
	str2 := "select user from userdata where username = ?"
	if err := db.QueryRow(str1,user.User).Scan(val);err != sql.ErrNoRows{
		//前端显示用户名已被注册
		problem += 1
	}
	if err := db.QueryRow(str2,user.Username).Scan(val);err != sql.ErrNoRows{
		//前端显示账号已被注册
		problem += 2
	}
	return
}
func main()  {
	r := gin.Default()
	r.LoadHTMLGlob("./Templates/*")
	err := InitDb()
	if err != nil{
		log.Println("Create DB failed",err.Error())
		return
	}else {
		fmt.Println("成功连接数据库!")
	}
	r.GET("/login", LoginHome())
	r.POST("/login", HandleLogin(r))
	r.GET("/register", RegisterHome())
	r.POST("/register", HandleRegister())
	r.POST("/SetPassword", SetPassword())
	r.GET("/SetPassword", func(c *gin.Context) {
		c.HTML(http.StatusOK, "SetPassword.html", nil)
	})
	r.GET("/SetSignature", SetSignatureHome())
	r.POST("/SetSignature", HandleSetSignature())
	r.POST("/joined/:id", PersonalHome())
	_ = r.Run()
}

func PersonalHome() func(c *gin.Context) {
	return func(c *gin.Context) {
		var user LoginUser
		username := c.Param("id")
		if username == "" {
			c.HTML(http.StatusOK, "index.html", nil)
			return
		}
		fmt.Println("joined", username)
		QueryOne(username, &user)
		row := db.QueryRow("select signature from userdata where username = ?", username)
		var signature string
		err := row.Scan(&signature)
		if err != nil {
			log.Println(err.Error())
			c.HTML(http.StatusOK, "index.html", nil)
		} else {
			c.HTML(http.StatusOK, "home.html", gin.H{
				"user":      user.User,
				"signature": signature,
				"username":  user.Username,
			})
		}
	}
}

func HandleSetSignature() func(c *gin.Context) {
	return func(c *gin.Context) {
		var signature string
		username := c.Query("id")
		fmt.Println("signature", username)
		signature = c.PostForm("signature")
		_, err := db.Exec("UPDATE userdata set signature = ? where username = ?", signature, username)
		if err != nil {
			log.Println(err.Error())
			c.HTML(http.StatusOK, "login.html", nil)
		} else {
			c.HTML(http.StatusOK, "signature.html", gin.H{
				"username": username,
				"sign":     "提交成功!",
			})
		}
	}
}

func SetSignatureHome() func(c *gin.Context) {
	return func(c *gin.Context) {
		username := c.Query("id")
		c.HTML(http.StatusOK, "signature.html", gin.H{
			"username": username,
			"sign":     "",
		})
	}
}

func SetPassword() func(c *gin.Context) {
	return func(c *gin.Context) {
		password := c.PostForm("password")
		username := c.PostForm("username")
		NewPassword := c.PostForm("NewPassword")
		fmt.Println(password, username, NewPassword)
		var user LoginUser
		QueryOne(username, &user)
		if user.Username != username {
			c.HTML(http.StatusOK, "SetPassword.html", "账号不存在!")
		} else if user.Password != password {
			c.HTML(http.StatusOK, "SetPassword.html", "原密码错误!")
		} else {
			Update(NewPassword, username)
			c.HTML(http.StatusOK, "manageSet.html", user.Username)
		}
	}
}

func HandleLogin(r *gin.Engine) func(c *gin.Context) {
	return func(c *gin.Context) {
		var user LoginUser
		err := c.ShouldBind(&user)
		if err != nil {
			fmt.Println(err.Error())
			c.HTML(http.StatusOK, "index.html", "账号或密码为空")
		} else {
			var InUser LoginUser
			QueryOne(user.Username, &InUser)
			if InUser.Username != user.Username {
				fmt.Println(user.Username, " != ", InUser.Username)
				c.HTML(http.StatusOK, "index.html", "账号不存在!")
				//前端显示账号不存在
			} else if InUser.Password != user.Password {
				//前端显示密码错误
				fmt.Println(user.Password, " != ", InUser.Password)
				c.HTML(http.StatusOK, "index.html", "密码错误!")
			} else {
				c.Request.URL.Path = fmt.Sprintf("/joined/%s",user.Username)
				r.HandleContext(c)
			}
		}
	}
}

func LoginHome() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	}
}

func HandleRegister() func(c *gin.Context) {
	return func(c *gin.Context) {
		var user RegisterUser
		if err := c.ShouldBind(&user); err != nil {
			c.HTML(http.StatusOK, "register.html", "输入不能为空!")
			log.Println(err.Error())
		} else {
			switch check(user) {
			case 0:
				ok := user.Password != user.RepeatPassword
				if ok {
					c.HTML(http.StatusOK, "register.html", "两次输入的密码不一致!")
				} else {
					Insert(user)
					c.HTML(http.StatusOK, "manageRegister.html", nil)
				}
			case 1:
				c.HTML(http.StatusOK, "register.html", "用户名已被注册!")
			case 2:
				c.HTML(http.StatusOK, "register.html", "账号已被注册!")
			case 3:
				c.HTML(http.StatusOK, "register.html", "用户名和账号都已被注册!")
			}
		}
	}
}

func RegisterHome() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", nil)
	}
}
