package main

import (
	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/text/cases"
	//"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/text/language"

	//"os"
	"strconv"
	"strings"
	"text/template"

	database "github.com/replit/database-go"
	//database "github.com/replit/database-go"
)

// var MONGODB_URI=`mongodb+srv://dbcoder523:DXGXkHcdXWWhAoER@db325-cluster-1.gsh2ghg.mongodb.net/?retryWrites=true&w=majority`

var tmpl *template.Template
var c CentersJson

func init() {

	// f, err := ioutil.ReadFile("ROI-2/js/dc.json")

	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(2)
	// }

	// database.Set("all", string(f))

	//value, _ := database.Get("all")

	// fmt.Println("this is all centers", value)
	// _ = json.Unmarshal([]byte(value), &c)
	tmpl = template.Must(template.ParseGlob("templates/*"))

}

type PageData struct {
	Title      string
	IntroText  string
	Center     DonorCenter
	Locations  []string
	Index      int
	DriveDate  string
	Donorpool  []string
	Location   int
	StaffAdded string
	TDY        string
	TDYDrives  []BloodDrive
	Drive      BloodDrive
	Comments   []Comment
	JsonValue  string
	Today      string
	CommentSig string
	Expire     string
}

// var uri=`mongodb+srv://dbcoder523:DXGXkHcdXWWhAoER@db325-cluster-1.gsh2ghg.mongodb.net/?retryWrites=true&w=majority`

//Begin DB Code

// var mongoclient = getClient()
// var collection = mongoclient.Database("DonationCenters").Collection("centers")
func main() {
	db()

	// AS.AddAccount(Administrator)
	// updateUsers(Administrator.Center, Administrator)
	//GenerateDummyAccounts()

	populateCenterUsers()

	// f, err := ioutil.ReadFile("./dc.json")

	//defer mongoclient.Disconnect(context.TODO())

	fs := http.FileServer(http.Dir("js"))
	http.HandleFunc("/admin-login", adminLoginHandler)
	http.HandleFunc("/admin-portal", adminPortalHandler)
	http.HandleFunc("/reset-pwd", resetPwdHandler)
	//http.HandleFunc("/quick-peek", quickViewHandler)
	//	http.HandleFunc("/center-report", centerReportHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/keep-logged-in", keepLoggedInHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/create-user", createAccountHandler)
	http.HandleFunc("/oic-main", oicmain)
	//http.HandleFunc("/get-drives", getDrivesHandler)
	http.HandleFunc("/index", indexHandler)
	http.HandleFunc("/data-entry", dataentry)
	http.HandleFunc("/add-staff", addStaff)
	http.HandleFunc("/create-drive", createdriveHandler)
	http.HandleFunc("/process-default", processDefault)
	http.HandleFunc("/center-main", centerMain)
	http.HandleFunc("/staff-portal", staffPortalHandler)
	http.HandleFunc("/add-to-drive", addStaffHandler)
	http.HandleFunc("/staff-hours", staffHoursHandler)
	http.HandleFunc("/cost-data", costDataHandler)
	http.HandleFunc("/full-view", fullDetailsHandler)
	http.HandleFunc("/deletedrive", deleteDriveHandler)
	http.HandleFunc("/filter", filterHandler)
	http.HandleFunc("/delete-redirect", deleteRedirectHandler)
	http.HandleFunc("/add-note", addNotesHandler)
	http.HandleFunc("/drive-select", driveSelectHandler)
	http.HandleFunc("/drive-select-staff", driveSelectAddStaff)
	http.HandleFunc("/remove-staff", removeStaffHandler)
	http.HandleFunc("/toggle-active", toggleHandler)

	http.Handle("/js/", http.StripPrefix("/js/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the correct Content-Type header for JavaScript files
		w.Header().Set("Content-Type", "application/javascript")

		// Serve the file using the file server
		fs.ServeHTTP(w, r)
	})))

	//http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./js"))))
	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("./templates"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))

	http.Handle("/png/", http.StripPrefix("/png/", http.FileServer(http.Dir("./png"))))
	http.Handle("/svg/", http.StripPrefix("/svg/", http.FileServer(http.Dir("./svg"))))
	http.HandleFunc("/sw.js", sw)
	log.Fatal(http.ListenAndServe(":80", nil))

}

var AS = &AccountStore{}

func createAccountHandler(resp http.ResponseWriter, req *http.Request) {

	if req.Method == "GET" {
		d := struct {
			Title string
		}{
			"Create Account",
		}
		fmt.Println("GET: Create Account")
		tmpl.ExecuteTemplate(resp, "account-creation", d)
		return
	}

	if req.Method == "POST" {
		name := req.FormValue("name")
		email := req.FormValue("email")
		role := req.FormValue("role")
		center, _ := strconv.Atoi(req.FormValue("center"))

		acc := CreateAccount(role, center, "", email, name) // Fix: Adjusted the arguments passed to CreateAccount

		AS.AddAccount(acc)
		http.Redirect(resp, req, "/create-user", 303)
		return
	}

}

// func centerReportHandler(resp http.ResponseWriter, req *http.Request) {
// 	//get the center number from the cookie
// 	cookie, err := req.Cookie("ROI")
// 	if err != nil {
// 		fmt.Println("No cookie found")
// 		http.Redirect(resp, req, "/login", 303)
// 		return
// 	}
// 	//get the value of the cookie
// 	decoded := DecodeHexToString(cookie.Value)
// 	role := strings.Split(decoded, "-")[0]
// 	//center := strings.Split(strings.Split(decoded, "-")[1], "*")[0]

// 	//check if the role is admin
// 	if role != "admin" {
// 		http.Redirect(resp, req, "/login", 303)
// 		return
// 	}

// 	//get the center number from the request
// 	centerNumber, _ := strconv.Atoi(req.FormValue("center"))
// 	// Get the center's report data
// 	donorcenter := c.Centers[centerNumber]

// 	donorcenter.GenerateReport()
// 	update(centerNumber, donorcenter)

// 	//report := donorcenter.Report // Replace with the actual function to get the center's report data

// 	// Execute the center-report-card template with the center's report as the data
// 	err = tmpl.ExecuteTemplate(resp, "center-report-card", donorcenter.Report)
// 	if err != nil {
// 		http.Error(resp, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// }

func adminLoginHandler(resp http.ResponseWriter, req *http.Request) {

	if req.Method == "GET" {

		d := struct {
			Title string
		}{
			"Admin Login",
		}
		fmt.Println("GET: Admin Login")
		tmpl.ExecuteTemplate(resp, "admin-login", d)
		return
	}

	//create cookie for session management

	if req.Method == "POST" {

		acc := &Account{}
		email := req.FormValue("email")
		pass := req.FormValue("pass")
		there := AS.AccountExistsByEmail(email)
		//TODO:Consider changing how I check for the account. Maybe use the name insted of the email

		if there {
			acc = AS.GetAccountByEmail(email)

		}
		if VerifyPassword(acc.HashedPwd, pass) {
			fmt.Println("now redirecting")
			//if the password is correct, create a cookie and redirect to the admin portal

			c := http.Cookie{
				Name:    "ROI",
				Value:   EncodeStringToHex("admin"),
				Expires: time.Now().Add(10 * time.Minute),
			}

			http.SetCookie(resp, &c)
			http.Redirect(resp, req, "/admin-portal", 303)
		}
		//otherwise, redirect to the login page
		http.Redirect(resp, req, "/admin-login", 303)
		return
	}
}

