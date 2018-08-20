function login() {
    const login = $('#loginInput').val(),
        pass = $('#passwordInput').val();
    if (DEBUG) console.log(login, pass);

    createAlert('alert-info', 'Пожалуйста, подождите...')

    sendPOST('/authorize', {login, pass}, true)
        .then(resp => {
            if (!resp.ok)
                return createAlert('alert-danger', 'Ошибка!', resp.err);

            createAlert('alert-success', 'Успех!', 'Все прошло успешно!');
            location = "/account";
        })
        .catch(console.error);

    return false;
}

function createAlert(type, title, text = "") {
    $('#alertWrapper').html(`
        <div class="alert alert-dismissible fade show ${type} mb-0" role="alert">
            <button type="button" class="close" data-dismiss="alert" aria-label="Закрыть">
                <span aria-hidden="true">×</span>
            </button>
            <strong>${title}</strong> ${text}
        </div>
    `);
}