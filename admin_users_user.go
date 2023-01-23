package cms

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gouniverse/api"
	"github.com/gouniverse/bs"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/utils"
)

func (cms Cms) pageUsersUserManager(w http.ResponseWriter, r *http.Request) {
	endpoint := r.Context().Value(keyEndpoint).(string)
	// log.Println(endpoint)

	users, err := cms.UserStore.EntityList("user", 0, 200, "", "id", "asc")

	if err != nil {
		api.Respond(w, r, api.Error("User list failed to be retrieved "+err.Error()))
		return
	}

	header := cms.cmsHeader(endpoint)
	breadcrums := cms.cmsBreadcrumbs(map[string]string{
		endpoint: "Home",
		(endpoint + "?path=" + PathUsersUserManager): "Users",
	})

	container := hb.NewDiv().Attr("class", "container").Attr("id", "user-manager")
	heading := hb.NewHeading1().HTML("User Manager")
	button := hb.NewButton().HTML("New user").Attr("class", "btn btn-success float-end").Attr("v-on:click", "showUserCreateModal")
	heading.AddChild(button)

	container.Children([]*hb.Tag{
		hb.NewHTML(header),
		heading,
		hb.NewHTML(breadcrums),
		pageUsersUserCreateModal(),
		pageUsersUserTrashModal(),
	})

	table := hb.NewTable().Attr("id", "TableUsers").Attr("class", "table table-responsive table-striped mt-3")
	thead := hb.NewThead()
	tbody := hb.NewTbody()
	table.AddChild(thead).AddChild(tbody)

	tr := hb.NewTR()
	th1 := hb.NewTD().HTML("Name")
	th2 := hb.NewTD().HTML("Status")
	th3 := hb.NewTD().HTML("Actions").Attr("style", "width:150px;")
	thead.AddChild(tr.AddChild(th1).AddChild(th2).AddChild(th3))

	for _, user := range users {
		firstName, _ := user.GetString("first_name", "n/a")
		lastName, _ := user.GetString("last_name", "n/a")
		status, _ := user.GetString("status", "n/a")
		buttonEdit := hb.NewButton().HTML("Edit").Attr("type", "button").Attr("class", "btn btn-primary btn-sm").Attr("v-on:click", "userEdit('"+user.ID+"')").Attr("style", "margin-right:5px")
		buttonTrash := hb.NewButton().HTML("Trash").Attr("type", "button").Attr("class", "btn btn-danger btn-sm").Attr("v-on:click", "showUserTrashModal('"+user.ID+"')")

		tr := hb.NewTR()
		td1 := hb.NewTD().HTML(firstName + " " + lastName)
		td2 := hb.NewTD().HTML(status)
		td3 := hb.NewTD().AddChild(buttonEdit).AddChild(buttonTrash)

		tbody.AddChild(tr.AddChild(td1).AddChild(td2).AddChild(td3))
	}
	container.AddChild(table)

	h := container.ToHTML()

	inlineScript := `
var userCreateUrl = "` + endpoint + `?path=users/user-create-ajax"
var userTrashUrl = "` + endpoint + `?path=users/user-trash-ajax"
var userUpdateUrl = "` + endpoint + `?path=users/user-update"
const UserManager = {
	data() {
		return {
		  userCreateModel:{
			  firstName:"",
			  lastName:"",
		  },
		  userTrashModel:{
			id:""
		  }
		}
	},
	created(){
		this.initDataTable();
	},
	methods: {
		initDataTable(){
			$(() => {
				$('#TableUsers').DataTable({
					"order": [[ 0, "asc" ]] // 1st column
				});
			});
		},
        showUserCreateModal(){
			//alert("Create user");
			var modalUserCreate = new bootstrap.Modal(document.getElementById('ModalUserCreate'));
			modalUserCreate.show();
		},
		userCreate(){
			const firstName = this.userCreateModel.firstName;
			const lastName = this.userCreateModel.lastName;
			$.post(userCreateUrl, {first_name: firstName, last_name: lastName}).done((result)=>{
				if (result.status==="success"){
					var modalUserCreate = new bootstrap.Modal(document.getElementById('ModalUserCreate'));
					modalUserCreate.hide();

					return location.href = userUpdateUrl+ "&user_id=" + result.data.user_id;
				}
				alert("Failed. " + result.message)
			}).fail((result)=>{
				alert("Failed" + result)
			});
		},
		userEdit(userId){
			return location.href = userUpdateUrl+ "&user_id=" + userId;
		},
		showUserTrashModal(userId){
			this.userTrashModel.id = userId;
			var modalUserTrash = new bootstrap.Modal(document.getElementById('ModalUserTrash'));
			modalUserTrash.show();
		},
		userTrash(){
            let userId = this.userTrashModel.id;
			$.post(userTrashUrl, {user_id: userId}).done((result)=>{
				if (result.status==="success"){
					var ModalUserTrash = new bootstrap.Modal(document.getElementById('ModalUserTrash'));
				    ModalUserTrash.hide();
					location.href = location.href;
					return;
				}
				alert("Failed. " + result.message)
			}).fail((result)=>{
				alert("Failed" + result)
			});
		}
	}
};
Vue.createApp(UserManager).mount('#user-manager')
	`

	webpage := Webpage("User Manager", h)
	webpage.AddStyleURL("https://cdnjs.cloudflare.com/ajax/libs/datatables/1.10.21/css/jquery.dataTables.css")
	webpage.AddScriptURL("https://cdnjs.cloudflare.com/ajax/libs/datatables/1.10.21/js/jquery.dataTables.js")
	webpage.AddScript(inlineScript)
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(webpage.ToHTML()))
}

