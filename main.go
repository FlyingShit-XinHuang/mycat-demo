package main

import (
	"fmt"
	"whispir/mycat-demo/models/template"

	"os"
	"whispir/mycat-demo/models"
	spacemodel "whispir/mycat-demo/models/space"
)

func main() {
	err := models.Connect("demo", "demo", "localhost:8066", "demo")
	if nil != err {
		panic(err)
	}

	// insert
	space := spacemodel.NewSpace("habor-demo")
	if _, err := models.Create(space); nil != err {
		fmt.Println("Failed to insert space %v:", space, err)
		os.Exit(1)
	}
	fmt.Println("New space is created with ", space.Id)

	tmpl := template.NewTemplate("habor-demo", "this is a message template of habor", space)
	if _, err := models.Create(tmpl); nil != err {
		fmt.Println("Failed to insert template %v:", tmpl, err)
		os.Exit(1)
	}
	fmt.Println("New template is created with ", tmpl.Id)

	newSpace := spacemodel.Space{}
	if err := models.FindByPK(&newSpace, space.Id); nil != err {
		fmt.Printf("Failed to query space %d: %v\n", space.Id, err)
		os.Exit(1)
	}
	fmt.Println("Selected space:", newSpace)

	newTmpl := template.Template{}
	if err := models.FindByPK(&newTmpl, tmpl.Id); nil != err {
		fmt.Printf("Failed to query template %d: %v\n", tmpl.Id, err)
		os.Exit(1)
	}
	fmt.Println("Selected template", newTmpl)

	// update
	if _, err := models.Update(
		spacemodel.Space{
			Name: "modified demo",
		}, "id=?", space.Id); nil != err {

		fmt.Printf("Failed to update space %d: %v\n", space.Id, err)
		os.Exit(1)
	}

	if _, err := models.Update(
		template.Template{
			Name:    "modified demo template",
			Content: []byte("this template was modified"),
		}, "id=? and name=?", tmpl.Id, tmpl.Name); nil != err {

		fmt.Printf("Failed to update template %d: %v\n", tmpl.Id, err)
		os.Exit(1)
	}

	if err := models.FindByPK(&newSpace, space.Id); nil != err {
		fmt.Println("Failed to query updated space:", err)
		os.Exit(1)
	}
	fmt.Println("Updated space:", newSpace)

	if err := models.FindByPK(&newTmpl, tmpl.Id); nil != err {
		fmt.Println("Failed to query updated template:", err)
		os.Exit(1)
	}
	fmt.Println("Updated tempate:", newTmpl)

	// customized query
	if _, err := models.Create(
		template.NewTemplate("habor-demo2", "this is another message template of habor", space)); nil != err {

		fmt.Println("Failed to insert 2nd template:", err)
		os.Exit(1)
	}

	list, err := template.ListInSpace(space)
	if nil != err {
		fmt.Printf("Failed to get all templates in space %d: %v\n", space.Id, err)
		os.Exit(1)
	}
	fmt.Printf("There are %d templates in space %d\n", len(list), space.Id)

	//soft delete
	models.Delete(&spacemodel.Space{}, "id=?", space.Id) // It's ok because of 'soft delete'
	models.Delete(&template.Template{}, "space_id=?", space.Id)

	list, err = template.ListInSpace(space)
	if nil != err {
		fmt.Printf("Failed to get all templates in space %d: %v\n", space.Id, err)
		os.Exit(1)
	}
	fmt.Printf("There are %d templates in space %d after soft deleted\n", len(list), space.Id)
}