var CenterUsers = make(map[int][]Account)

func populateCenterUsers() {
	for _, v := range c.Centers {
		CenterUsers[v.Num] = AS.GetAccountsByCenter(v.Num)
	}
}

// create a func that takes interface{} and returns a json string
func toJSON(data interface{}) (string, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func adminPortalHandler(resp http.ResponseWriter, req *http.Request) {

	populateCenterUsers()
	jsonData, _ := toJSON(AS.accounts)
	// cents, err := toJSON(c.Centers)

	// if err != nil {
	// 	fmt.Println("error converting to json")
	// }

	d := struct {
		Title string
		Users map[int][]Account `bson:"users"`
		Names []string          `json:"names" bson:"names"`
		Json  string            `json:"json" bson:"json"`
		//Reports string            `json:"centers" bson:"centers"`
	}{
		"Admin Portal",
		CenterUsers,
		c.getLocations(),
		jsonData,
		//reports,
	}

	if req.Method == "GET" {
		//check for the cookie
		cookie, err := req.Cookie("ROI")
		if err != nil {
			fmt.Println("No cookie found")
			http.Redirect(resp, req, "/admin-login", 303)
			return
		}
		//check the value of the cookie
		if cookie.Value != EncodeStringToHex("admin") {
			fmt.Println("Cookie value is not admin")
			http.Redirect(resp, req, "/admin-login", 303)
			return
		}

		fmt.Println("GET: Admin Portal")

		tmpl.ExecuteTemplate(resp, "admin-portal", d)
		return
	}

}

func loginHandler(resp http.ResponseWriter, req *http.Request) {

	if req.Method == "GET" {
		d := struct {
			Title string
		}{
			"Login",
		}
		fmt.Println("GET:Staff/OIC Login")
		tmpl.ExecuteTemplate(resp, "portal-login", d)
		return
	}

	if req.Method == "POST" {

		uname := req.FormValue("username")
		pass := req.FormValue("password")
		acc := AS.GetAccountByName(uname)
		if acc == nil {
			http.Redirect(resp, req, "/login", 303)
			return
		}
		if acc != nil {
			if VerifyPassword(acc.HashedPwd, pass) {
				//if the password is correct, create a cookie and redirect to the admin portal
				//build string for the cookie value

				//REMEMBER: changr the expiration time on both cv and c
				cv := acc.Role + "-" + strconv.Itoa(acc.Center) + "*" + time.Now().Add(15*time.Minute).Format(time.RFC3339) //this format is used to convert the time to a string. it is the same format used in the time package
				c := &http.Cookie{
					Name:  "ROI",
					Value: EncodeStringToHex(cv),

					//Expires: time.Now().Add(15 * time.Minute),

					//set expires to 1.5 minutes
					Expires: time.Now().Add(15 * time.Minute),
				}

				http.SetCookie(resp, c)
				//check the header to see if the cookie is being set
				fmt.Println("Cookie being set: ", c)
				//check the header
				fmt.Println("Header being set: ", resp.Header())

				//redirect based on the role
				fmt.Println("THIS IS THE Role OF THE PERSON LOGGING IN: ", acc.Role, " AND THE CENTER: ", acc.Center, " AND THE COOKIE VALUE: ", &c.Expires)
				if acc.Role == OIC {
					//build the string
					fmt.Println("Checking the timestamp", c.Expires.UTC())
					newstring := "oic-main?locations=" + strconv.Itoa(acc.Center) //+ "&exp=" + c.Expires.Format(time.RFC3339)
					http.Redirect(resp, req, newstring, 303)
					return
				}
				if acc.Role == DStaff {
					//build the string
					newstring := "staff-portal?locations=" + strconv.Itoa(acc.Center)
					http.Redirect(resp, req, newstring, 303)
					return
				}
			} else {
				http.Redirect(resp, req, "/login", 303)
				return
			}
		}
	}
}

// oicmain is the function that handles the main page for the OIC
func oicmain(resp http.ResponseWriter, req *http.Request) {

	if req.Method == "GET" {
		//check for the cookie
		cookie, err := req.Cookie("ROI")
		fmt.Println("you will be logged out at : ", req.FormValue("exp"))
		if err != nil {
			fmt.Println("No cookie found")
			http.Redirect(resp, req, "/login", 303)

			return
		}
		//check the value of the cookie
		decoded := DecodeHexToString(cookie.Value)

		fmt.Println("the decoded cookie is: ", decoded)

		role := strings.Split(decoded, "-")[0]
		center := strings.Split(strings.Split(decoded, "-")[1], "*")[0]
		expirationTime := strings.Split(decoded, "*")[1]

		fmt.Println("the decoded cookie is: ", decoded)
		if role == DStaff {
			//redirect back to where they came from
			newstring := "staff-portal?locations=" + center
			http.Redirect(resp, req, newstring, 303)
			return
		}

		loc := req.FormValue("locations")

		if loc != center {
			http.Redirect(resp, req, "oic-main?locations="+center, 303)
			return
		}

		// if strings.Contains(decoded, "admin") {
		// 	fmt.Println("Admin is trying to access the OIC page")
		// 	//decoded += "-" + req.FormValue("locations")
		// 	center = req.FormValue("locations")
		// }

		newloc, _ := strconv.Atoi(loc)
		donor_center := c.Centers[newloc]

		fmt.Println("coockie set to expire at: ", expirationTime)

		xb, _ := json.Marshal(c.Centers[newloc])
		s, e := toJSON(expirationTime)
		if e != nil {
			fmt.Println("error converting to json")
		}
		data := PageData{
			IntroText: "Welcome to the Blood Drive Management Portal.",
			Title:     "OIC Dashboard",
			Locations: c.getLocations(),
			Center:    donor_center,
			TDY:       c.Centers[newloc].HASTDY(),
			JsonValue: string(xb),
			Expire:    s,
		}

		//check if decoded has the word admin in it

		tmpl.ExecuteTemplate(resp, "oic-main.gohtml", data)
		return
	}
}

func sw(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("sw hit")
	resp.Header().Set("Service-Worker-Allowed", "/")
	http.ServeFile(resp, req, "./sw.js")
}

func indexHandler(resp http.ResponseWriter, req *http.Request) {

	loc, _ := strconv.ParseInt(req.FormValue("location"), 0, 32)
	center := c.Centers[int(loc)]

	// 	//check for the cookie
	cookie, err := req.Cookie("ROI")
	if err != nil {
		fmt.Println("No cookie found")
		http.Redirect(resp, req, "/login", 303)
		return
	}

	data := PageData{
		IntroText: "Drive Expenses",
		Title:     "Average Costs",
		Center:    center,
		Expire:    cookie.Expires.String(),
	}

	tmpl.ExecuteTemplate(resp, "index.gohtml", data)
	return

}

// dataentry is the function that handles the data entry for the drive
func dataentry(resp http.ResponseWriter, req *http.Request) {

	type Data struct {
		Locations []string
		Donorpool []string
		Center    DonorCenter
		Index     int
		DriveDate string
	}

	if req.Method == "GET" {

		//check for the cookie
		cookie, err := req.Cookie("ROI")
		if err != nil {
			fmt.Println("No cookie found")
			http.Redirect(resp, req, "/login", 303)
			return
		}
		//check the value of the cookie
		decoded := DecodeHexToString(cookie.Value)
		role := strings.Split(decoded, "-")[0]
		donor_center := strings.Split(strings.Split(decoded, "-")[1], "*")[0]

		fmt.Print("This is the donor center", donor_center)

		expirationTime := strings.Split(decoded, "*")[1]

		l, _ := strconv.Atoi(req.FormValue("locations"))

		if role == DStaff {
			//check if the location is the same as the center
			if donor_center != req.FormValue("locations") {
				http.Redirect(resp, req, "staff-portal?locations="+donor_center, 303)
				return
			}

		}
		dd := req.FormValue("Drive") //The value of Drive is the drive date
		center := c.Centers[l]
		littleDrive, _ := center.getDrive(dd)
		xb, _ := json.Marshal(center)
		d := PageData{
			IntroText: "Blood Drive Data Entry",
			Title:     "Data Entry",
			Locations: c.getLocations(),
			Donorpool: []string{"Trainee", "Non-Trainee"},
			Center:    center,
			Index:     l,
			DriveDate: dd,
			Drive:     littleDrive,
			JsonValue: string(xb),
			Expire:    expirationTime,
		}

		tmpl.ExecuteTemplate(resp, "drivedata.html", d)
		return
	}

	if req.Method == "POST" {
		// get the date from params so we can get the drive

		l, _ := strconv.Atoi(req.FormValue("index"))
		center := c.Centers[l]

		drivedate := req.FormValue("drive-date")
		where := req.FormValue("where")
		var registered, deferred, ffpcreated int

		donorpool := req.FormValue("donor-pool")
		drivehours, _ := strconv.ParseFloat(req.FormValue("drive-hours"), 32)
		drive, ok := center.BloodDrives[drivedate]
		if !ok {
			fmt.Println("sorry no drive")
			http.Redirect(resp, req, "data-entry?locations="+strconv.Itoa(l), 303)
			return
		}
		setup, _ := strconv.ParseFloat(req.FormValue("trvl-setup-teardown"), 32)
		registered, _ = strconv.Atoi(req.FormValue("registered-donors"))
		deferred, _ = strconv.Atoi(req.FormValue("deferred-donors"))
		ffpcreated, _ = strconv.Atoi(req.FormValue("ffp-created"))
		signature := req.FormValue("signature")
		//driveType := DriveType{}
		drive.DriveType.Location = where
		drive.DriveType.DonorPool = donorpool
		drive.DriveType.DeferredDonors = deferred
		drive.DriveType.NumOfRegisteredDonors = registered

		o_pos, _ := strconv.Atoi(req.FormValue("o-pos"))
		o_neg, _ := strconv.Atoi(req.FormValue("o-neg"))
		a_pos, _ := strconv.Atoi(req.FormValue("a-pos"))

		a_neg, _ := strconv.Atoi(req.FormValue("a-neg"))
		b_pos, _ := strconv.Atoi(req.FormValue("b-pos"))
		b_neg, _ := strconv.Atoi(req.FormValue("b-neg"))
		ab_pos, _ := strconv.Atoi(req.FormValue("ab-pos"))
		ab_neg, _ := strconv.Atoi(req.FormValue("ab-neg"))

		drive.DriveType.WholeBloodCount.OPositive = o_pos
		drive.DriveType.WholeBloodCount.ONegative = o_neg
		drive.DriveType.WholeBloodCount.APositive = a_pos
		drive.DriveType.WholeBloodCount.ANegative = a_neg
		drive.DriveType.WholeBloodCount.BPositive = b_pos
		drive.DriveType.WholeBloodCount.BNegative = b_neg
		drive.DriveType.WholeBloodCount.ABPositive = ab_pos
		drive.DriveType.WholeBloodCount.ABNegative = ab_neg

		drive.DriveType.WholeBloodCount.Total = drive.DriveType.WholeBloodCount.GetTotal()
		drive.DriveType.NumOfFFPPF24 = ffpcreated
		drive.DriveType.Signature = signature
		drive.DriveType.TravelSetupTeardown = float32(setup)
		drive.DriveType.DriveLength = float32(drivehours)
		//drive.Rep = drive.GenerateReport()
		center.BloodDrives[drivedate] = drive

		c.Centers[l] = center
		center.GenerateReport()
		centerDB, _ := json.Marshal(c)
		database.Set("all", string(centerDB))

		update(center.Num, center)

		newstring := "staff-portal?locations=" + strconv.Itoa(l)
		http.Redirect(resp, req, newstring, 303)

		return
	}

}

func processDefault(res http.ResponseWriter, req *http.Request) {
	//get the drive, then set the values,then save the state, then return

	//get the location first
	location, _ := strconv.Atoi(req.FormValue("locations"))
	//get date next
	//date:= req.FormValue("date")

	kind := req.FormValue("drive-type") // in-house or tdy
	collectionSupplies, _ := strconv.ParseFloat(req.FormValue("collection-sup"), 32)
	processingSupplies, _ := strconv.ParseFloat(req.FormValue("processing-sup"), 32)
	testing, _ := strconv.ParseFloat(req.FormValue("testing"), 32)
	averageCivilianRBC, _ := strconv.ParseFloat(req.FormValue("averageCivilianRbc"), 32)
	averageCivilianFFPpF24, _ := strconv.ParseFloat(req.FormValue("averageCivilianFFPpF24"), 32)
	shipping, _ := strconv.ParseFloat(req.FormValue("shipping"), 32)
	flightCosts, _ := strconv.ParseFloat(req.FormValue("flights"), 32)
	airportParking, _ := strconv.ParseFloat(req.FormValue("airport-parking"), 32)
	perDiem, _ := strconv.ParseFloat(req.FormValue("m-and-e-perdiem"), 32)
	hotelParking, _ := strconv.ParseFloat(req.FormValue("hotel-parking"), 32)
	lodging, _ := strconv.ParseFloat(req.FormValue("lodgingg"), 32)

	center := c.Centers[location]

	//create the default template and add to donor center
	sig := req.FormValue("signature") //string
	Lock.Lock()
	t := createTemplate(kind, float32(collectionSupplies), float32(processingSupplies), float32(testing), float32(averageCivilianRBC), float32(averageCivilianFFPpF24), float32(shipping), float32(flightCosts), float32(airportParking), float32(perDiem), float32(hotelParking), float32(lodging))
	t.Sig = sig
	center.Default = t

	c.Centers[location] = center
	//center.GenerateReport()
	Lock.Unlock()

	centerDB, _ := json.Marshal(c)
	database.Set("all", string(centerDB))
	fmt.Println("Generating reoprt before update....", center.GenerateReport())
	update(center.Num, center)

	newstring := "oic-main?locations=" + strconv.Itoa(location)
	http.Redirect(res, req, newstring, 303)
	return

}

type DataStruct struct {
	Location int
	Title    string
	Center   DonorCenter
}

func addStaff(resp http.ResponseWriter, req *http.Request) {

	if req.Method == "GET" {

		//check for the cookie
		cookie, err := req.Cookie("ROI")
		if err != nil {
			fmt.Println("No cookie found")
			http.Redirect(resp, req, "/login", 303)
			return
		}

		//check the value of the cookie
		decoded := DecodeHexToString(cookie.Value)
		role := strings.Split(decoded, "-")[0]
		donor_center := strings.Split(strings.Split(decoded, "-")[1], "*")[0]
		expirationTime := strings.Split(decoded, "*")[1]
		if role == DStaff {
			//redirect back to where they came from
			newstring := "staff-portal?locations=" + donor_center
			http.Redirect(resp, req, newstring, 303)
			return
		}

		//ensure that nobody can access a center that they are not assigned to

		staffString := req.FormValue("location")
		if role == OIC && donor_center != staffString {
			http.Redirect(resp, req, "oic-main?locations="+donor_center, 303)
			return
		}

		loc, _ := strconv.Atoi(staffString)
		d := PageData{
			IntroText: "Staff Creation Page",
			Title:     "Add Staff",
			Center:    c.Centers[loc],

			Location: loc,
			Expire:   expirationTime,
		}
		tmpl.ExecuteTemplate(resp, "add-staff.html", d)
		return
	}
	if req.Method == "POST" {
		locs := req.FormValue("locations")
		location, _ := strconv.Atoi(locs)
		fn := req.FormValue("first-name")
		ln := req.FormValue("last-name")

		hrly, _ := strconv.ParseFloat(req.FormValue("hourly-wage"), 32)

		caser := cases.Title(language.English)

		//trim last name to just the first letter
		ln = strings.Split(ln, "")[0]
		staff := makeStaff(fmt.Sprint(caser.String(fn)+" "+caser.String(ln)), float32(hrly), location)
		center := &c.Centers[location]
		center.addStaff(staff)
		center.GenerateReport()
		c.Centers[location] = *center
		centerDB, _ := json.Marshal(c)
		database.Set("all", string(centerDB))
		staff.GenerateReport()
		fmt.Println("Generating reoprt before update....", center.GenerateReport())

		update(center.Num, *center)

		http.Redirect(resp, req, "add-staff?location="+strconv.Itoa(center.Num), 303)
		return

	}

}

func centerMain(resp http.ResponseWriter, req *http.Request) {
	//get location from request params
	loc, err := strconv.Atoi(req.FormValue("location"))
	if err != nil {
		fmt.Printf("Sorry, theres been an error %v", err)
	}

	center := c.Centers[loc]

	data := struct {
		Title  string
		Center DonorCenter
		Expire string
	}{
		"Main Page",
		center,
		"",
	}

	tmpl.ExecuteTemplate(resp, "center-main.html", data)
}

func createdriveHandler(resp http.ResponseWriter, req *http.Request) {

	if req.Method == "GET" {

		//check for the cookie
		cookie, err := req.Cookie("ROI")
		if err != nil {
			fmt.Println("No cookie found")
			http.Redirect(resp, req, "/login", 303)
			return
		}

		//check the value of the cookie
		decoded := DecodeHexToString(cookie.Value)
		role := strings.Split(decoded, "-")[0]
		donor_center := strings.Split(strings.Split(decoded, "-")[1], "*")[0]
		expirationTime := strings.Split(decoded, "*")[1]
		if role == DStaff {
			//redirect back to where they came from
			newstring := "staff-portal?locations=" + donor_center
			http.Redirect(resp, req, newstring, 303)
			return
		}

		//ensure that nobody can access a center that they are not assigned to

		if role == OIC && donor_center != req.FormValue("location") {
			http.Redirect(resp, req, "oic-main?locations="+donor_center, 303)
			return
		}

		loc := req.FormValue("location")

		loc1, _ := strconv.Atoi(loc)

		data := PageData{
			IntroText: "BloodDrive Creation",
			Title:     "Create Drive",
			Index:     loc1,
			Center:    c.Centers[loc1],
			Expire:    expirationTime,
		}

		tmpl.ExecuteTemplate(resp, "create-drive.html", data)
		return
	}
	if req.Method == "POST" {
		url := "/oic-main?locations="
		loc, _ := strconv.Atoi(req.FormValue("location"))
		kind := req.FormValue("kind")
		signature := req.FormValue("signature")
		when := req.FormValue("when")
		time := strings.Split(when, "T")[1]

		when = strings.Split(when, "T")[0]
		whens := strings.Split(when, "-")
		when = fmt.Sprint(strings.Join([]string{whens[1], whens[2], whens[0]}, "/"))
		_ = createDrive(kind, signature, when, time, loc)
		center := c.Centers[loc]
		center.GenerateReport()
		loc2 := strconv.Itoa(loc)

		newstring := url + loc2
		http.Redirect(resp, req, newstring, 303)

		return
	}

}

func staffPortalHandler(resp http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("ROI")
	if req.Method == "GET" {

		if err != nil {
			fmt.Println("No cookie found")
			http.Redirect(resp, req, "/login", 303)
			return
		}

		//check the value of the cookie
		decoded := DecodeHexToString(cookie.Value)
		role := strings.Split(decoded, "-")[0]
		donor_center := strings.Split(strings.Split(decoded, "-")[1], "*")[0]
		expirationTime := strings.Split(decoded, "*")[1]
		fmt.Println("Role: ", role, "Center: ", donor_center)
		if role == OIC {
			//redirect back to where they came from
			newstring := "oic-main?locations=" + donor_center
			http.Redirect(resp, req, newstring, 303)
			return
		}

		loc := req.FormValue("locations")

		if loc != donor_center {
			http.Redirect(resp, req, "staff-portal?locations="+donor_center, 303)
			return
		}

		newloc, _ := strconv.Atoi(loc)
		center := c.Centers[newloc]
		c.Centers[newloc] = center
		staffAdded := req.FormValue("staff-added")
		xb, _ := json.Marshal(c.Centers[newloc])

		data := PageData{
			Title:     "Blood Drive Portal",
			Locations: c.getLocations(),
			Center:    center,
			Index:     newloc,
			IntroText: "Welcome to the Blood Drive Data Entry System.",
			JsonValue: string(xb),
			Expire:    expirationTime,
		}
		if staffAdded == "yes" {
			data.StaffAdded = "yes"
		}
		//	fmt.Println("INSIDE staf portal______>>>>>>>.", data)

		tmpl.ExecuteTemplate(resp, "staff-main.gohtml", data)
		return

	}

}

func addStaffHandler(resp http.ResponseWriter, req *http.Request) {

	if req.Method == "POST" {
		err := req.ParseForm()
		if err != nil {
			fmt.Println("someting went wrong when parsing the form: ", err)
			fmt.Println("you will be redirected to the staff portal")
			//http.Redirect(resp, req, "staff-portal", 303)

		}

		//check for the cookie
		cookie, err := req.Cookie("ROI")
		if err != nil {
			fmt.Println("No cookie found")
			http.Redirect(resp, req, "/login", 303)
			return
		}

		//check the value of the cookie
		decoded := DecodeHexToString(cookie.Value)
		role := strings.Split(decoded, "-")[0]
		number := strings.Split(strings.Split(decoded, "-")[1], "*")[0]

		if role == OIC {
			//redirect back to where they came from
			newstring := "oic-main?locations=" + number
			http.Redirect(resp, req, newstring, 303)
			return
		}

		if role == DStaff && number != req.FormValue("center-num") {
			fmt.Println("You are not allowed to access this page", number, req.FormValue("center-num"))
			http.Redirect(resp, req, "staff-portal?locations="+number, 303)
			return
		}

		values := req.Form
		var date string
		var centernum int
		var names []string
		var addremove string

		for k, v := range values {
			if k == "date" {
				date = v[0]
			}
			if k == "add" {
				addremove = "add"
				names = v
			} else if k == "remove" {
				addremove = "remove"
				names = v
			}
			if k == "center-num" {
				centernum, _ = strconv.Atoi(v[0])
			}

		}
		center := c.Centers[centernum]
		drive, _ := center.getDrive(date)

		if addremove == "add" {

			for _, v := range names {
				yn := center.StaffList[v].DoIBelong(date)
				if !yn {

					drive.addStaff(center.StaffList[v], centernum, center)
					//drive.Rep = drive.GenerateReport()
					s := center.StaffList[v]
					s.GenerateReport()

				}

			}
			fmt.Println("ADDED: ", names, " to ", drive, date)

			centerDB, _ := json.Marshal(c)
			database.Set("all", string(centerDB))
			fmt.Println("Generating reoprt before update....", center.GenerateReport())

			center.GenerateReport()
			c.Centers[centernum] = center
			update(center.Num, center)
		}

		if addremove == "remove" {
			for _, v := range names {

				yn := center.StaffList[v].DoIBelong(date)

				if yn {
					fmt.Printf("Removing: %v from: %v", center.StaffList[v].Name, drive)

					staff := drive.Staff[v]

					drive.removeStaff(staff, centernum, center)
					//drive.Rep = drive.GenerateReport()
					staff.GenerateReport()

				}

			}

			centerDB, _ := json.Marshal(c)
			database.Set("all", string(centerDB))
			fmt.Println("Generating reoprt before update....", center.GenerateReport())
			center.GenerateReport()
			update(center.Num, center)

			fmt.Println("removed: ", names, " from ", drive)

		}

		center.BloodDrives[date] = drive
		c.Centers[centernum] = center

		centerDB, _ := json.Marshal(c)
		database.Set("all", string(centerDB))
		c.Centers[centernum] = center
		fmt.Println("Generating reoprt before update....", center.GenerateReport())

		center.GenerateReport()
		update(centernum, center)

		newstring := "add-to-drive?date=" + date + "&location=" + strconv.Itoa(centernum)
		http.Redirect(resp, req, newstring, 303)
		return
	}

	if req.Method == "GET" {
		fmt.Println("GET: Add Staff")
		//check for the cookie
		cookie, err := req.Cookie("ROI")
		if err != nil {
			fmt.Println("No cookie found")
			http.Redirect(resp, req, "/login", 303)
			return
		}

		//check the value of the cookie
		decoded := DecodeHexToString(cookie.Value)
		role := strings.Split(decoded, "-")[0]
		donor_center := strings.Split(strings.Split(decoded, "-")[1], "*")[0]
		expirationTime := strings.Split(decoded, "*")[1]
		if role == OIC {
			//redirect back to where they came from
			newstring := "oic-main?locations=" + donor_center
			http.Redirect(resp, req, newstring, 303)
			return
		}

		if role == DStaff && donor_center != req.FormValue("location") {
			fmt.Println("GET: You are not allowed to access this page", donor_center, req.FormValue("location"))
			http.Redirect(resp, req, "staff-portal?locations="+donor_center, 303)
			return
		}

		location, _ := strconv.Atoi(req.FormValue("location"))
		//get the center
		center := c.Centers[location]
		xb, _ := json.Marshal(center)
		drive, tf := center.getDrive(req.FormValue("date"))
		if !tf {
			newstring := "add-to-drive?date=" + req.FormValue("date") + "&location=" + strconv.Itoa(center.Num)
			http.Redirect(resp, req, newstring, 303)
		}

		data := PageData{
			IntroText:  "Drive Assigmnent",
			Center:     center,
			Title:      "Drive Assignment",
			JsonValue:  string(xb),
			Drive:      drive,
			CommentSig: "~ Staff",
			Expire:     expirationTime,
		}

		tmpl.ExecuteTemplate(resp, "drive-assignment.html", data)
		return
	}

}

func staffHoursHandler(resp http.ResponseWriter, req *http.Request) {

	if req.Method == "GET" {
		//check for the cookie
		cookie, err := req.Cookie("ROI")
		if err != nil {
			fmt.Println("No cookie found")
			http.Redirect(resp, req, "/login", 303)
			return
		}

		//check the value of the cookie
		decoded := DecodeHexToString(cookie.Value)
		role := strings.Split(decoded, "-")[0]
		donor_center := strings.Split(strings.Split(decoded, "-")[1], "*")[0]
		expirationTime := strings.Split(decoded, "*")[1]
		if role == OIC {
			//redirect back to where they came from
			newstring := "oic-main?locations=" + donor_center
			http.Redirect(resp, req, newstring, 303)
			return
		}

		if role == DStaff && donor_center != req.FormValue("center-num") {
			http.Redirect(resp, req, "staff-portal?locations="+donor_center, 303)
			return
		}

		date := req.FormValue("date")
		//get the center
		centernum, _ := strconv.Atoi(req.FormValue("center-num"))
		center := c.Centers[centernum]
		bd := center.BloodDrives[date]
		pd := PageData{
			Title:  "Staff tasks",
			Drive:  bd,
			Center: center,
			Expire: expirationTime,
		}
		tmpl.ExecuteTemplate(resp, "staff-hours.html", pd)
		return
	}

	if req.Method == "POST" {

		date := req.FormValue("date")
		centernum, _ := strconv.Atoi(req.FormValue("center-num"))
		name := req.FormValue("staff")

		processingHrs, _ := strconv.ParseFloat(req.FormValue("processing-hrs"), 32)
		collectionHrs, _ := strconv.ParseFloat(req.FormValue("collection-hrs"), 32)
		supportingHrs, _ := strconv.ParseFloat(req.FormValue("supporting-hrs"), 32)

		//get the center
		center := c.Centers[centernum]
		//get the staff member
		staff := center.StaffList[name]

		//get the drive hours struct of the date
		driveHours := DriveHoursWorked{}

		//set the values
		driveHours.HrsCollecting = float32(collectionHrs)
		driveHours.HrsSupporting = float32(supportingHrs)
		driveHours.HrsProcessing = float32(processingHrs)

		//save the changes

		Lock.Lock()

		staff.DriveHours[date] = driveHours
		//check is staff has reports
		staff.GenerateReport()
		center.StaffList[name] = staff
		bd := center.BloodDrives[date]
		bd.Staff[staff.Name] = staff

		Lock.Unlock()
		fmt.Println("Generating reoprt before update....", center.GenerateReport())

		center.GenerateReport()
		c.Centers[centernum] = center

		centerDB, _ := json.Marshal(c)
		database.Set("all", string(centerDB))

		update(centernum, c.Centers[centernum])

		//TODO:ADD LOGIC FOR UPDATING DB
		fmt.Println(staff)

		newstring := "add-to-drive?date=" + bd.GetStartTime() + "&location=" + strconv.Itoa(centernum)
		http.Redirect(resp, req, newstring, 303)

		return
	}

}

// costDataHandler is the function that handles the cost data entry for TDY drives
func costDataHandler(resp http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		//check for the cookie
		cookie, err := req.Cookie("ROI")
		if err != nil {
			fmt.Println("No cookie found")
			http.Redirect(resp, req, "/login", 303)
			return
		}

		//check the value of the cookie
		decoded := DecodeHexToString(cookie.Value)
		role := strings.Split(decoded, "-")[0]
		donor_center := strings.Split(strings.Split(decoded, "-")[1], "*")[0]
		expirationTime := strings.Split(decoded, "*")[1]
		if role == DStaff {
			//redirect back to where they came from
			newstring := "staff-portal?locations=" + donor_center
			http.Redirect(resp, req, newstring, 303)
			return
		}

		if role == OIC && donor_center != req.FormValue("location") {
			http.Redirect(resp, req, "oic-main?locations="+donor_center, 303)
			return
		}

		type Data struct {
			Title     string
			Center    DonorCenter
			Index     int
			DriveDate string
			Expire    string
		}

		l, _ := strconv.Atoi(req.FormValue("location"))
		drivedate := req.FormValue("Drive") //The value of Drive is the drive date
		center := c.Centers[l]

		if center.GetTDY() == nil {
			newstring := "oic-main?locations=" + strconv.Itoa(l)
			http.Redirect(resp, req, newstring, 303)
		}
		drive, _ := center.getDrive(drivedate)
		js, _ := toJSON(drive.Cost)
		d := PageData{
			Title:     "Drive Cost",
			IntroText: "TDY Drive Cost",
			Center:    center,
			DriveDate: drivedate,
			TDYDrives: center.GetTDY(),
			Drive:     drive,
			JsonValue: js,
			Expire:    expirationTime,
		}
		tmpl.ExecuteTemplate(resp, "cost-data.html", d)
		return
	}

	if req.Method == "POST" {

		l, _ := strconv.Atoi(req.FormValue("location"))
		drivedate := req.FormValue("date") //The value of Drive is the drive date

		center := c.Centers[l]

		if center.BloodDrives == nil {
			fmt.Println("No drives found for this center")
			http.Redirect(resp, req, "/oic-main?locations={{string(l)}}", 303)
			return
		}

		if drivedate == "" {
			fmt.Println("No drives found for this center")
			http.Redirect(resp, req, "/oic-main?locations={{string(l)}}", 303)
			return
		}
		Lock.Lock()
		drive := center.BloodDrives[drivedate]
		drive.Cost.EnteredBy = req.FormValue("signature")

		//TDY values
		shipping, _ := strconv.ParseFloat(req.FormValue("shipping"), 32)
		drive.Cost.TDY.Shipping = float32(shipping)

		flightCosts, _ := strconv.ParseFloat(req.FormValue("flights"), 32)
		drive.Cost.TDY.FlightCosts = float32(flightCosts)

		airportParking, _ := strconv.ParseFloat(req.FormValue("airport-parking"), 32)
		drive.Cost.TDY.AirportParking = float32(airportParking)

		perDiem, _ := strconv.ParseFloat(req.FormValue("m-and-e-perdiem"), 32)
		drive.Cost.TDY.MEPerDiem = float32(perDiem)

		hotelParking, _ := strconv.ParseFloat(req.FormValue("hotel-parking"), 32)
		drive.Cost.TDY.HotelParking = float32(hotelParking)

		lodging, _ := strconv.ParseFloat(req.FormValue("lodging"), 32)
		drive.Cost.TDY.Lodging = float32(lodging)

		drive.Cost.BeenChanged = true

		center.BloodDrives[drivedate] = drive
		Lock.Unlock()
		fmt.Println("Generating reoprt before update....", center.GenerateReport())

		center.GenerateReport()
		c.Centers[l] = center

		centerDB, _ := json.Marshal(c)
		database.Set("all", string(centerDB))

		update(center.Num, center)

		newstring := "oic-main?locations=" + strconv.Itoa(l)
		http.Redirect(resp, req, newstring, 303)
		return

	}

}

