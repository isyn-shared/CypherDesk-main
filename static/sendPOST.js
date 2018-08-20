const xhr = new XMLHttpRequest();

function sendPOST(url, body = {}, convertToJSON = false) {
    return new Promise((resolve, reject) => {
        let formData = new FormData();
        for (let name in body)
            formData.append(name, body[name])

        xhr.open("POST", url, true);
        xhr.onreadystatechange = (event) => {
            if (xhr.readyState == XMLHttpRequest.DONE) {
                if (DEBUG) console.log(xhr.status, event.currentTarget.responseText);
                if (xhr.status != 200) return reject(event.currentTarget.responseText);

                let resp = event.currentTarget.responseText;
                if (convertToJSON) resp = JSON.parse(resp);

                resolve(resp);
            } 
        }
        xhr.send(formData);
    });
}