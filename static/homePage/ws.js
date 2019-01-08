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

// If it turns out to be empty, in 'create' event we need to clear UL
var sentNeedsToBeEmpty = false,
    incomingNeedsToBeEmpty = false;

let incomingCount = 0, sentCount = 0;
const myEvents = {
    // Special multi-purpose event 
    "error": (err) => { console.error(err); },
    "get": (tickets) => {
        console.log("Got tickets!", tickets);
        let sentLi = '';
        let incomingLi = '';
        for (let extTicket of tickets) {
            extTicket.ticket.Description = extTicket.ticket.Description.replace(/(?:\r\n|\r|\n)/g, '<br>');
            // if (myUserData.ID == extTicket.ticket.Sender) extTicket.ticket.Sender = `${myUserData.Name} ${myUserData.Surname}`;
            if (userNames[extTicket.ticket.Sender]) extTicket.ticket.Sender = userNames[extTicket.ticket.Sender];
            else extTicket.ticket.Sender = '<span class="green-text">Модератор</span>';

            let isSentLi = (extTicket.forwardFrom == myUserData.ID || extTicket.ticket.Sender == myUserData.ID);

            let ticket = `
                <li>
                    <div class="collapsible-header"><i class="material-icons">folder</i>${extTicket.ticket.Sender}: «<b>${extTicket.ticket.Caption}</b>»</div>
                    <div class="collapsible-body">
                        <span>${extTicket.ticket.Description}</span>
                        <br><br><br>
                        <span class="right">Информация:</span><br><br>
                        <div>
                            <span class="right">Время: ${getTime( new Date(extTicket.time) )}</span><br>
                            <span class="right">Отправитель: ${extTicket.ticket.Sender}</span><br>
                            <span class="right">Статус: ${extTicket.ticket.Status}</span><br>
                        </div>
                        
                        ${!isUser && !isSentLi ? '<a class="waves-effect waves-light btn" style="position: absolute; bottom: 15px; left: 15px" onclick="forwardTicket(' + extTicket.ticket.ID + ')">' +
                            '<i class="material-icons left">forward</i>' +
                            'переслать' +
                        '</a>' : ''}
                    </div>
                </li>
            `;


            if (isSentLi)
                sentLi += ticket;
            else
                incomingLi += ticket;
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

        // Set counters
        /*$('#sentHeader').html(`<h5>У вас <b>${sentCount}</b> отправленных тикетов</h5>`);
        $('#incomingHeader').html(`<h5>У вас <b>${incomingCount}</b> полученных тикетов</h5>`);*/
        // $('#sentCounter').html(sentCount);
        // $('#receivedCounter').html(incomingCount);
    },
    "create": (extTicket) => {
        swal('Успешно!', 'Тикет был отправлен!', 'success');

        console.log("Ticket was successfully created!", extTicket);

        // if (sentNeedsToBeEmpty) {
        //     $('.sentTicketsUl').html('');
        //     sentNeedsToBeEmpty = false;
        // }

        extTicket.ticket.Description = extTicket.ticket.Description.replace(/(?:\r\n|\r|\n)/g, '<br>');
        if (userNames[extTicket.ticket.Sender]) extTicket.ticket.Sender = userNames[extTicket.ticket.Sender];
        else extTicket.ticket.Sender = '<span class="green-text">Модератор</span>';

        let li = `
            <li>
                <div class="collapsible-header"><i class="material-icons">folder</i>${extTicket.ticket.Sender}: «<b>${extTicket.ticket.Caption}</b>»</div>
                <div class="collapsible-body">
                    <span>${extTicket.ticket.Description}</span>
                    <br><br><br>
                    <span class="right">Информация:</span><br><br>
                    <div>
                        <span class="right">Время: ${getTime( new Date(extTicket.time) )}</span><br>
                        <span class="right">Отправитель: ${extTicket.ticket.Sender}</span><br>
                        <span class="right">Статус: ${extTicket.ticket.Status}</span><br>
                    </div>                   
                </div>
            </li>
        `;
        li += $('.sentTicketsUl').html();

        $('.sentTicketsUl').html(li);
        // Set counter
        $('#sentCounter').html(++sentCount).addClass('new');
    },
    "incoming": (extTicket) => {
        console.log("Got new ticket!", extTicket);

        // if (incomingNeedsToBeEmpty) {
        //     $('.incomingTicketsUl').html('');
        //     incomingNeedsToBeEmpty = false;
        // }

        extTicket.ticket.Description = extTicket.ticket.Description.replace(/(?:\r\n|\r|\n)/g, '<br>');
        if (userNames[extTicket.ticket.Sender]) extTicket.ticket.Sender = userNames[extTicket.ticket.Sender];
        else extTicket.ticket.Sender = '<span class="green-text">Модератор</span>';

        let li = `
            <li>
                <div class="collapsible-header"><i class="material-icons">folder</i>${extTicket.ticket.Sender}: «<b>${extTicket.ticket.Caption}</b>»</div>
                <div class="collapsible-body">
                    <span>${extTicket.ticket.Description}</span>
                    <br><br><br>
                    <span class="right">Информация:</span><br><br>
                    <div>
                        <span class="right">Время: ${getTime( new Date(extTicket.time) )}</span><br>
                        <span class="right">Отправитель: ${extTicket.ticket.Sender}</span><br>
                        <span class="right">Статус: ${extTicket.ticket.Status}</span><br>
                    </div>
                    
                    ${!isUser ? '<a class="waves-effect waves-light btn" style="position: absolute; bottom: 15px; left: 15px" onclick="forwardTicket(' + extTicket.ticket.ID + ')">' +
                        '<i class="material-icons left">forward</i>' +
                        'переслать' +
                    '</a>' : ''}
                </div>
            </li>
        `;
        li += $('.incomingTicketsUl').html();

        $('.incomingTicketsUl').html(li);
        // Set counter
        $('#receivedCounter').html(++incomingCount).addClass('new');
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

ws.onmessage = (event) => {
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
    sendEvent('get', {});
}

ws.onclose = () => {
    console.log("Connection was closed")
}

function sendEvent(event, data) {
    let obj = { event, data: JSON.stringify(data) };
    send(JSON.stringify(obj));
}