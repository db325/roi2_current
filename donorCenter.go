package main

import (
	"fmt"
	"strings"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
)

var Lock sync.RWMutex

type CenterController struct {
	session *mongo.Collection
}

func NewCenterControler(s *mongo.Collection) *CenterController {
	return &CenterController{s}

}

var DonorCentersDB = map[string]*DonorCenter{}

type DonorCenter struct {
	//sync.RWMutex
	session       *mongo.Collection
	Branch        string                `bson:"branch" json:"branch"`
	Name          string                `bson:"name" json:"name"`
	Hours         string                `bson:"hours" json:"hours"`
	Phone         string                `bson:"phone" json:"phone"`
	DonationTypes []string              `bson:"donationtypes" json:"donationTypes"`
	Addr          Address               `bson:"addr" json:"address"`
	SocialMedia   SocMedia              `bson:"socialmedia" json:"socialMedia"`
	BloodDrives   map[string]BloodDrive `bson:"blooddrives" json:"bloodDrives"`
	Region        string                `bson:"region" json:"region"`
	Recruiters    []Recruiter           `bson:"recruiters" json:"recruiters"`
	Num           int                   `bson:"num" json:"num"`
	Default       ConfigTemplate        `bson:"defaultcosts" json:"defaultCosts"`
	StaffList     map[string]Staff      `bson:"stafflist" json:"staffList"`
	OID           string                `bson:"oid" json:"oid"`
	Report        Report                `bson:"report" json:"report"`
}

func (dc *DonorCenter) AverageStaffPerDrive() float64 {
	var totalStaff int
	for _, v := range dc.BloodDrives {
		totalStaff += len(v.Staff)
	}
	return float64(totalStaff) / float64(len(dc.BloodDrives))
}

func (dc *DonorCenter) TotalCostOnAverage() float64 {
	var totalCost float64
	for _, v := range dc.BloodDrives {

		totalCost = float64(v.GetTotalExpenses())
	}
	return totalCost / float64(len(dc.BloodDrives))
}

func (dc *DonorCenter) AverageWholeBloodDonations() float64 {
	var totalDonations int
	for _, v := range dc.BloodDrives {
		totalDonations += v.DriveType.WholeBloodCount.Total
	}
	return float64(totalDonations) / float64(len(dc.BloodDrives))
}

func (dc *DonorCenter) AverageDriveScore() passOrFail {
	var totalScore float64
	var avgColor string
	var green, yellow, red int
	for _, v := range dc.BloodDrives {
		totalScore += float64(zeroOut(v.GetPercentage().Number))
		//print out the actuall percentage with 1 decimal point using the GetPercentage.Number method
		fmt.Printf("Drive: %s, Score: %.1f%%\n", v.GetStartTime(), v.GetPercentage().Number)
		//return number as a percent string
		pctg := fmt.Sprintf("%.1f%%", v.GetPercentage().Number)

		fmt.Println("Drive: ", v.GetStartTime(), "Score: ", pctg)

		avgColor = v.GetPercentage().Color
		switch avgColor {
		case "green":
			green++
		case "yellow":
			yellow++
		case "red":
			red++
		}
	}
	//get the color with the highest count
	var maxColor string
	if green > yellow && green > red {
		maxColor = "green"
	}
	if yellow > green && yellow > red {
		maxColor = "yellow"
	}
	if red > green && red > yellow {
		maxColor = "red"
	}
	if totalScore/float64(zeroOut(float64(len(dc.BloodDrives)))) == 0.5 {
		maxColor = "orange"
	}

	fmt.Println("the average color is: ", maxColor, "\n\n and here are the counts..", "green: ", green, "yellow: ", yellow, "red: ", red)
	return passOrFail{maxColor, totalScore / float64(zeroOut(float64(len(dc.BloodDrives))))}

}

func (dc *DonorCenter) Tasks() interface{} {
	var collecting, processing, supporting float64
	for _, v := range dc.BloodDrives {
		collecting += float64(v.GetTotalCollectionHours())
		processing += float64(v.GetTotalProcessingHours())
		supporting += float64(v.GetTotalSupportHours())
	}
	//check to see which has the highest value and return it
	if collecting > processing && collecting > supporting {
		//return a struct with the task and the total hours
		return struct {
			Task  string
			Hours float64
		}{"collecting", collecting}

	}
	if processing > collecting && processing > supporting {
		return struct {
			Task  string
			Hours float64
		}{"processing", processing}
	}
	if supporting > collecting && supporting > processing {
		return struct {
			Task  string
			Hours float64
		}{"supporting", supporting}
	}
	return nil
}