func fullDetailsHandler(resp http.ResponseWriter, req *http.Request) {
	//get the cookie

	cookie, err := req.Cookie("ROI")
	if err != nil {
		fmt.Println("No cookie found")
		http.Redirect(resp, req, "/login", 303)
		return
	}

	//check the value of the cookie
	decoded := DecodeHexToString(cookie.Value)
	role := strings.Split(decoded, "-")[0]
	donor_center := strings.Split(strings.Split(decoded, "-")[1], "*")[0]
	expirationTime := strings.Split(decoded, "*")[1]
	if role == DStaff {
		//redirect back to where they came from
		newstring := "staff-portal?locations=" + donor_center
		http.Redirect(resp, req, newstring, 303)
		return
	}

	if role == OIC && donor_center != req.FormValue("center-num") {
		http.Redirect(resp, req, "oic-main?locations="+donor_center, 303)
		return
	}

	location, _ := strconv.Atoi(req.FormValue("center-num"))
	center := c.Centers[location]

	date := req.FormValue("date")
	time := req.FormValue("time")

	drive := center.BloodDrives[fmt.Sprintf("%v %v", date, time)]
	pd := PageData{
		Title:      "Full Details",
		Center:     center,
		Drive:      drive,
		IntroText:  "Collection Data",
		CommentSig: "~ OIC",
		Expire:     expirationTime,
	}
	tmpl.ExecuteTemplate(resp, "full-details.html", pd)
	return
}

