package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)
var db *sql.DB
func InitDb(){
	defer func() {
		if err := recover();err != nil{
			fmt.Println("Other wrong:",err)
		}
	}()
	dsn := "root:qazpl.123456@tcp(127.0.0.1:3306)/threes"
	var err error
	db ,err = sql.Open("mysql",dsn)
	if err != nil{
		fmt.Println("open:",err.Error())
		return
	}
	errs := db.Ping()
	if errs != nil{
		fmt.Println("ping:",errs.Error())
		return
	}
	fmt.Println("连接数据库成功!")
}
type User struct{
	id int
	name string
	password string
	num int
}
type updateData struct {
	id int
	index int
	data interface{}
	newData interface{}
}
func QueryOne(n int,users *User) error {
	str := "select * from user where id = ?;"
	err := db.QueryRow(str, n).Scan(&users.id, &users.name, &users.password, &users.num)
	if err != nil{
		return err
	}
	return nil
}
func QueryOnes(m,n int,user []User) error {
	str := "select * from user limit ?,?"
	rows, err := db.Query(str, m, n)
	if err != nil {
		return err
	}
	i := 0
	for rows.Next(){
		_ = rows.Scan(&user[i].id, &user[i].name, &user[i].password, &user[i].num)
		i++
	}
	return nil
}
func Insert(user User) error  {
	str := "insert into user value(?,?,?,?)"
	_, err := db.Exec(str, user.id, user.name, user.password, user.num)
	if err != nil{
		return err
	}
	fmt.Println("成功插入1条数据")
	return nil
}
func Update(id,index int, data,newData interface{}) error  {
	str := "Update user set ? = ? where ? = ?"
	_, err := db.Exec(str,data,newData,id,index)
	if err != nil{
		return err
	}
	fmt.Println("成功更改1条数据")
	return nil
}
func Delete(index int) error {
	str := "Delete from user where id = ?"
	_, err := db.Exec(str,index)
	if err != nil{
		return err
	}
	fmt.Println("成功删除1条数据")
	return nil
}
func prepareInserts(user []User) error  {
	str := "Insert into user value(?,?,?,?)"
	stmt,err := db.Prepare(str)
	if err != nil{
		return err
	}
	var i = 0
	for _,v := range user{
		_, _ = stmt.Exec(v.id, v.name, v.password, v.num)
		i++
	}
	fmt.Printf("成功插入%d条数据",i)
	return nil
}
func prepareUpdates(data []updateData) error {
	str := "Update user set ? = ? where ? = ?"
	stmt,err := db.Prepare(str)
	if err != nil{
		return err
	}
	var i = 0
	for _,v := range data{
		_, _ = stmt.Exec(v.data, v.newData, v.id, v.index)
		i++
	}
	fmt.Printf("成功更改%d条数据",i)
	return nil
}
func prepareDeletes(index []int) error {
	str := "Delete from user where id = ?"
	stmt,err := db.Prepare(str)
	if err != nil{
		return err
	}
	var i = 0
	for _,v := range index{
		_, _ = stmt.Exec(v)
		i++
	}
	fmt.Printf("成功删除%d条数据",i)
	return nil
}
func Affair(str1 ,str2 string) error {
	tx,err := db.Begin()
	if err != nil {
		return err
	}
	_,err = tx.Exec(str1)
	if err != nil {
		fmt.Println("操作一失败,进行回卷")
		_ = tx.Rollback()
		return err
	}
	_,err = tx.Exec(str2)
	if err != nil{
		fmt.Println("操作二失败,进行回卷")
		_ = tx.Rollback()
		return err
	}
	_ = tx.Commit()
	fmt.Println("事务成功完成!")
	return nil
}
func main () {
	var user User
	users := make([]User,5)
	InitDb()
	err := QueryOne(3,&user)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(user)
	m,n := 0,4
	err = QueryOnes(m,n,users)
	if err != nil{
		fmt.Println(err.Error())
	}
	for i:=0;i<n;i++{
		fmt.Println(users[i])
	}
	fmt.Println("------分割线------")
	for _,v := range users{
		fmt.Println(v)
	}
	/*var adduser User = User{
		id : 5,
		name : "艾嚄",
		password : "897",
		num : 33,
	}
	err1 := Insert(adduser)
	if err1 != nil{
		fmt.Println(err1.Error())
		return
	}*/
	err = QueryOnes(1,4,users)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("------分割线------")
	for _,v := range users{
		fmt.Println(v)
	}
	str1 := "update user set num = num - 2 where id = 3"
	str2 := "update user set num = num + 2 where id = 4"
	err = Affair(str1,str2)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = QueryOnes(0,5,users)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("------分割线------")
	for _,v := range users{
		fmt.Println(v)
	}
}
