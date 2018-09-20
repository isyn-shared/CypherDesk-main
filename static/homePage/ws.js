const url = `ws://${window.location.host}/tickets/ws`;
const ws = new WebSocket(url);

function send(data) {
    console.log("Sending", data);
    ws.send(data)
}

let tickets = [{
	ticket: {
		ID: 1,
		Caption: "Hello",
		Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
		Sender: 1,
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
		Caption: "Hello",
		Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
		Sender: 1,
		Status: "opened"
    },
	forwardFrom: 1,
	forwardTo: 2,
    Action: "create",
	Time: "Когда-то"
}]

const myEvents = {
    // Special multi-purpose event 
    "error": (err) => { console.error(err); },
    "get": (tickets) => {
        console.log("Got tickets!", tickets);
        let li = '';
        for (let extTicket of tickets) {
            li += `
                <li class="list-group-item mb-2" style="float: left">
                    <h5 class="font-weight-bold">Тема: <span class="font-weight-normal">${extTicket.ticket.Caption}</span></h5>
                    <p class="mb-0 lead">${extTicket.ticket.Description}</p>

                    <div style="float: right">
                        <a class="text-muted" data-toggle="collapse" href="#collapseInfo${extTicket.ticket.ID}" aria-expanded="false" aria-controls="collapseExample">
                            Расширить информацию ▼
                        </a>
                        
                        <div class="collapse" id="collapseInfo${extTicket.ticket.ID}">
                            <p class="font-weight-bold mb-0">Время: <span class="font-weight-normal">${extTicket.Time}</span></p>
                            <p class="font-weight-bold mb-0">Отправитель: <span class="font-weight-normal">${extTicket.ticket.Sender}</span></p>
                            <p class="font-weight-bold mb-0">Статус: <span class="font-weight-normal">${extTicket.ticket.Status}</span></p>
                        </div>
                    </div>
                </li>
            `;
        }
        $('.sentTicketsUl').html(li);

    }
}
myEvents['get'](tickets);

ws.onmessage = (event) => {
    console.log(event);
    let msg = JSON.parse(event.data);
    console.log(msg);

    if (msg.ok === false) {
        return console.warn("Ошибка!", msg.errorMessage)
    }

    console.log(`Attempt to execute '${msg.event}' with data:`, msg.data);

    if (!myEvents[msg.event])
        return console.log("No event " + msg.event);

    myEvents[msg.event](msg.data);
}

ws.onopen = () => {
    console.log("Connected successfully!");

    /* Send initial events here */
}

ws.onclose = () => {
    console.log("Connection was closed")
}

function sendEvent(event, data) {
    let obj = { event, data: JSON.stringify(data) };
    send(JSON.stringify(obj));
}