func (dc *DonorCenter) GenerateReport() Report {
	// Create a report for the center
	//read lock the map
	Lock.Lock()

	data := struct {
		CenterName                string  `json:"centerName"`
		CenterNum                 int     `json:"centerNum"`
		AverageCostPerDrive       float64 `json:"averageCostPerDrive"`
		AverageDriveScore         float64 `json:"averageDriveScore"`
		AverageDriveColor         string  `json:"averageDriveColor"`
		MostTimeSpentOn           interface{}
		AverageTotalRBC           float64       `json:"averageTotalRBC"`
		TotalNumberofStaffMembers float64       `json:"totalNumberofStaffMembers"`
		TotalNumberofBloodDrives  float64       `json:"totalNumberofBloodDrives"`
		TDY                       float64       `json:"tdy"`
		InHouse                   float64       `json:"inHouse"`
		Mobile                    float64       `json:"mobile"`
		GreenDrives               float64       `json:"greenDrives"`
		YellowDrives              float64       `json:"yellowDrives"`
		RedDrives                 float64       `json:"redDrives"`
		AverageStaffAtEachDrive   float64       `json:"averageStaffAtEachDrive"`
		DriveReports              []interface{} `json:"driveReports"`
	}{
		CenterName:                dc.Name,
		CenterNum:                 dc.Num,
		AverageCostPerDrive:       float64(zeroOut(dc.TotalCostOnAverage())),
		AverageDriveScore:         float64(zeroOut(dc.AverageDriveScore().Number)),
		AverageDriveColor:         dc.AverageDriveScore().Color,
		MostTimeSpentOn:           dc.Tasks(),
		AverageTotalRBC:           float64(zeroOut(dc.AverageWholeBloodDonations())),
		TotalNumberofStaffMembers: float64(zeroOut(float64(len(dc.StaffList)))),
		TotalNumberofBloodDrives:  float64(zeroOut(float64(len(dc.BloodDrives)))),
		TDY:                       float64(zeroOut(float64(len(dc.GetTDY())))),
		InHouse:                   float64(zeroOut(float64(len(dc.GetInHouse())))),
		Mobile:                    float64(zeroOut(float64(len(dc.GetMobile())))),
		GreenDrives:               float64(zeroOut(float64(len(dc.FilterByScore("green"))))),
		YellowDrives:              float64(zeroOut(float64(len(dc.FilterByScore("yellow"))))),
		RedDrives:                 float64(zeroOut(float64(len(dc.FilterByScore("red"))))),
		AverageStaffAtEachDrive:   float64(zeroOut(dc.AverageStaffPerDrive())),
	}
	//get the drive reports
	for _, v := range dc.BloodDrives {
		data.DriveReports = append(data.DriveReports, v.GenerateReport())
	}
	// Create the report
	report := NewReport(dc.Name, "Donor Center Report", "This report contains information about the donor center", "Donor Center", data)

	Lock.Unlock()
	dc.Report = *report

	// Return the report
	return *report
}

type ConfigTemplate struct {
	TDY     Tdy     `bson:"tdy"`
	Inhouse InHouse `bson:"inhouse"`
	Sig     string  `bson:"signature"`
}

type Recruiter struct {
	FirstName string `bson:"firstname"`
	LastName  string `bson:"lastname"`
	OfficeNum string `bson:"officenumber"`
	Email     string `bson:"email"`
}

type Address struct {
	Street    string `bson:"street"`
	CityState string `bson:"citystate"`
	Zip       string `bson:"zip"`
	Region    string `bson:"region"`
}

type SocMedia struct {
	Site string `bson:"site"`
	Link string `bson:"link"`
}

type CentersJson struct {
	Centers []DonorCenter `bson:"centers"`
}

func (cjson *CentersJson) getCenterByName(name string) *DonorCenter {

	for _, v := range cjson.Centers {
		if name == v.Name {
			return &v
		}
	}

	return nil
}

func (dc *DonorCenter) getNum() int {
	return dc.Num
}

// Func getLocations returns an array with a list of BDC Names
func (cjson *CentersJson) getLocations() []string {
	var locations []string
	for _, v := range cjson.Centers {
		locations = append(locations, v.Name)
	}
	return locations
}

