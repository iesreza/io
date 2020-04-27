package admin

type Settings struct {
	SessionAge string `type:"select" col:"6" label:"Session Age" hint:"JWT Token Age" options:"1800:30 Mins,3600:1 Hour,86400:24 Hours,604800:1 Week,2592000:1 Month,31536000:1 Year"`
}