func pageUsersUserTrashModal() *hb.Tag {
	modal := hb.NewDiv().Attr("id", "ModalUserTrash").Attr("class", "modal fade")
	modalDialog := hb.NewDiv().Attr("class", "modal-dialog")
	modalContent := hb.NewDiv().Attr("class", "modal-content")
	modalHeader := hb.NewDiv().Attr("class", "modal-header").AddChild(hb.NewHeading5().HTML("Trash Template"))
	modalBody := hb.NewDiv().Attr("class", "modal-body")
	modalBody.AddChild(hb.NewParagraph().HTML("Are you sure you want to move this user to trash bin?"))
	modalFooter := hb.NewDiv().Attr("class", "modal-footer")
	modalFooter.AddChild(hb.NewButton().HTML("Close").Attr("class", "btn btn-secondary").Attr("data-bs-dismiss", "modal"))
	modalFooter.AddChild(hb.NewButton().HTML("Move to trash bin").Attr("class", "btn btn-danger").Attr("v-on:click", "userTrash"))
	modalContent.AddChild(modalHeader).AddChild(modalBody).AddChild(modalFooter)
	modalDialog.AddChild(modalContent)
	modal.AddChild(modalDialog)
	return modal
}

func pageUsersUserCreateModal() *hb.Tag {
	modal := hb.NewDiv().Attr("id", "ModalUserCreate").Attr("class", "modal fade")
	modalDialog := hb.NewDiv().Attr("class", "modal-dialog")
	modalContent := hb.NewDiv().Attr("class", "modal-content")
	modalHeader := hb.NewDiv().Attr("class", "modal-header").Children([]*hb.Tag{
		hb.NewHeading5().HTML("New User"),
	})
	modalBody := hb.NewDiv().Attr("class", "modal-body").Children([]*hb.Tag{
		hb.NewDiv().Attr("class", "form-group").Children([]*hb.Tag{
			hb.NewLabel().HTML("First Name"),
			hb.NewInput().Attr("class", "form-control").Attr("v-model", "userCreateModel.firstName"),
		}),
		hb.NewDiv().Attr("class", "form-group").Children([]*hb.Tag{
			hb.NewLabel().HTML("Last Name"),
			hb.NewInput().Attr("class", "form-control").Attr("v-model", "userCreateModel.lastName"),
		}),
	})
	modalFooter := hb.NewDiv().Attr("class", "modal-footer")
	modalFooter.AddChild(hb.NewButton().HTML("Close").Attr("class", "btn btn-prsecondary").Attr("data-bs-dismiss", "modal"))
	modalFooter.AddChild(hb.NewButton().HTML("Create & Continue").Attr("class", "btn btn-primary").Attr("v-on:click", "userCreate"))
	modalContent.AddChild(modalHeader).AddChild(modalBody).AddChild(modalFooter)
	modalDialog.AddChild(modalContent)
	modal.AddChild(modalDialog)
	return modal
}