func (dc *DonorCenter) getDrive(drive string) (BloodDrive, bool) {

	d, ok := dc.BloodDrives[drive]
	if !ok {
		fmt.Println("drive not found....sorry")
		return BloodDrive{}, false
	}

	return d, true

}

func createTemplate(Kind string, collectionSupplies, processingSupplies, testing, averageCivilianRBC, averageCivilianFFPpF24, shipping, flightCosts, airportParking, lodging, mEPerDiem, hotelParking float32) ConfigTemplate {
	t := &Tdy{
		Kind:           Kind,
		Shipping:       shipping,
		FlightCosts:    flightCosts,
		AirportParking: airportParking,
		MEPerDiem:      mEPerDiem,
		HotelParking:   hotelParking,
	}
	ih := &InHouse{
		CollectionSupplies:     collectionSupplies,
		ProcessingSupplies:     processingSupplies,
		Testing:                testing,
		AverageCivilianRBC:     averageCivilianRBC,
		AverageCivilianFFPpF24: averageCivilianFFPpF24,
	}
	ct := &ConfigTemplate{
		TDY:     *t,
		Inhouse: *ih,
	}
	fmt.Println("template created: ", ct)
	return *ct
}

func (dc *DonorCenter) addDrive(bd BloodDrive) {
	Lock.Lock()
	if dc.BloodDrives == nil {
		dc.BloodDrives = make(map[string]BloodDrive)
		fmt.Println("no data..look----->: ", dc.BloodDrives)
	}
	if dc.BloodDrives != nil {
		//check to see if its already there
		_, ok := dc.BloodDrives[strings.Split(bd.creationTime, " ")[0]]
		if ok {
			fmt.Println("drive already exixts and will not be added")
			return
		}
		if !ok {
			dc.BloodDrives[strings.Split(bd.creationTime, " ")[0]] = bd
			//setBloodDrives(*dc)
		}
	}
	Lock.Unlock()
	//addDriveToCenter(*dc)
}

func (dc *DonorCenter) addStaff(staff Staff) {
	Lock.Lock()
	//check to see if map has been initialized
	if dc.StaffList == nil {
		dc.StaffList = make(map[string]Staff)

	}

	dc.StaffList[staff.Name] = staff
	Lock.Unlock()
	return
}

func (dc *DonorCenter) removeStaff(name string) {
	Lock.Lock()
	//check to see if map has been initialized
	if dc.StaffList == nil {
		dc.StaffList = make(map[string]Staff)

	}
	//check for staff member
	fmt.Println("looking for: ", name)
	if _, ok := dc.StaffList[name]; !ok {
		fmt.Println(name, "not found")
		return
	}

	//remove them from any blood drives they belong to
	for _, v := range dc.BloodDrives {
		if _, ok := v.Staff[name]; ok {
			fmt.Println(name, "removing : ", name)
			delete(v.Staff, name)

		}
	}

	delete(dc.StaffList, name)
	Lock.Unlock()
	fmt.Println("deleted staff member: ", name)
}

// HASTDY returns yes if there are any TDY drives in the center, and no if there are none
func (dc *DonorCenter) HASTDY() string {
	for _, v := range dc.BloodDrives {
		if v.DriveType.Kind == "TDY" {

			return fmt.Sprint("yes")
		}
	}
	return fmt.Sprint("no")
}

func (dc *DonorCenter) GetTDY() []BloodDrive {
	var tdy []BloodDrive
	for _, v := range dc.BloodDrives {
		if v.DriveType.Kind == "TDY" {
			tdy = append(tdy, v)
		}
	}
	return tdy
}

func (dc *DonorCenter) GetInHouse() []BloodDrive {
	var inHouse []BloodDrive
	for _, v := range dc.BloodDrives {
		if v.DriveType.Kind == IH {
			inHouse = append(inHouse, v)
		}
	}
	return inHouse
}

func (dc *DonorCenter) GetMobile() []BloodDrive {
	var mobile []BloodDrive
	for _, v := range dc.BloodDrives {
		if v.DriveType.Kind == "Mobile" {
			mobile = append(mobile, v)
		}
	}
	return mobile
}

