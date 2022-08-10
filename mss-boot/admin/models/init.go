package models

func Init() {
	(&Tenant{}).InitStore()
}