func deleteDriveHandler(resp http.ResponseWriter, req *http.Request) {
	//check for the cookie

	cookie, err := req.Cookie("ROI")
	if err != nil {
		fmt.Println("No cookie found")
		http.Redirect(resp, req, "/login", 303)
		return
	}

	//check the value of the cookie
	decoded := DecodeHexToString(cookie.Value)
	role := strings.Split(decoded, "-")[0]
	donor_center := strings.Split(strings.Split(decoded, "-")[1], "*")[0]

	if role == DStaff {
		//redirect back to where they came from
		newstring := "staff-portal?locations=" + donor_center
		http.Redirect(resp, req, newstring, 303)
		return
	}

	number, _ := strconv.Atoi(req.FormValue("number"))

	if role == OIC && donor_center != req.FormValue("number") {
		http.Redirect(resp, req, "oic-main?locations="+donor_center, 303)
		return
	}

	drivedate := req.FormValue("drivedate")
	fmt.Println(drivedate, "DRIVEDATE(inside delete)")

	//check to see if the drive exists
	center := c.Centers[number]
	_, ok := center.BloodDrives[drivedate]
	if !ok {
		fmt.Println("No drive found")
		http.Redirect(resp, req, "oic-main?locations="+strconv.Itoa(number), 303)
		return
	}

	//check to see if the drive has any staff assigned to it, if so, dont
	//delete the drive
	if len(center.BloodDrives[drivedate].Staff) > 0 {
		// fmt.Println("Drive has staff assigned to it", center.BloodDrives[drivedate].Staff)
		// http.Redirect(resp, req, "oic-main?locations="+strconv.Itoa(number), 303)
		// return
		fmt.Println("Drive has staff assigned to it", center.BloodDrives[drivedate].Staff)
	}

	deleteDrive(drivedate, number)

	http.Redirect(resp, req, "oic-main?locations="+strconv.Itoa(number), 303)

}

