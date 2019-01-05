const isUser = true;

$(window).on('scroll', () => {
    let s = $(window).scroliflTop(),
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

    if (DEBUG) console.log(caption, description);

    if (!caption) return swal('Упс!', 'Вы не заполнили заголовок тикета!', 'info');
    // Set length in html
    if ($('#ticketCaptionInput').attr('data-length') < caption.length) return swal('Упс!', 'Ваш заголовок слишком длинный для отправки!', 'info');
    
    if (!description) return swal('Упс!', 'Вы не заполнили описание проблемы!', 'info');
    if ($('#ticketDesc').attr('data-length') < caption.length) return swal('Упс!', 'Ваш заголовок слишком длинный для отправки!', 'info');

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
            sendEvent('create', {caption, description});
            
            // Empty inputs
            $('#ticketCaptionInput').val('');
            $('#ticketDesc').val('');
        }
    });
});