package svc

//package base service

import (
	helper "beego-mgo-rabbitmq-jwt/utilities/helper"

	"beego-mgo-rabbitmq-jwt/utilities/mgodb"

	log "github.com/goinggo/tracelog"
	"gopkg.in/mgo.v2"
)

type (
	// Service contains common properties for all services.
	Service struct {
		MongoSession *mgo.Session //every service has a db connection
		UserID       string       //for log
		UserName     string       //for log like nickname
	}
)

// get a session when a service inited
func (service *Service) Prepare() (err error) {
	if service.UserID == "" {
		service.UserID = "unknown"
	}
	service.MongoSession, err = mgodb.CopyMasterSession(service.UserID)
	if err != nil {
		log.Error(err, "UserID:"+service.UserID+" UserName:"+service.UserName, "Service.Prepare")
		return err
	}

	return err
}

// Finish is called after the controller.
// release the db connection.
func (service *Service) Finish() (err error) {
	defer helper.CatchPanic(&err, service.UserID, "Service.Finish")

	if service.MongoSession != nil {
		mgodb.CloseSession(service.UserID, service.MongoSession)
		service.MongoSession = nil
	}
	return err
}

// DBAction executes the MongoDB literal function
func (service *Service) DBAction(databaseName string, collectionName string, dbCall mgodb.DBCall) (err error) {
	return mgodb.Execute(service.UserID, service.MongoSession, databaseName, collectionName, dbCall)
}
