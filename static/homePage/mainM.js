const isUser = false;

$(window).on('scroll', () => {
    let s = $(window).scrollTop(),
        d = $(document).height(),
        c = $(window).height(),
        w = $(window).width();

    // If height is greater than width
    if (c > w) 
        return $("#main").css("opacity", 1);

    let scrolledArea = (s / (d - c));
    // console.log(s, d, c, scrolledArea);
    $("#main").css("opacity", scrolledArea);
    $("body").css("background-color", `rgb(${242 + scrolledArea * 13}, ${252 + scrolledArea * 3}, ${255 - scrolledArea * 13})`);
});

$('#sendTicketForm').submit(e => {
    e.preventDefault();

    const caption = $('#ticketCaptionInput').val(),
        description = $('#ticketDesc').val();

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