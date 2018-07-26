$('#send_feedback_form').click(() => {
    let name = $('#user_name').val(),
        email = $('#user_email').val(),
        title = $('#message_title').val(),
        text = $('.pseudotextarea').text();

    if (!name || !email || !title || !text)
        return swal("Постойте!", "Заполните все поля перед отправкой", "warning");

    if (!validateEmail(email))
        return swal("Постойте!", "Вы уверены в правильном написании Вашей почты?", "warning");


    $("#send_feedback_form").attr('disabled', 'true');
    swal('Спасибо!', 'Ваш запрос был отправлен серверу на рассмотрение. Пожалуйста, подождите (Обычно это занимает около минуты).', 'info');

    let xhr = new XMLHttpRequest();

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
            console.log(xhr.status, xhr.response);
            if (xhr.status != 200) return;

            // 0 - success
            // 1 - pochta
            // 2 - stop the spam man

            if (xhr.response == "0") {
                swal("Отлично!", "Все прошло успешно и команда CypherDesk получила Ваш запрос! (Вы будете перенаправлены на главную через пять секунд)", "success").then(() => {
                    location = '../';
                });
                setTimeout(() => {
                    location = '../';
                }, 5000);
            } else if (xhr.response == "1") {
                swal("Упс!", "Что-то пошло совсем не так! Мы не смогли отправить Вам почту. Попробуйте отправить запрос еще раз", "error");
                $("#send_feedback_form").attr('disabled', null);
            } else if (xhr.response == "2") {
                swal("Упс!", "С Вашего IP адреса уже была совершена отправка в течении этих двух часов. Если вы считаете, что это ошибка, пожалуйста, напишите нам:\n http://cypherdesk.ru/#contacts", "error");
            }
        }
    };

    xhr.send(body);
});


function validateEmail(email) {
    var re = /^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
    return re.test(String(email).toLowerCase());
}