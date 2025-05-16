package models



type Place struct {
	ID int `json:"id"`
	Name string `json:"name"`
}


// DistanceRequest represents the request body for creating a place distance
type DistanceRequest struct {
	FromPlaceID int `json:"from_place" binding:"required" example:"1" description:"ID of the starting place"`
	ToPlaceID   int `json:"to_place" binding:"required" example:"2" description:"ID of the destination place"`
	Distance    int `json:"distance" binding:"required,min=1" example:"100" description:"Distance between places in meters"`
}

type TaxistAnnouncmentRequest struct {
	DepartDate string `json:"depart_date" binding:"required" example:"2025-01-01" description:"depart date"`
	DepartTime string `json:"depart_time" binding:"required" example:"15:30:00" description:"depart time"`
	FromPlaceID int `json:"from_place" binding:"required" example:"0" description:"departure"`
	ToPlaceID int `json:"to_place" binding:"required" example:"0" description:"destination"`
	Space int `json:"space" binding:"required" example:"4" description:"how many seats are available"`
	Type AnnouncementType `json:"type" binding:"required" example:"submit one of person, package or person_and_package" description:"person, package or person_and_package"`
}

// DistanceResponse represents the response after creating a place distance
type DistanceResponse struct {
	ID          int `json:"id"`
	FromPlaceID int `json:"from_place"`
	ToPlaceID   int `json:"to_place"`
	Distance    int `json:"distance"`
}

type AnnouncementType string


// Define constants for the AnnouncementType
const (
    Person         AnnouncementType = "person"
    Package        AnnouncementType = "package"
    PersonAndPackage AnnouncementType = "person_and_package"
)

type TaxistAnnouncement struct {
	DepartDate string `json:"depart_date"`
	DepartTime string `json:"depart_time"`
	FromPlaceID int `json:"from_place"`
	ToPlaceID int `json:"to_place"`
	Space int `json:"space"`
	Type AnnouncementType `json:"type"` //person, package and person or package
}


type Ugur struct {
	ID int `json:"id"`
	TaxistID int `json:"taxist_id"`
	DepartDate string `json:"depart_date"`
	DepartTime string `json:"depart_time"`
	Space int `json:"space"`
	Distance int `json:"distance"`
	Type AnnouncementType `json:"type"`
	FullName string `json:"full_name"`
	CarMake string `json:"car_make"`
	CarModel string `json:"car_model"`
	CarYear int `json:"car_year"`
	CarNumber string `json:"car_number"`
	FromPlace string `json:"from_place"`
	ToPlace string `json:"to_place"`
	Rating float32 `json:"rating"`
}

type ReservePassengers struct {
	ID int `json:"id"`
	FullName string `json:"full_name"`
	Phone string `json:"phone"`
}

type UgurDetails struct {
	ID int `json:"id"`
	TaxistID int `json:"taxist_id"`
	DepartDate string `json:"depart_date"`
	DepartTime string `json:"depart_time"`
	Space int `json:"space"`
	Distance int `json:"distance"`
	Type AnnouncementType `json:"type"`
	FullName string `json:"full_name"`
	TaxistPhone string `json:"taxist_phone"`
	CarMake string `json:"car_make"`
	CarModel string `json:"car_model"`
	CarYear int `json:"car_year"`
	CarNumber string `json:"car_number"`
	FromPlace string `json:"from_place"`
	ToPlace string `json:"to_place"`
	Rating float32 `json:"rating"`
	Passengers []ReservePassengers `json:"passengers"`
}

type NotificationDetails struct {
	ID int `json:"id"`
	TaxistID int `json:"taxist_id"`
	FullName string `json:"full_name"`
	Phone string `json:"phone"`
	Package string `json:"package"`
	Count int `json:"count"`
	CreatedAt string `json:"created_at"`
	Passengers []ReservePassengers `json:"passengers"`
}

type FullReservePassengers struct {
	ID int `json:"id"`
	FullName string `json:"full_name"`
	Phone string `json:"phone"`
	Package string `json:"package"`
	AnnId int `json:"taxi_ann_id"`
	WhoReserved int `json:"who_reserved"`
}

type ReservePackages struct {
	PackageSender string `json:"package_sender"`
	PackageReciever string `json:"package_reciever"`
	SenderPhone string `json:"sender_phone"`
	RecieverPhone string `json:"reciever_phone"`
	AboutPackage string `json:"about_package"`
}

type Comment struct {
	ID int `json:"id"`
	FullName string `json:"full_name"`
	Comment string `json:"comment"`
}

type PassengerProfile struct {
	ID int `json:"id"`
	FullName string `json:"full_name"`
	Phone string `json:"phone"`
	UserType string `json:"user_type"`
}


type Notification struct {
	ID int `json:"id"`
	TaxistID int `json:"taxist_id"`
	FullName string `json:"full_name"`
	Package string `json:"package"`
	Count int `json:"count"`
	CreatedAt string `json:"created_at"`
}