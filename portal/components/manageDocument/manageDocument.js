function linkManageDocument(context) {
    "use strict";
    console.log("Running linkManageDocument");
    var subjectId = context.subjectId;
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
        var url = "/" + context.updateUrlComponent +  "/" + subjectId;
        console.log("Doing upload with url: " + url);
        var formData = new FormData();
        formData.append("file", file);
        post(url, formData)
            .then(onResponse)
            .catch(onResponse);
    }

    function post(url, data) {
        return axios.post(url, data)
            .then(function (response) {
                return response;
            }).catch(function (error) {
                return error.response;
            });
    }

    function onResponse(response) {
        if(response.status >= 400) {
            showResponse("error", response.data);
        } else {
            var theResponse = response.data;
            var responseClass = "success";
            if(theResponse.IsWarning) {
                responseClass = "error"
            }
            showResponse(responseClass, theResponse.Message);
            descriptionControl.innerHTML = theResponse.Description;
            setUploadNeeded(theResponse.UploadNeeded)
        }
    }

    function showResponse(responseClass, response) {
        alertControl.classList.remove("success");
        alertControl.classList.remove("error");
        alertControl.classList.add(responseClass);
        alertControl.innerHTML = response;
    }

    function verify() {
        var request = {
            Description: descriptionControl.innerHTML
        };
        post("/" + context.verifyUrlComponent + "/" + subjectId, JSON.stringify(request))
            .then(onResponse)
            .catch(onResponse);
    }
}
