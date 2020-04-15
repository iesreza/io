var App = {}

App.Alert = function (message,callback) {
    alert(message)
    if(callback){
        callback()
    }
}

IO.onAjaxError = function (error) {
    console.warn(error)
    IO.Loading(false)
    App.Alert("Unable to complete the request.")
}