var urlTextfield = document.getElementById("url-textfield")

var shortenButton = document.getElementById("submit-button")
shortenButton.addEventListener("click", onSubmitBtnClicked)

function onSubmitBtnClicked() {
    urlValue = urlTextfield.value
    reqBody = { url: urlValue }

    fetch("https://urlshortener-production-f981.up.railway.app/api/v1/url/create", {
        method: 'POST',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(reqBody)
    }).
        then(res => res.json()).
        then(json => {
            if (json.status_code != 201) {
                urlTextfield.value = "URL isn't valid"
            } else {
                urlTextfield.value = "https://shrt.go/" + json.data.short_url
            }
        })
}

console.log(document.cookie)