func filterHandler(res http.ResponseWriter, req *http.Request) {
	gyr := req.FormValue("color")
	newmap := map[string]BloodDrive{}
	centernumber, _ := strconv.Atoi(req.FormValue("center-number"))
	center := c.Centers[centernumber]
	drives := center.FilterByScore(gyr)
	for _, v := range drives {
		newmap[v.DriveType.Date+" "+v.DriveType.StartTime] = v
	}
	center.BloodDrives = newmap

	data := PageData{
		IntroText: "Welcome to the Blood Drive Management Portal.",
		Title:     "OIC Dashboard",
		Locations: c.getLocations(),
		Center:    center,
		TDY:       c.Centers[centernumber].HASTDY(),
		//Expire:    time.Now().Add(15 * time.Minute).String(),
	}
	fmt.Println("this is the drives: ", drives)
	tmpl.ExecuteTemplate(res, "oic-main.gohtml", data)

}

func deleteRedirectHandler(res http.ResponseWriter, req *http.Request) {
	number, _ := strconv.Atoi(req.FormValue("number"))

	drivedate := req.FormValue("drivedate")

	//check for the cookie
	cookie, err := req.Cookie("ROI")
	if err != nil {
		fmt.Println("No cookie found")
		http.Redirect(res, req, "/login", 303)
		return
	}
	decoded := DecodeHexToString(cookie.Value)
	expirationTime := strings.Split(decoded, "*")[1]
	data := PageData{
		DriveDate: drivedate,
		Index:     number,
		Title:     "Delete Drive",
		Expire:    expirationTime,
	}
	tmpl.ExecuteTemplate(res, "delete-drive", data)
	return
}