// pageUsersUserTrashAjax - moves the user to the trash
func (cms Cms) pageUsersUserTrashAjax(w http.ResponseWriter, r *http.Request) {
	userID := strings.Trim(utils.Req(r, "user_id", ""), " ")

	if userID == "" {
		api.Respond(w, r, api.Error("User ID is required"))
		return
	}

	user, _ := cms.UserStore.EntityFindByID(userID)

	if user == nil {
		api.Respond(w, r, api.Error("User NOT FOUND with ID "+userID))
		return
	}

	isOk, err := cms.UserStore.EntityTrash(userID)

	if err != nil {
		api.Respond(w, r, api.Error("User failed to be trashed "+err.Error()))
		return
	}

	if !isOk {
		api.Respond(w, r, api.Error("User failed to be moved to trash"))
		return
	}

	api.Respond(w, r, api.SuccessWithData("User moved to trash successfully", map[string]interface{}{
		"user_id": user.ID,
	}))
}

func (cms Cms) pageUsersUserCreateAjax(w http.ResponseWriter, r *http.Request) {
	firstName := strings.Trim(utils.Req(r, "first_name", ""), " ")
	lastName := strings.Trim(utils.Req(r, "last_name", ""), " ")

	if firstName == "" {
		api.Respond(w, r, api.Error("first name is required field"))
		return
	}

	if lastName == "" {
		api.Respond(w, r, api.Error("last name is required field"))
		return
	}

	user, err := cms.UserStore.EntityCreate("user")

	if err != nil {
		api.Respond(w, r, api.Error("User failed to be created: "+err.Error()))
		return
	}

	// log.Println(page)

	if user == nil {
		api.Respond(w, r, api.Error("User failed to be created"))
		return
	}

	user.SetString("first_name", firstName)
	user.SetString("last_name", lastName)
	user.SetString("status", "inactive")

	api.Respond(w, r, api.SuccessWithData("User saved successfully", map[string]interface{}{"user_id": user.ID}))
}

func (cms Cms) pageUsersUserUpdateAjax(w http.ResponseWriter, r *http.Request) {
	userID := strings.Trim(utils.Req(r, "user_id", ""), " ")
	firstName := strings.Trim(utils.Req(r, "first_name", ""), " ")
	lastName := strings.Trim(utils.Req(r, "last_name", ""), " ")
	status := strings.Trim(utils.Req(r, "status", ""), " ")

	if userID == "" {
		api.Respond(w, r, api.Error("User ID is required"))
		return
	}

	user, _ := cms.UserStore.EntityFindByID(userID)

	if user == nil {
		api.Respond(w, r, api.Error("User NOT FOUND with ID "+userID))
		return
	}

	if firstName == "" {
		api.Respond(w, r, api.Error("First name is required field"))
		return
	}

	if lastName == "" {
		api.Respond(w, r, api.Error("Last name is required field"))
		return
	}

	if status == "" {
		api.Respond(w, r, api.Error("status is required field"))
		return
	}

	user.SetString("first_name", firstName)
	user.SetString("last_name", lastName)
	isOk, err := user.SetString("status", status)

	if err != nil {
		api.Respond(w, r, api.Error("User failed to be updated: "+err.Error()))
		return
	}

	if !isOk {
		api.Respond(w, r, api.Error("User failed to be updated"))
		return
	}

	api.Respond(w, r, api.SuccessWithData("User saved successfully", map[string]interface{}{"user_id": user.ID}))
}

