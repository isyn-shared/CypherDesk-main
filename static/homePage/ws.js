const url = `ws://${window.location.host}/tickets/ws`;
const ws = new WebSocket(url);

function send(data) {
    console.log("Sending", data);
    ws.send(data)
}

/*let tickets = [{
	ticket: {
		ID: 1,
		Caption: "Сломался чайник в отделе 'Лох'",
		Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. \n\nUt enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
		Sender: "Иллейный Илля",
		Status: "opened"
    },
	forwardFrom: 1,
	forwardTo: 2,
    Action: "create",
	Time: "Когда-то"
},
{
	ticket: {
		ID: 2,
		Caption: "Сдать документы до пятого числа",
		Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. \nUt enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
		Sender: "Одмен Одменович",
		Status: "opened"
    },
	forwardFrom: 1,
	forwardTo: 2,
    Action: "create",
	Time: "Когда-то"
}];*/

let incomingCount = 0, sentCount = 0;

function makeForwardable() {
    $('a.forwardTicket').click(function (e) {
        e.stopPropagation();

        forwardTicket($(this).attr('ticketID'));
    });
}

function closeTicket(ticketID) {
    swal({
        title: 'Закрыть тикет?',
        text: 'Данное действие будет невозможно отменить',
        type: 'info',
        showCancelButton: true,
        confirmButtonColor: '#3085d6',
        cancelButtonColor: '#d33',
        confirmButtonText: 'Да, закрыть',
        cancelButtonText: 'Отмена'
    }).then(e => {
        if (e.value) {
            if (DEBUG) console.log(`Closing ${ticketID}`);

            sendEvent('close', {id: ticketID});
        }
    }).catch(console.error);
}

function makeClosable() {
    $('a.closeTicket').click(function (e) {
        closeTicket($(this).attr('ticketID'));
    });
}

function flash(elementID) {
    $('#tap-target').attr('data-target', elementID).tapTarget();
}

const myEvents = {
    // Special multi-purpose event 
    "error": (err) => { console.error(err); },
    "get": (tickets) => {
        console.log("Got tickets!", tickets);
        let sentLi = '';
        let incomingLi = '';

        let getSentCnt = 0, getIncomingCnt = 0;
        for (let i = tickets.length-1; i >= 0; i--) {
            const extTicket = tickets[i];

            extTicket.ticket.Description = extTicket.ticket.Description.replace(/(?:\r\n|\r|\n)/g, '<br>');
            // if (myUserData.ID == extTicket.ticket.Sender) extTicket.ticket.Sender = `${myUserData.Name} ${myUserData.Surname}`;
            if (userNames[extTicket.ticket.Sender]) extTicket.ticket.Sender = userNames[extTicket.ticket.Sender];
            else
                if (isUser)
                    extTicket.ticket.Sender = '<span class="green-text">Модератор</span>';
                else
                    extTicket.ticket.Sender = '<span class="green-text">Администратор</span>';

            let isSentLi = (extTicket.forwardFrom == myUserData.ID || extTicket.ticket.Sender == myUserData.ID);

            let ticket = prepareTicket({extTicket, isSentLi});


            if (isSentLi) {
                sentLi += ticket;
                getSentCnt++;
            }
            else {
                getIncomingCnt++;
                incomingLi += ticket;
            }
        }
        // if (incomingLi.length == 0) {
        //     incomingLi = `<h2 class="text-center">У вас нет тикетов!</h2>`
        //     incomingNeedsToBeEmpty = true;
        // }
        // if (sentLi.length == 0) {
        //     sentLi = `<h2 class="text-center">У вас нет тикетов!</h2>`
        //     sentNeedsToBeEmpty = true;
        // }

        $('.incomingTicketsUl').html(incomingLi);
        $('.sentTicketsUl').html(sentLi);
        makeForwardable();
        makeClosable();

        // Set counters
        /*$('#sentHeader').html(`<h5>У вас <b>${sentCount}</b> отправленных тикетов</h5>`);
        $('#incomingHeader').html(`<h5>У вас <b>${incomingCount}</b> полученных тикетов</h5>`);*/
        $('#sentCounter').html(getSentCnt);
        $('#receivedCounter').html(getIncomingCnt);

        $('#ticketAmountB').html(tickets.length);
    },
    "create": (extTicket) => {
        swal('Успешно!', 'Тикет был отправлен!', 'success');
        const isSentLi = true;

        if (DEBUG) console.log("Ticket was successfully created!", extTicket);

        // if (sentNeedsToBeEmpty) {
        //     $('.sentTicketsUl').html('');
        //     sentNeedsToBeEmpty = false;
        // }

        extTicket.ticket.Description = extTicket.ticket.Description.replace(/(?:\r\n|\r|\n)/g, '<br>');
        if (userNames[extTicket.ticket.Sender]) extTicket.ticket.Sender = userNames[extTicket.ticket.Sender];
        else 
            if (isUser)
                extTicket.ticket.Sender = '<span class="green-text">Модератор</span>';
            else
                extTicket.ticket.Sender = '<span class="green-text">Администратор</span>';

        let li = prepareTicket({ isSentLi, extTicket});
        li += $('.sentTicketsUl').html();

        $('.sentTicketsUl').html(li);
        makeForwardable();
        makeClosable();
        // Set counter
        $('#sentCounter').html(++sentCount).addClass('new');
    },
    "incoming": (extTicket) => {
        if (DEBUG) console.log("Got new ticket!", extTicket);

        M.toast({displayLength: 10000, html: `<span>Поступил новый тикет!</span><a class="btn-flat toast-action smoothScroll" href="#ticketsBlock" onclick="$('#receivedTicketsSpan').click()">Посмотреть</a>`});
        makeSmoothScrollable();

        const isSentLi = false;

        // if (incomingNeedsToBeEmpty) {
        //     $('.incomingTicketsUl').html('');
        //     incomingNeedsToBeEmpty = false;
        // }

        extTicket.ticket.Description = extTicket.ticket.Description.replace(/(?:\r\n|\r|\n)/g, '<br>');
        if (userNames[extTicket.ticket.Sender]) extTicket.ticket.Sender = userNames[extTicket.ticket.Sender];
        else
            if (isUser)
                extTicket.ticket.Sender = '<span class="green-text">Модератор</span>';
            else
                extTicket.ticket.Sender = '<span class="green-text">Администратор</span>';

        let li = prepareTicket({extTicket, isSentLi});
        li += $('.incomingTicketsUl').html();

        $('.incomingTicketsUl').html(li);
        makeForwardable();
        makeClosable();

        // Set counter
        $('#receivedCounter').html(++incomingCount).addClass('new');
    },
    "forward": (data) => {
        if (DEBUG) console.log("Forwarded!", data);

        swal('Успешно!', 'Тикет был перенаправлен!', 'success');
    },
    "close": (ticket) => {
        if (DEBUG) console.log("Closed!", ticket);

        M.toast({displayLength: 2000, html: `<span>Тикет «<b class="yellow-text">${ticket.Caption}</b>» был закрыт</span>`});
        $(`.ticketStatus${ticket.ID}`).html('Статус: <b>closed</b>');
    },
    "closedTicket": (ticket) => {
        if (DEBUG) console.log("Somebody closed:", ticket);

        M.toast({displayLength: 10000, html: `<span>Тикет «<b class="yellow-text">${ticket.Caption}</b>» был закрыт. </span><a class="btn-flat toast-action smoothScroll" href="#ticketsBlock" onclick="$('#sentTicketsSpan').click()">Посмотреть</a>`});
        makeSmoothScrollable();
        $(`.ticketStatus${ticket.ID}`).html('Статус: <b>closed</b>');
    },
    "publicKey": (serverCipherText) => {
        serverPublicKey = new NodeRSA();
        serverPublicKey = serverPublicKey.importKey(encryptionKey.decrypt(serverCipherText), 'pkcs1-public');

        if (DEBUG) console.log("Got server's publicKey:", serverPublicKey.exportKey('pkcs1-public'));
        sendEvent('get', {});
    }
}
// myEvents['get'](tickets);