func addNotesHandler(resp http.ResponseWriter, req *http.Request) {

	//check for the cookie
	cookie, err := req.Cookie("ROI")
	if err != nil {
		fmt.Println("No cookie found")
		http.Redirect(resp, req, "/login", 303)
		return
	}

	centernumber, _ := strconv.Atoi(req.FormValue("center-number"))
	drivedate := req.FormValue("drive-date")
	notes := req.FormValue("comments")
	signature := req.FormValue("signature")
	center := c.Centers[centernumber]
	drive, _ := center.getDrive(drivedate)
	comments := drive.DriveType.Comments

	decoded := DecodeHexToString(cookie.Value)
	exp := strings.Split(decoded, "*")[1]

	data := PageData{
		Title:     "Add Comment",
		Center:    center,
		DriveDate: drivedate,
		Comments:  comments,
		Expire:    exp,
	}
	if req.Method == "GET" {
		tmpl.ExecuteTemplate(resp, "notes.html", data)
		return
	}

	if req.Method == "POST" {
		drive, _ := center.getDrive(drivedate)
		drive.addComment(notes, signature, centernumber)
		//TODO:switch redirect based on comment signature

		if signature == "~ OIC" {

			fmt.Println(drivedate, "DRIVEDATE(inside add notes)")
			var date, time string
			date = strings.Split(drivedate, " ")[0]
			time = strings.Split(drivedate, " ")[1]

			http.Redirect(resp, req, "full-view?center-num="+strconv.Itoa(centernumber)+"&date="+date+"&time="+time, 303)

		}

		if signature == "~ Staff" {

			http.Redirect(resp, req, "add-to-drive?location="+strconv.Itoa(centernumber)+"&date="+drivedate, 303)

		}
	}
}

