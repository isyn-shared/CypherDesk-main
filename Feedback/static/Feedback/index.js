$('#send_feedback_form').click(() => {
    let xhr = new XMLHttpRequest();

    let name = $('#user_name').text(),
        email = $('#user_email').text(),
        title = $('#message_title').text(),
        text = $('.pseudotextarea').text();

    console.log(name, email, title, text);

    let body = 'user_name=' + encodeURI(name) +
      '&user_email=' + encodeURI(email) +
      '&message_title=' + encodeURI(title) +
      '&message_text=' + encodeURI(text);
    
    xhr.open("POST", '/feedback/send/', true);
    xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');

    let csrftoken = Cookies.get('csrftoken');
    xhr.setRequestHeader("X-CSRFToken", csrftoken);
    
    xhr.onreadystatechange = event => {
        if (xhr.readyState == 4) {
            console.log(xhr.status, event.currentTarget.responseText);
            if (xhr.status != 200) return;

            // Do stuff
        }
    };
    
    xhr.send(body);
});