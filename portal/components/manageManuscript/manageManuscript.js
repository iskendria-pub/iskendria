function linkManageDocument(context) {
    "use strict";
    console.log("Running linkManageDocument");
    var subjectId = context.subjectId;
    var uploadTriggerControl = document.querySelector("#uploadTrigger");
    var alertControl = document.querySelector('#alert');
    var downloadControlId = context.downloadControlId
    var downloadControl = document.querySelector("#" + downloadControlId)

    uploadTriggerControl.addEventListener("change", addFile);
    setUploadNeeded(context.initialIsUploadNeeded);

    function setUploadNeeded(uploadNeeded) {
        console.log("setUploadNeeded: " + uploadNeeded);
        if(uploadNeeded) {
            downloadControl.disabled = true
        } else {
            downloadControl.disabled = false
        }
        console.log("done setUploadNeeded")
    }

    function addFile(e) {
        console.log("Doing addFile")
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
            console.log("Got successful response: " + theResponse)
            var theResponse = response.data;
            var responseClass = "success";
            if(theResponse.IsWarning) {
                responseClass = "error"
            }
            showResponse(responseClass, theResponse.Message);
            setUploadNeeded(theResponse.UploadNeeded)
        }
    }

    function showResponse(responseClass, response) {
        alertControl.classList.remove("success");
        alertControl.classList.remove("error");
        alertControl.classList.add(responseClass);
        alertControl.innerHTML = response;
    }
}