func driveSelectHandler(resp http.ResponseWriter, req *http.Request) {

	if req.Method == "GET" {
		if req.FormValue("locations") == "" {

		}

		//check for the cookie
		cookie, err := req.Cookie("ROI")
		if err != nil {
			fmt.Println("No cookie found")
			http.Redirect(resp, req, "/login", 303)
			return
		}

		decoded := DecodeHexToString(cookie.Value)
		expirationTime := strings.Split(decoded, "*")[1]

		l, _ := strconv.Atoi(req.FormValue("locations"))
		dd := req.FormValue("Drive") //The value of Drive is the drive date
		center := c.Centers[l]
		littleDrive, _ := center.getDrive(dd)
		xb, _ := json.Marshal(center)
		d := PageData{
			IntroText: "Collection Data",
			Title:     "Full Details",
			Locations: c.getLocations(),
			Donorpool: []string{"Trainee", "Non-Trainee"},
			Center:    center,
			Index:     l,
			DriveDate: dd,
			Drive:     littleDrive,
			JsonValue: string(xb),
			Expire:    expirationTime,
		}

		tmpl.ExecuteTemplate(resp, "full-details.html", d)
		return
	}

}

func driveSelectAddStaff(resp http.ResponseWriter, req *http.Request) {

	if req.Method == "GET" {
		if req.FormValue("locations") == "" {

		}

		//check for the cookie
		cookie, err := req.Cookie("ROI")
		if err != nil {
			fmt.Println("No cookie found")
			http.Redirect(resp, req, "/login", 303)
			return
		}

		l, _ := strconv.Atoi(req.FormValue("locations"))
		dd := req.FormValue("Drive") //The value of Drive is the drive date
		center := c.Centers[l]
		littleDrive, _ := center.getDrive(dd)
		xb, _ := json.Marshal(center)

		decoded := DecodeHexToString(cookie.Value)
		expirationTime := strings.Split(decoded, "*")[1]

		d := PageData{
			// may need to change this to drive-assignment
			IntroText: "Drive Assigmnent",
			Title:     "Full Details",
			Locations: c.getLocations(),
			Donorpool: []string{"Trainee", "Non-Trainee"},
			Center:    center,
			Index:     l,
			DriveDate: dd,
			Drive:     littleDrive,
			JsonValue: string(xb),
			Expire:    expirationTime,
		}

		tmpl.ExecuteTemplate(resp, "drive-assignment.html", d)
		return
	}

}

func removeStaffHandler(resp http.ResponseWriter, req *http.Request) {
	name := req.FormValue("name")
	num, _ := strconv.Atoi(req.FormValue("center"))

	center := c.Centers[num]
	var staff Staff

	if _, ok := center.StaffList[name]; ok {
		staff = center.StaffList[name]
	}

	staff.removeAllDrives()
	center.removeStaff(name)
	center.GenerateReport()
	centerDB, _ := json.Marshal(c)
	database.Set("all", string(centerDB))

	c.Centers[num] = center

	http.Redirect(resp, req, "add-staff?location="+strconv.Itoa(num), 303)
}