func (dc *DonorCenter) FilterByScore(gyr string) []BloodDrive {
	var drives []BloodDrive
	gyr = strings.ToLower(gyr)
	switch gyr {
	case "green":
		for _, v := range dc.BloodDrives {
			percent := v.GetPercentage()
			if percent.Color == "green" {
				drives = append(drives, v)
			}
		}
		if len(drives) <= 0 {
			//return empty slice
			return drives

		}
		fmt.Println("NUMBER OF GREEN DRIVES \n\n", len(drives))
		break // Break statement for "green" case

	case "yellow":
		for _, v := range dc.BloodDrives {
			percent := v.GetPercentage()
			if percent.Color == "yellow" {
				drives = append(drives, v)
			}

		}
		if len(drives) <= 0 {
			return drives
		}
		fmt.Println("NUMBER OF YELLOW DRIVES \n\n", len(drives))
		break // Break statement for "yellow" case

	case "red":
		for _, v := range dc.BloodDrives {
			percent := v.GetPercentage()
			if percent.Color == "red" {
				drives = append(drives, v)
			}

		}
		if len(drives) <= 0 {
			return drives
		}
		fmt.Println("NUMBER OF RED DRIVES: \n\n", len(drives))
		break // Break statement for "red" case

	default:
		for _, v := range dc.BloodDrives {
			drives = append(drives, v)

		}
		break // Break statement for default case

	}
	return drives

}

func (dc *DonorCenter) FilterByType(imt string) []BloodDrive {
	var drives []BloodDrive
	imt = strings.ToLower(imt)
	switch imt {
	case "in-house":
		for _, v := range dc.BloodDrives {
			t := v.DriveType.Kind
			if t == "in-house" {
				drives = append(drives, v)
			}
		}
		if len(drives) <= 0 {
			return drives
		}
		fmt.Println("in-house: ", drives)
		break

	case "mobile":
		for _, v := range dc.BloodDrives {
			m := v.DriveType.Kind
			if m == "mobile" {
				drives = append(drives, v)
			}

		}
		if len(drives) <= 0 {
			return drives
		}
		fmt.Println("mobile drives: ", drives)
		break
	case "tdy":
		for _, v := range dc.BloodDrives {
			t := v.DriveType.Kind
			if t == "tdy" {
				drives = append(drives, v)
			}

		}
		if len(drives) <= 0 {
			return drives
		}
		fmt.Println("tdy drives: ", drives)
		break

	default:
		for _, v := range dc.BloodDrives {
			drives = append(drives, v)

		}

	}
	return drives

}

// write a function that uses the filterByTypeAndScore method to filter the blood drives by type and score. The function will take two parameters: the type of blood drive and the color of the blood drive. The function will return a list of blood drives that match the type and score. The function will be stored in the DonorCenter struct and will be returned as a JSON object.
func (dc *DonorCenter) FilterByTypeAndScore(imt, gyr string) map[string][]BloodDrive {
	drives := make(map[string][]BloodDrive)
	imt = strings.ToLower(imt)
	gyr = strings.ToLower(gyr)
	switch imt {
	case "in-house":
		var inHouseDrives []BloodDrive
		for _, v := range dc.BloodDrives {
			t := v.DriveType.Kind
			if t == "in-house" {
				percent := v.GetPercentage()
				if percent.Color == gyr {
					inHouseDrives = append(inHouseDrives, v)
				}
			}
		}
		if len(inHouseDrives) > 0 {
			drives[gyr] = inHouseDrives
		}
		fmt.Println("in-house drives: ", inHouseDrives)
		break

	case "mobile":
		var mobileDrives []BloodDrive
		for _, v := range dc.BloodDrives {
			m := v.DriveType.Kind
			if m == "mobile" {
				percent := v.GetPercentage()
				if percent.Color == gyr {
					mobileDrives = append(mobileDrives, v)
				}
			}
		}
		if len(mobileDrives) > 0 {
			drives[gyr] = mobileDrives
		}
		fmt.Println("mobile drives: ", mobileDrives)
		break

	case "tdy":
		var tdyDrives []BloodDrive
		for _, v := range dc.BloodDrives {
			t := v.DriveType.Kind
			if t == "tdy" {
				percent := v.GetPercentage()
				if percent.Color == gyr {
					tdyDrives = append(tdyDrives, v)
				}
			}
		}
		if len(tdyDrives) > 0 {
			drives[gyr] = tdyDrives
		}
		fmt.Println("tdy drives: ", tdyDrives)
		break

	default:
		var otherDrives []BloodDrive
		for _, v := range dc.BloodDrives {
			percent := v.GetPercentage()
			if percent.Color == gyr {
				otherDrives = append(otherDrives, v)
			}
		}
		if len(otherDrives) > 0 {
			drives[gyr] = otherDrives
		}
	}
	return drives
}