func (cms Cms) pageUsersUserUpdate(w http.ResponseWriter, r *http.Request) {
	endpoint := r.Context().Value(keyEndpoint).(string)
	// log.Println(endpoint)

	userID := utils.Req(r, "user_id", "")
	if userID == "" {
		api.Respond(w, r, api.Error("User ID is required"))
		return
	}

	user, _ := cms.UserStore.EntityFindByID(userID)

	if user == nil {
		api.Respond(w, r, api.Error("User NOT FOUND with ID "+userID))
		return
	}

	header := cms.cmsHeader(r.Context().Value(keyEndpoint).(string))
	breadcrums := cms.cmsBreadcrumbs(map[string]string{
		endpoint: "Home",
		(endpoint + "?path=" + PathUsersUserManager):                       "Users",
		(endpoint + "?path=" + PathUsersUserUpdate + "&user_id=" + userID): "Edit user",
	})

	container := hb.NewDiv().ID("user-update").Class("container")
	heading := hb.NewHeading1().HTML("Edit User")
	saveButton := hb.NewButton().HTML("Save").Class("btn btn-success float-end").Attr("v-on:click", "userSave")
	heading.Child(saveButton)

	// Status
	formGroupStatus := bs.FormGroup().Children([]*hb.Tag{
		bs.FormLabel("Status"),
		bs.FormSelect().Attr("v-model", "pageModel.status").Children([]*hb.Tag{
			bs.FormSelectOption("active", "Active"),
			bs.FormSelectOption("inactive", "Inactive"),
			bs.FormSelectOption("trash", "Trash"),
		}),
	})

	// First Name
	formGroupFirstName := bs.FormGroup().Children([]*hb.Tag{
		bs.FormLabel("First Name"),
		bs.FormInput().Attr("v-model", "pageModel.firstName"),
		bs.FormText("The first name of the user"),
	})

	// Last Name
	formGroupLastName := bs.FormGroup().Children([]*hb.Tag{
		bs.FormLabel("Last Name"),
		bs.FormInput().Attr("v-model", "pageModel.lastName"),
		bs.FormText("The last name of the user"),
	})

	// Template
	// templateList, err := cms.EntityStore.EntityList("template", 0, 100, "", "id", "asc")
	// if err != nil {
	// 	api.Respond(w, r, api.Error("Entity list failed to be retrieved "+err.Error()))
	// 	return
	// }
	// formGroupTemplateSelect := bs.FormSelect().Attr("v-model", "pageModel.templateId")
	// formGroupTemplateOptionsEmpty := bs.FormSelectOption("", "- none -")
	// formGroupTemplateSelect.Child(formGroupTemplateOptionsEmpty)
	// for _, template := range templateList {
	// 	templateName, _ := template.GetString("name", "n/a")
	// 	formGroupTemplateOptionsTemplate := bs.FormSelectOption(template.ID, templateName)
	// 	formGroupTemplateSelect.Child(formGroupTemplateOptionsTemplate)
	// }

	container.Children([]*hb.Tag{
		hb.NewHTML(header),
		heading,
		hb.NewHTML(breadcrums),
		formGroupStatus,
		formGroupFirstName,
		formGroupLastName,
	})

	h := container.ToHTML()

	firstName, _ := user.GetString("first_name", "")
	lastName, _ := user.GetString("last_name", "")
	status, _ := user.GetString("status", "")

	firstNameJSON, _ := json.Marshal(firstName)
	lastNameJSON, _ := json.Marshal(lastName)

	inlineScript := `
const userUpdateUrl = "` + endpoint + `?path=users/user-update-ajax";
const userId = "` + userID + `";
const firstName = ` + string(firstNameJSON) + `;
const lastName = ` + string(lastNameJSON) + `;
const status = "` + status + `";
const UserUpdate = {
	data() {
		return {
			pageModel:{
				userId: userId,
				firstName: firstName,
				lastName: lastName,
				status: status,
		    }
		}
	},
	created(){
	},
	methods: {
		userSave(){
			var userId = this.pageModel.userId;
			var status = this.pageModel.status;
			var firstName = this.pageModel.firstName;
			var lastName = this.pageModel.lastName;
			
			$.post(userUpdateUrl, {
				user_id:userId,
				first_name: firstName,
				last_name: lastName,
				status: status,
			}).done((response)=>{
				if (response.status !== "success") {
					return Swal.fire({
						icon: 'error',
						title: 'Oops...',
						text: response.message,
					});
				}

				return Swal.fire({
					icon: 'success',
					title: 'User saved',
				});
			}).fail((result)=>{
				console.log(result);
				return Swal.fire({
					icon: 'error',
					title: 'Oops...',
					text: result,
				});
			});
		}
	}
};
Vue.createApp(UserUpdate).mount('#user-update')
	`

	webpage := Webpage("Edit Page", h)
	webpage.AddStyleURLs([]string{
		"//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/codemirror.min.css",
	})
	webpage.AddScriptURLs([]string{
		"//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/codemirror.min.js",
		"//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/mode/xml/xml.min.js",
		"//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/mode/htmlmixed/htmlmixed.min.js",
		"//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/mode/javascript/javascript.js",
		"//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/mode/css/css.js",
		"//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/mode/clike/clike.min.js",
		"//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/mode/php/php.min.js",
		"//cdnjs.cloudflare.com/ajax/libs/codemirror/2.36.0/formatting.min.js",
		"//cdnjs.cloudflare.com/ajax/libs/codemirror/3.22.0/addon/edit/matchbrackets.min.js",
	})
	webpage.AddScript(inlineScript)
	// webpage.AddScript(ve.VisualeditorScripts())
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(webpage.ToHTML()))
}