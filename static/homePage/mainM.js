const isUser = false;

$('#sendTicketForm').submit(e => {
    e.preventDefault();

    const caption = $('#ticketCaptionInput').val(),
        description = $('#ticketDesc').val(),
        selectedUser = $('#userSelect').val();
    
    // console.log(selectedUser);

    if (!caption) return swal('Упс!', 'Вы не заполнили заголовок тикета!', 'info');
    // Set length in html
    if ($('#ticketCaptionInput').attr('data-length') < caption.length) return swal('Упс!', 'Ваш заголовок слишком длинный для отправки!', 'info');
    
    if (!description) return swal('Упс!', 'Вы не заполнили описание проблемы!', 'info');
    if ($('#ticketDesc').attr('data-length') < caption.length) return swal('Упс!', 'Ваш заголовок слишком длинный для отправки!', 'info');
    
    if (!selectedUser) return swal('Упс!', 'Вы не выбрали пользователя, которому отправляется тикет!', 'info');

    // TODO: Сообщение, куда, что и зачем
    swal({
        title: 'Отправить тикет?',
        text: "Это действие нельзя будет отменить",
        type: 'info',
        showCancelButton: true,
        confirmButtonColor: '#3085d6',
        cancelButtonColor: '#d33',
        confirmButtonText: 'Да, отправить',
        cancelButtonText: 'Отмена'
    }).then(obj => {
        if (obj.value) {
            sendEvent('createM', {caption, description, id: selectedUser});
            
            // Empty inputs
            $('#ticketCaptionInput').val('');
            $('#ticketDesc').val('');
        }
    });
});

let selectedUser = null;
let lastSelectedHtml = null;
function selectSendUser(id, htmlElement) {
    selectedUser = id;

    if (lastSelectedHtml) lastSelectedHtml.classList.remove('selected');

    htmlElement.classList.add('selected');
    lastSelectedHtml = htmlElement;
}

function forwardTicket(ticketID) {
    let html = `
        <div class="input-field col m6 s12" id="userSelectFieldSwal" style="margin-bottom: 20px">
            <select id="userSelectSwal">
                <option value="-1" disabled selected>Выберите пользователя</option>
    `;

    for (user of usersWithID)
        html += `<option value="${user.id}">${user.name}</option>\n`;

    html += `</select></div>`;


    swal({
        title: 'Переслать тикет?',
        type: 'info',
        html: html,
        showCloseButton: true,
        confirmButtonColor: '#3085d6',
        confirmButtonText: "Отправить!",
        cancelButtonText: "Отмена",
        showCancelButton: true,
        focusConfirm: false,
    }).then(obj => {
        if (obj.value) {
            let toID = $('#userSelectSwal').val();
            if (!toID) return swal('Упс!', 'Вы не выбрали пользователя!', 'info');

            sendEvent('forward', {ticketID: ticketID.toString(), to: toID});
        }
    });

    $('#userSelectSwal').formSelect();
    // const instance = M.FormSelect.getInstance($('#userSelectSwal')[0]);
    // instance.options.
}