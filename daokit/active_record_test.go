package daokit

import "fmt"

func check() {
	entity := M("auth_config")
	i, err := entity.Set("app_id", "x23").Set("key", "y").Set("token", nil).Insert()
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Println(i)
}
