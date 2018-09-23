$(window).on('scroll', () => {
    let s = $(window).scrollTop(),
        d = $(document).height(),
        c = $(window).height();

    let scrolledArea = (s / (d - c));
    // console.log(s, d, c, scrolledArea);
    $("#main").css("opacity", scrolledArea);
    $("body").css("background-color", `rgb(${242 + scrolledArea * 13}, ${252 + scrolledArea * 3}, ${255 - scrolledArea * 13})`);
});

$('#sendTicketForm').submit(e => {
    e.preventDefault();

    const caption = $('#ticketCaptionInput').val(),
        description = $('#ticketDesc').val();

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
            $('#sendTicketModal').modal('hide');
        }
    });
});