func toggleHandler(resp http.ResponseWriter, req *http.Request) {

	//check for the cookie
	cookie, err := req.Cookie("ROI")
	if err != nil {
		fmt.Println("No cookie found")
		http.Redirect(resp, req, "/login", 303)
		return
	}

	//check the value of the cookie
	decoded := DecodeHexToString(cookie.Value)
	role := strings.Split(decoded, "-")[0]
	number := strings.Split(decoded, "-")[1]

	if role == DStaff {
		//redirect back to where they came from
		newstring := "staff-portal?locations=" + number
		http.Redirect(resp, req, newstring, 303)
		return
	}
	number=strings.Split(number,"*")[0]
	
	if role == OIC && number != req.FormValue("num") {
		fmt.Println("You are not allowed to access this page", number, "centernumber:",req.FormValue("num"))
		http.Redirect(resp, req, "oic-main?locations="+number, 303)
		return
	}

	name := req.FormValue("name")
	num, _ := strconv.Atoi(req.FormValue("num"))
	fmt.Println("these are the values received: ", name, num)
	center := c.Centers[num]
	staff := center.StaffList[name]

	fmt.Println("BEFORE TOGGLED ACTIVITY", staff.Name, staff.Active)
	Lock.Lock()

	if staff.Active == true {
		staff.Active = false

	} else if staff.Active == false {
		staff.Active = true

	}

	center.StaffList[name] = staff
	
	c.Centers[num] = center

	Lock.Unlock()
	staff.GenerateReport()
	center.GenerateReport()
	fmt.Println("AFTER TOGGLED ACTIVITY", staff.Name, staff.Active)

	update(num, center)

	http.Redirect(resp, req, "add-staff?location="+strconv.Itoa(num), 303)
}

func resetPwdHandler(resp http.ResponseWriter, req *http.Request) {
	//check for the cookie
	cookie, err := req.Cookie("ROI")
	if err != nil {
		fmt.Println("No cookie found")
		http.Redirect(resp, req, "/login", 303)
		return
	}

	//check the value of the cookie
	decoded := DecodeHexToString(cookie.Value)
	role := strings.Split(decoded, "-")[0]
	donor_center := strings.Split(strings.Split(decoded, "-")[1], "*")[0]

	if role == DStaff {
		//redirect back to where they came from
		newstring := "staff-portal?locations=" + donor_center
		http.Redirect(resp, req, newstring, 303)
		return
	}

	if role == OIC {
		http.Redirect(resp, req, "oic-main?locations="+donor_center, 303)
		return
	}

	name := req.FormValue("name")
	acc := AS.GetAccountByName(name)
	if acc == nil {
		fmt.Println("\n\nAccount Not Found: ", name)
		http.Redirect(resp, req, "/admin-portal", 303)
		return
	}
	fmt.Println("\n\nAccount Found Old Password: ")
	rando := GenerateRandomName()
	encpwd, err := EncodePassword(rando)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	acc.pwd = rando
	acc.HashedPwd = encpwd
	fmt.Println("\n\naccount has been reset: ", acc.pwd, acc.HashedPwd, acc.Center, acc.Role)
	AS.UpdateAccount(*acc)
	http.Redirect(resp, req, "/admin-portal", 303)
}

func keepLoggedInHandler(resp http.ResponseWriter, req *http.Request) {
	//get values from the form body
	yesno := req.FormValue("keep")
	fmt.Println("keepLoggedIn: ", yesno)
	//update the cookie
	cookie, err := req.Cookie("ROI")
	if err != nil {
		fmt.Println("No cookie found")
		http.Redirect(resp, req, "/login", 303)
		return
	}

	//get the value of the cookie
	decoded := DecodeHexToString(cookie.Value)
	role := strings.Split(decoded, "-")[0]
	donor_center := strings.Split(strings.Split(decoded, "-")[1], "*")[0]
	c := &http.Cookie{
		Name:    "ROI",
		Value:   cookie.Value,
		Expires: time.Now().Add(10 * time.Minute),
	}
	//add 10 minutes to the cookie
	cookie.MaxAge = 600 //10 minutes
	http.SetCookie(resp, c)
	fmt.Println("Five minutes added to the cookie")
	//redirect to the appropriate page
	if role == DStaff {
		//redirect back to where they came from
		newstring := "staff-portal?locations=" + donor_center
		http.Redirect(resp, req, newstring, 303)
		return
	}

	if role == OIC {
		http.Redirect(resp, req, "oic-main?locations="+donor_center, 303)
		return
	}

}

func logoutHandler(resp http.ResponseWriter, req *http.Request) {
	//get the cookie
	cookie, err := req.Cookie("ROI")
	if err != nil {
		fmt.Println("No cookie found")
		http.Redirect(resp, req, "/login", 303)
		return
	}
	//delete the cookie
	cookie.MaxAge = -1
	http.SetCookie(resp, cookie)

	//redirect to the login page
	http.Redirect(resp, req, "/login", 303)
}

// func getDrivesHandler(resp http.ResponseWriter, req *http.Request) {
// 	//get the cookie
// 	cookie, err := req.Cookie("ROI")
// 	if err != nil {
// 		fmt.Println("No cookie found")
// 		http.Redirect(resp, req, "/login", 303)
// 		return
// 	}
// 	fmt.Println("GET DRIVES...valid cookie", cookie)
// 	center, _ := strconv.Atoi(req.FormValue("center"))

// 	// get the center
// 	c := c.Centers[center]
// 	//var reports []interface{}

// 	//fmt.Println("Drives: ", reports)

// 	//convert to json then send back to the client
// 	// jsonDrives, er := toJSON(drives)
// 	// if er != nil {
// 	// 	fmt.Println("Error: ", er)
// 	// }
// 	js, e := toJSON(reports)
// 	if e == nil {
// 		//send the json back to the client
// 		fmt.Println("Error: ", e)
// 		resp.Header().Set("Content-Type", "application/json")
// 		resp.Write([]byte(js))
// 	}
// 	if e != nil {
// 		fmt.Println("Error: ", e)
// 	}
// 	return
// }

// func quickViewHandler(resp http.ResponseWriter, req *http.Request) {
// 	//get the cookie
// 	cookie, err := req.Cookie("ROI")
// 	if err != nil {
// 		fmt.Println("No cookie found")
// 		http.Redirect(resp, req, "/login", 303)
// 		return
// 	}

// 	//check the value of the cookie
// 	decoded := DecodeHexToString(cookie.Value)
// 	role := strings.Split(decoded, "-")[0]
// 	number := strings.Split(strings.Split(decoded, "-")[1], "*")[0]

// 	// if role == DStaff {
// 	// 	//redirect back to where they came from
// 	// 	newstring := "staff-portal?locations=" + number
// 	// 	http.Redirect(resp, req, newstring, 303)
// 	// 	return
// 	// }

// 	// if role == OIC && number != req.FormValue("center-num") {
// 	// 	http.Redirect(resp, req, "oic-main?locations="+number, 303)
// 	// 	return
// 	// }

// 	//get the center
// 	num, _ := strconv.Atoi(number)
// 	center := c.Centers[num]
// 	var reports []interface{}
// 	//get the drive
// 	for _, v := range center.BloodDrives {
// 		rep := v.GenerateReport()
// 		reports = append(reports, rep)

// 	}
// 	fmt.Println("DRIVES: ", reports)
// 	//convert to json
// 	js, _ := toJSON(reports)
// 	//send the json back to the client
// 	//resp.Header().Set("Content-Type", "application/json")
// 	tmpl.ExecuteTemplate(resp, "quickpeek.html", reports)
// 	//resp.Write([]byte(js))
// 	fmt.Println("sendeing back the drive: ", js, role)
// 	return
// }
