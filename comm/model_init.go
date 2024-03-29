/**
* @Author: xhzhang
* @Date: 2019/7/15 15:19
 */
package comm

func checkTableAndCreate(obj interface{}) bool {
	if !DB.HasTable(obj) {
		DB.CreateTable(obj)
		return false
	}
	return true
}

func inittable() {
	//Organization
	if !checkTableAndCreate(Organization{}) {
		orgname := "cdporg"
		org := Organization{Name: orgname}
		DB.Create(&org)
	}
	//Environment
	if !checkTableAndCreate(Environment{}) {
		envname := "cdpenv"
		env := Environment{Name: envname}
		DB.Create(&env)
	}
	//Project
	if !checkTableAndCreate(Project{}) {
		proname := "cdppro"
		pro := Project{Name: proname}
		DB.Create(&pro)
	}
	checkTableAndCreate(Agent{})
	//if !checkTableAndCreate(Agent_Operation{}) {
	//	DB.Model(&Agent_Operation{}).AddForeignKey("agent_id", "cdp_agents(id)", "CASCADE", "CASCADE")
	//}

	//group
	if !checkTableAndCreate(Group{}) {
		//DB.Model(&Group{}).AddForeignKey("organization_id", "cdp_organizations(id)", "CASCADE", "CASCADE")
		//DB.Model(&Group{}).AddForeignKey("environment_id", "cdp_environments(id)", "CASCADE", "CASCADE")
		//DB.Model(&Group{}).AddForeignKey("project_id", "cdp_projects(id)", "CASCADE", "CASCADE")
		groname := "cdpgro"
		gro := Group{Name: groname, OrganizationID: 1, EnvironmentID: 1, ProjectID: 1}
		DB.Create(&gro)
	}
	//service
	checkTableAndCreate(Service{})
	//if !checkTableAndCreate(Service{}) {
	//	DB.Model(&Service{}).AddForeignKey("agent_id", "cdp_agents(id)", "CASCADE", "CASCADE")
	//	DB.Model(&Service{}).AddForeignKey("group_id", "cdp_groups(id)", "CASCADE", "CASCADE")
	//}

	//release

	if !checkTableAndCreate(Release{}) {
		//DB.Model(&Release{}).AddForeignKey("organization_id", "cdp_organizations(id)", "CASCADE", "CASCADE")
		//DB.Model(&Release{}).AddForeignKey("project_id", "cdp_projects(id)", "CASCADE", "CASCADE")
		DB.Model(&Release{}).AddUniqueIndex("release_unique","version","project_id","organization_id")
	}

	//release-code
	checkTableAndCreate(ReleaseCode{})
	//if !checkTableAndCreate(ReleaseCode{}) {
	//	DB.Model(&ReleaseCode{}).AddForeignKey("release_id", "cdp_releases(id)", "CASCADE", "CASCADE")
	//	DB.Model(&ReleaseCode{}).AddUniqueIndex("releasecode_u","name","release_id")
	//}

	// task
	checkTableAndCreate(Task{})

	// execution
	checkTableAndCreate(Execution{})
	//if !checkTableAndCreate(Execution{}) {
	//	DB.Model(&Execution{}).AddForeignKey("service_id", "cdp_services(id)", "CASCADE", "CASCADE")
	//	DB.Model(&Execution{}).AddForeignKey("task_id", "cdp_tasks(id)", "CASCADE", "CASCADE")
	//}
	// execution_detail
	checkTableAndCreate(Execution_Detail{})
	//if !checkTableAndCreate(Execution_Detail{}) {
	//	DB.Model(&Execution_Detail{}).AddForeignKey("execution_id", "cdp_executions(id)", "CASCADE", "CASCADE")
	//}

	// task_timed
	checkTableAndCreate(Cron_Task{})
	//if !checkTableAndCreate(Cron_Task{}) {
	//	DB.Model(&Cron_Task{}).AddForeignKey("task_id", "cdp_tasks(id)", "CASCADE", "CASCADE")
	//}

}
