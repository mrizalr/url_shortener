var urlTextfield = document.getElementById("url-textfield")

var shortenButton = document.getElementById("submit-button")
shortenButton.addEventListener("click", onSubmitBtnClicked)

function onSubmitBtnClicked(){
    urlValue = urlTextfield.value
    reqBody = { url: urlValue}
    
}

