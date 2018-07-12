$('#send_authorization_form').click(() => {
    $("#send_authorization_form").attr('disabled', 'true');
    swal('Секундочку...', 'Проходит авторизация', 'info');

    let xhr = new XMLHttpRequest();

    let login = $("#login").val(),
        password = $("#password").val();

    let body = 'login=' + encodeURI(login) +
      '&password=' + encodeURI(password);

    xhr.open("POST", '/custom_admin/authorization/', true);
    xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');

    let csrftoken = Cookies.get('csrftoken');
    xhr.setRequestHeader("X-CSRFToken", csrftoken);


    xhr.onreadystatechange = event => {
        if (xhr.readyState == 4) {
            console.log(xhr.status, xhr.response);
            if (xhr.status != 200) return;

            let data = xhr.response;
            data = JSON.parse(data);

            if (data.authorization) {
                setTimeout(() => {
                    location = './';
                }, 1000);
            }
            else {
                if(data.type == 3){
                    swal("Упс!", "Вы уже авторизованы!", "fail");
                    setTimeout(() => {
                        location = './';
                    }, 5000);
                }
                else
                {
                    let error_mess;
                    if(data.type == 1)
                        error_mess = "Неправильный пароль! Попробуйте еще раз.";
                    if(data.type == 2)
                        error_mess = "Такого пользователя не существует! Проверьте логин.";
                    swal("Упс!", error_mess, "fail");
                    $("#send_authorization_form").attr('disabled', null);
                }
            }
        }
    };

    xhr.send(body);
});