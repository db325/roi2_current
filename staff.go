package main

import (
	"fmt"
	"sync"
)

const (

	//Person Types
	Mil = "Military"
	Civ = "Civilian"
	Ctr = "Contractor"
	//GroupTypes
	Phlem   = "Phlebotomist"
	Supp    = "Support Staff"
	Lab     = "Lab Staff"
	LabTech = "Lab Tech"
	Med     = "Med Tech"
	Lead    = "Lead Tech"
	Rec     = "Recruiter"
	Super   = "Supervisor"
)

type Staff struct {
	r          *sync.RWMutex               `json:"-" bson:"-"`
	Name       string                      `json:"name" bson:"name"`
	PayRate    float32                     `json:"payrate" bson:"payrate"`
	DriveHours map[string]DriveHoursWorked `json:"drivehours" bson:"drivehours"`
	Active     bool                        `json:"active" bson:"active"`
	Center     int                         `json:"center" bson:"center"`
	Report     Report                      `json:"reports" bson:"report"`
}

type DriveHoursWorked struct {
	HrsCollecting float32 `json:"hrscollecting" bson:"hrscollecting"`
	HrsProcessing float32 `json:"hrsprocessing" bson:"hrsprocessing"`
	HrsSupporting float32 `json:"hrssupporting" bson:"hrssupporting"`
}

type Group struct {
	Staff     []*Staff `json:"staff" bson:"staff"`
	GroupType string   `json:"grouptype" bson:"grouptype"`
}

func makeStaff(name string, pay float32, center int) Staff {

	s := &Staff{
		Name:       name,
		PayRate:    pay,
		DriveHours: make(map[string]DriveHoursWorked),
		Active:     true,
		Center:     center,
		Report:     Report{}}

	return *s
}

func GetDriveHoursValues(staff Staff) []float64 {
	var values []float64
	for _, hours := range staff.DriveHours {
		values = append(values, float64(hours.HrsCollecting), float64(hours.HrsProcessing), float64(hours.HrsSupporting))
	}
	fmt.Println("VALUES: ", values)
	return values
}

func (s Staff) HaveIWorked() bool {
	if len(s.DriveHours) == 0 {
		return false
	}
	fmt.Println("i have hours", s.DriveHours)
	return true

}

func (s Staff) DoIBelong(drivedate string) bool {
	_, ok := s.DriveHours[drivedate]
	if ok {
		fmt.Println(ok, "YES")
		return true
	} else {
		return false
	}
}

func (s Staff) addDrive(date string) {
	//s.r.Lock()
	if v, found := s.DriveHours[date]; found {
		fmt.Println("Already added...", v)

	}

	s.DriveHours[date] = DriveHoursWorked{}
	//s.r.Unlock()
	return
}

func (s Staff) removeDrive(date string) {
	fmt.Println("INSIDE REMOVE DRIVE ...\n\n\n", date)
	//s.r.Lock()
	if dhw, found := s.DriveHours[date]; found {
		fmt.Println("WE FOUND IT: ", dhw, "\n\n")
		delete(s.DriveHours, date)

	} else {
		fmt.Println("NOT HERE!!\n\n\n\n\n")
		return
	}

	//s.DriveHours[date] = DriveHoursWorked{}
	//s.r.Unlock()
}

func (s Staff) removeAllDrives() {
	Lock.Lock()
	s.DriveHours = make(map[string]DriveHoursWorked)
	Lock.Unlock()
}

func (s Staff) toggle() bool {
	if s.Active == true {
		s.Active = false
	} else {
		s.Active = false
	}
	return s.Active
}

func (a *Staff) GenerateReport() Report {
	//check for drive hours first
	if !a.HaveIWorked() {
		fmt.Println("NO HOURS WORKED")
		return Report{}
	}
	values := GetDriveHoursValues(*a)
	dhw := struct {
		Collecting float64
		Processing float64
		Supporting float64
	}{values[0], values[1], values[2]}
	data := struct {
		Name             string
		PayRate          float32
		DriveHoursWorked interface{}
		Active           bool
		Center           int
	}{
		Name:             a.Name,
		PayRate:          a.PayRate,
		DriveHoursWorked: dhw,
		Active:           a.Active,
		Center:           a.Center,
	}
	r := NewReport(a.Name, "Staff Report", "Staff Report", "Staff", data)
	a.Report = *r
	return *r
}
