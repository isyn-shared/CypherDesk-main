$('#send_feedback_form').click(() => {
    let xhr = new XMLHttpRequest();

    let name = $('#user_name').val(),
        email = $('#user_email').val(),
        title = $('#message_title').val(),
        text = $('.pseudotextarea').text();

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
            if (xhr.status != 200) return;
        }
    };
    
    xhr.send(body);
});