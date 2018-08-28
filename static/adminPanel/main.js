let selectedWindow = "#status";

function showWindow(windowID) {
    if (selectedWindow != windowID) {
        if (DEBUG) console.log(windowID);


        let prevSelected = selectedWindow;
        selectedWindow = windowID;

        // This piece of code lights up navbar elements
        // In order to get <a> from window id you need to add "A" to the end
        // i.e. #status is a window, while #statusA is a corresponding <a>
        $(prevSelected + "A").removeClass('active');
        $(selectedWindow + "A").addClass('active');

        $(prevSelected).removeClass('slideIn').addClass('slideAway');
        // We need to place that "out" only when we do "slideAway".
        // If we don't remove that event, it will trigger in all different cases
        $(prevSelected).on('animationend', () =>
            $(prevSelected).addClass('out').removeClass('slideAway').off('animationend'));

        $(selectedWindow).removeClass('out').addClass('slideIn');
    }
}

$(document).ready(() => {
    $('#searchInput').keyup(event => {
        if (event.key)
            search(false);
    });

    $('#addDepForm').submit(e => {
        e.preventDefault();

        const name = $('#depNameInput').val();

        sendPOST('/createDepartment', {name})
            .then(e => {
                createAlert('alert-success', "Отлично!", "Сервер успешно получил данные, хз если все окей")
            })
            .catch(err => {
                createAlert('alert-danger', "Упс!", "Произошла ошибка: " + err);
                console.error(err);
            });
    });

    $('#addUserForm').submit(e => {
        e.preventDefault();

        const mail = $('#userEmailInput').val(),
            role = $('select[name=roleSelect]').val(),
            department = $('select[name=departmentSelect]').val();

        if (DEBUG) console.log(mail, role, department);
        if (role == "0" || department == "0")
            return createAlert('alert-danger', 'Упс!', 'Убедитесь, что выбрали все!')
    
        sendPOST('/createUser', {mail, role, department})
            .then(e => {
                createAlert('alert-success', "Отлично!", "Сервер успешно получил данные, хз если все окей")
            })
            .catch(err => {
                createAlert('alert-danger', "Упс!", "Произошла ошибка: " + err);
                console.error(err);
            });
    });
});

let timer = null;
function search(bypass = true) {
    console.log(bypass);
    const key = $('#searchInput').val() || "*";

    if (!bypass) {

        if (timer) clearTimeout(timer);

        timer = setTimeout(() => {
            action();
        }, 200);
    } else 
        action();

    function action() {
        sendPOST('/findUser', {key})
            .then(console.log)
            .catch(console.error);
    }
}

function createAlert(type, title, text = "") {
    $('.alertWrapper').html(`
        <div class="alert alert-dismissible fade show ${type} mb-0" role="alert">
            <button type="button" class="close" data-dismiss="alert" aria-label="Закрыть">
                <span aria-hidden="true">×</span>
            </button>
            <strong>${title}</strong> ${text}
        </div>
    `);

    const offset = $('#alertWrapper').offset();

    $('html, body').animate({
        scrollTop: offset.top,
        scrollLeft: offset.left
    });
}