function getTime(date) {
    return `${b( date.getHours() )}:${b( date.getMinutes() )} ${b( date.getDate() )}.${b( date.getMonth() + 1 )}`;
}
// Beautify
function b(n) {
    return n < 10 ? "0" + n : n;
}

ws.onmessage = (cipherText) => {
    let event = null;

    if (serverPublicKey) // If we already got server key we are switched to decrypt mode
        event = encryptionKey.decrypt(cipherText);
    else
        event = cipherText;

    console.log(event);
    let msg = JSON.parse(event.data);
    console.log(msg);

    if (msg.ok === false) {
        console.warn("Ошибка!", msg.data);
        swal('Упс!', `Что-то пошло не так: ${msg.data}`, 'error');
        return;
    }

    console.log(`Attempt to execute '${msg.event}' with data:`, msg.data);

    if (!myEvents[msg.event])
        return console.error("No event " + msg.event);

    myEvents[msg.event](JSON.parse(msg.data));
}

ws.onopen = () => {
    console.log("Connected successfully!");

    /* Sending initial events */
    sendEvent("publicKey", {key: encryptionKey.exportKey('pkcs1-public')});
}

ws.onclose = () => {
    console.log("Connection was closed")
}

function sendEvent(event, data) {
    let obj = { event, data: JSON.stringify(data) };
    let text = JSON.stringify(obj);

    if (serverPublicKey)
        text = serverPublicKey.encrypt(text);

    send(text);
}

function prepareTicket(info) {
    let { extTicket, isSentLi } = info;

    let ticket = `
        <li>
            <div class="collapsible-header">
                <i class="material-icons">folder</i>
                <span>${extTicket.ticket.Sender}: «<b>${extTicket.ticket.Caption}</b>»</span>
                
                ${!isUser && !isSentLi ? '<a class="waves-effect waves-light btn-small forwardTicket hide-on-small-and-down" style="position: absolute; right: 15px;" ticketID="' + extTicket.ticket.ID + '">' +
                    '<i class="material-icons left" style="margin: 0;">forward</i>' +
                    '<span>переслать</span>' +
                '</a>' : ''}
            </div>
            <div class="collapsible-body">
                <span>${extTicket.ticket.Description}</span>
                <br><br><br>
                <span class="right">Информация:</span><br><br>
                <div>
                    <span class="right">Время: ${getTime( new Date(extTicket.time) )}</span><br>
                    <span class="right">Отправитель: ${extTicket.ticket.Sender}</span><br>
                    <span class="right ticketStatus${extTicket.ticket.ID}">Статус: ${extTicket.ticket.Status}</span><br>
                </div>
               ${!isUser && !isSentLi ? '<a class="right waves-effect waves-light btn-small hide-on-med-and-up forwardTicket" ticketID="' + extTicket.ticket.ID + '" style="margin-top: 1rem; margin-bottom: 1rem">' +
                    '<span>переслать</span>' +
                '</a><br><br><br>' : ''}
               ${!isSentLi ? '<a class="right waves-effect waves-light btn-small red lighten-1 closeTicket" style="margin-top: 1rem;" ticketID="' + extTicket.ticket.ID + '">' +
                    '<span>закрыть тикет</span>' +
                '</a>' : ''}
               <br>
               <br>
            </div>
        </li>
    `;

    return ticket;
}