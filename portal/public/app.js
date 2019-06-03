function linkUploadForm(theHash, context) {
    "use strict";
    if(context.hasInitialDescriptionError) {
       context.alertControl.innerHTML = context.initialDescritionError
       context.alertControl.classList.add("error")
    }
    context.inputFileControl.addEventListener("change", addFile);
    context.verifyButtonControl.addEventListener("click", verify)
    setUploadNeeded(context.initialIsUploadNeeded)

    function setUploadNeeded(uploadNeeded) {
        console.log("setUploadNeeded: " + uploadNeeded)
        if(uploadNeeded) {
            context.uploadFormControl.hidden = false
            context.verifyControl.hidden = true
        } else {
            context.uploadFormControl.hidden = true
            context.verifyControl.hidden = false
        }
        console.log("done setUploadNeeded")
    }

    function addFile(e) {
        var file = e.target.files[0]
        if(!file){
            return
        }
        upload(file);
    }

    function upload(file) {
        var formData = new FormData()
        formData.append("file", file)
        post("/upload/" + theHash, formData)
            .then(onUploadResponse)
            .catch(onUploadResponse);
    }

    function post(url, data) {
        return axios.post(url, data)
            .then(function (response) {
                return response;
            }).catch(function (error) {
                return error.response;
            });
    }

    function onUploadResponse(response) {
        if(response.status >= 400) {
            showResponse("error", response)
        } else {
            showResponse("success", response)
            setUploadNeeded(false)
        }
    }

    function showResponse(responseClass, response) {
        context.alertControl.classList.remove("success")
        context.alertControl.classList.remove("error")
        context.alertControl.classList.add(responseClass)
        context.alertControl.innerHTML = response.data;
    }

    function verify() {
        post("/verify/" + theHash)
            .then(onVerifyResponse)
            .catch(onVerifyResponse)
    }

    function onVerifyResponse(response) {
        if(response.status >= 400) {
            showResponse("error", response)
            setUploadNeeded(true)
        } else {
            showResponse("success", response)
        }
    }
}
