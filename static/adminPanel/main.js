let selectedWindow = Cookies.get('window') || "#profile";

$(selectedWindow).removeClass('out');
$(selectedWindow + "A").addClass('active');

function showWindow(windowID) {
    if (selectedWindow != windowID) {
        if (DEBUG) console.log(windowID);

        Cookies.set('window', windowID);

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

        sendPOST('/admin/createDepartment', {name})
            .then(resp => {
                if (!resp.ok) 
                    return createAlert('alert-danger', "Упс!", "Произошла ошибка: " + resp.err);

                createAlert('alert-success', "Отлично!", "Отдел создан!");
                location.reload();
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
            return createAlert('alert-danger', 'Упс!', 'Убедитесь, что выбрали все!');

        createAlert('alert-info', "Загрузка...", "Пожалуйста подождите");
    
        sendPOST('/admin/createUser', {mail, role, department})
            .then(resp => {
                // if (DEBUG) console.log(resp, typeof resp);
                if (!resp.ok) 
                    return createAlert('alert-danger', "Упс!", "Произошла ошибка: " + resp.err);

                createAlert('alert-success', "Отлично!", "Пользователь создан!");
                $('#refreshUsers').click();
            })
            .catch(err => {
                createAlert('alert-danger', "Упс!", "Произошла ошибка: " + err);
                console.error(err);
            });
    });

    $('#editUserForm').submit(e => {
        e.preventDefault();

        const newLogin = $('#editUserLoginInput').val(), 
            newName = $('#editUserNameInput').val(),
            newSurname = $('#editUserSurnameInput').val(),
            newPartonymic = $('#editUserPartonymicInput').val(),
            newRecourse = $('#editUserRecourseInput').val(),
            newDepartment = $('#editUserDepartmentSelect').val();

        if (DEBUG) 
            console.log(editingUserLogin, newLogin, newSurname, newPartonymic, newRecourse, newDepartment)

        if (newDepartment == "0")
            return createAlert('alert-danger', 'Упс!', 'Убедитесь, что выбрали отдел!')

        sendPOST('/admin/changeUser', {login: editingUserLogin, newLogin, newName, newSurname, newPartonymic, newRecourse, newDepartment})
            .then(resp => {
                if (!resp.ok) 
                    return createAlert('alert-danger', "Упс!", "Произошла ошибка: " + resp.err);

                createAlert('alert-success', "Отлично!", "Пользователь изменен!");
                $('#refreshUsers').click();
            })
            .catch(err => {
                createAlert('alert-danger', "Упс!", "Произошла ошибка: " + err);
                console.error(err);
            });
    });

    $('#editDepForm').submit(e => {
        e.preventDefault();

        const name = $('#editDepNameInput').val();
        sendPOST('/admin/changeDepartment', {prevName: editingDepName, name})
            .then(resp => {
                // if (DEBUG) console.log(resp, typeof resp);
                if (!resp.ok) 
                    return createAlert('alert-danger', "Упс!", "Произошла ошибка: " + resp.err);

                createAlert('alert-success', "Отлично!", "Отредактировано!");
                location.reload();
            })
            .catch(err => {
                createAlert('alert-danger', "Упс!", "Произошла ошибка: " + err);
                console.error(err);
            });
    })
});

let timer = null;
function search(bypass = true) {
    console.log(bypass);
    const key = $('#searchInput').val() || "*";

    if (!bypass) {

        if (timer) clearTimeout(timer);

        timer = setTimeout(() => {
            findUsers(key);
        }, 1000);
    } else 
        findUsers(key);
}

let autoCloseTimer = null;
function createAlert(type, title, text = "") {
    clearTimeout(autoCloseTimer);
    $('.alertWrapper').html(`
        <div class="alert alert-dismissible fade show ${type} mb-0" role="alert">
            <button id="dismissAlert" type="button" class="close" data-dismiss="alert" aria-label="Закрыть">
                <span aria-hidden="true">×</span>
            </button>
            <strong>${title}</strong> ${text}
        </div>
    `);
    autoCloseTimer = setTimeout(() => {
        $(".alert").alert('close');
    }, 3000);
}


/* Vue.JS */
// Render all users by making one call ourselves
findUsers("*");

let editingUserLogin = "";

function editUser(user) {
    if (DEBUG) console.log("Editing", user);

    // Save for sending to server in form submit
    editingUserLogin = user.Login;

    $('#editUserLoginInput').val(user.Login)
    $('#editUserNameInput').val(user.Name);
    $('#editUserSurnameInput').val(user.Surname);
    $('#editUserPartonymicInput').val(user.Partonymic);
    $('#editUserRecourseInput').val(user.Recourse);
    $('#editUserDepartmentSelect').val(user.Department);
}

const app = new Vue({
    el: '#renderedUsers',
    data: {
        users: [],
        editUser
    },
});
function findUsers(key) {
    sendPOST('/admin/findUser', {key})
        .then(users => {

            console.log(users);

            for (user of users) {
                // Departments were defined inside html page
                user.DepartmentName = departments[user.Department];
                user.Tag = "bg-primary";

                if (!user.Name) {
                    user.Name = "Пользователь не активирован";
                    user.Recourse = "Неопределено";
                    user.Tag = "bg-secondary";
                }

            }

            app.users = users;

            if (!users.length)
                $('#usersNotFound').removeClass('d-none');
            else
                $('#usersNotFound').addClass('d-none');
        })
        .catch(console.error);
}

function deleteUser() {
    swal({
        title: `Удалить ${editingUserLogin}?`,
        text: "Это действие невозможно отменить!",
        type: 'warning',
        showCancelButton: true,
        confirmButtonColor: '#3085d6',
        cancelButtonColor: '#d33',
        confirmButtonText: 'Да, удалить!',
        cancelButtonText: 'Отмена',
        allowOutsideClick: () => !swal.isLoading(),
        preConfirm: () => {
            return sendPOST('/admin/deleteUser', {login: editingUserLogin})
                .then(response => {console.log(response); return response})
                .catch(error => {
                    swal.showValidationError(`Запрос провалился: ${error}`);
                    console.error(error);
                })
        }
    }).then(result => {
        if (result.value) {
            swal(
                'Успех!',
                'Удалено!',
                'success'
            );

            $('#refreshUsers').click();
        }
    }).catch(console.error);
}

let editingDepName = null;
function editDep(name) {
    editingDepName = name;

    $("#editDepNameInput").val(name);
}

function deleteDep() {
    swal({
        title: `Удалить ${editingDepName}?`,
        text: "Это действие невозможно отменить!",
        type: 'warning',
        showCancelButton: true,
        confirmButtonColor: '#3085d6',
        cancelButtonColor: '#d33',
        confirmButtonText: 'Да, удалить!',
        cancelButtonText: 'Отмена',
        allowOutsideClick: () => !swal.isLoading(),
        preConfirm: () => {
            return sendPOST('/admin/deleteDepartment', {name: editingDepName})
                .then(response => {console.log(response); return response})
                .catch(error => {
                    swal.showValidationError(`Запрос провалился: ${error}`);
                    console.error(error);
                })
        }
    }).then(result => {
        if (result.value) {
            swal(
                'Успех!',
                'Удалено!',
                'success'
            );
            location.reload();
        }
    }).catch(console.error);
}