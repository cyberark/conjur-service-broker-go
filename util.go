package main

import "fmt"

const (
	expectedServiceID = "c024e536-6dc4-45c6-8a53-127e7f8275ab"
	expectedPlanID    = "3a116ac2-fc8b-496f-a715-e9a1b205d05c.community"
)

func formContext(context *Context, field string) *string {
	if context == nil {
		return nil
	}
	val, ok := (*context)[field]
	if !ok {
		return nil
	}
	res := fmt.Sprintf("%v", val)
	return &res
}

func validateServiceAndPlan(serviceID, planID string) error {
	if serviceID != expectedServiceID {
		return fmt.Errorf("invalid serviceID expected %v, got %v", expectedServiceID, serviceID)
	}
	if planID != expectedPlanID {
		return fmt.Errorf("invalid planID expected %v, got %v", expectedPlanID, planID)
	}
	return nil
}
