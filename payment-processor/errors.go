package main

import "errors"

var getResourceError = errors.New("Could not get resource")
var createResourceError = errors.New("Could not create resource")
var invalidResource = errors.New("invalid data")
var failedTransactionError = errors.New("(Fail) Transaction failed")
var failedError = errors.New("request failed")
var insufficientError = errors.New("insufficient funds!")
