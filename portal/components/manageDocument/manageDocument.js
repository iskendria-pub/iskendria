function linkManageDocument(context) {
    "use strict";
    var theHash = context.descriptionHash;
    var uploadTriggerControl = document.querySelector("#uploadTrigger");
    var verifyTriggerControl = document.querySelector("#verifyTrigger");
    var uploadControlControl = document.querySelector("#uploadControl");
    var verifyControlControl = document.querySelector("#verifyControl");
    var alertControl = document.querySelector('#alert');
    var descriptionControl = document.querySelector("#" + context.descriptionControlId);

    if(context.hasInitialDescriptionError) {
       context.alertControl.innerHTML = context.initialDescritionError;
       context.alertControl.classList.add("error");
    }
    uploadTriggerControl.addEventListener("change", addFile);
    verifyTriggerControl.addEventListener("click", verify);
    setUploadNeeded(context.initialIsUploadNeeded);

    function setUploadNeeded(uploadNeeded) {
        console.log("setUploadNeeded: " + uploadNeeded);
        if(uploadNeeded) {
            uploadControlControl.hidden = false;
            verifyControlControl.hidden = true;
        } else {
            uploadControlControl.hidden = true;
            verifyControlControl.hidden = false;
        }
        console.log("done setUploadNeeded")
    }

    function addFile(e) {
        var file = e.target.files[0];
        if(!file){
            return
        }
        upload(file);
    }

    function upload(file) {
        var formData = new FormData();
        formData.append("file", file);
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
            showResponse("error", response.data);
        } else {
            var theResponse = response.data;
            showResponse("success", theResponse.Message);
            setUploadNeeded(false);
            descriptionControl.innerHTML = theResponse.Text;
        }
    }

    function showResponse(responseClass, response) {
        alertControl.classList.remove("success");
        alertControl.classList.remove("error");
        alertControl.classList.add(responseClass);
        alertControl.innerHTML = response;
    }

    function verify() {
        post("/verify/" + theHash)
            .then(onVerifyResponse)
            .catch(onVerifyResponse)
    }

    function onVerifyResponse(response) {
        if(response.status >= 400) {
            showResponse("error", response.data);
            setUploadNeeded(true)
        } else {
            showResponse("success", response.data);
        }
    }
}
