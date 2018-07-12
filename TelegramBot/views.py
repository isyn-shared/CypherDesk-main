from . import TelegramBotClass
from django.conf import settings
from django.views.decorators.csrf import csrf_protect, csrf_exempt
from django.http import HttpResponse, HttpResponseRedirect
from django.core.exceptions import PermissionDenied

@csrf_exempt
def send(request):
    client_ip = str(request.META['REMOTE_ADDR'])
    if client_ip != '127.0.0.1':
        return PermissionDenied

    if request.POST:
        post_request = request.POST
        token_name = post_request['token_name']
        token = settings.FEEDBACK_TELEGRAM_BOT_KEY['feedback']

        chat_name = post_request['chat_name']
        chat_id = settings.FEEDBACK_TELEGRAM_CHAT_ID['feedback']

        text = post_request['text']

        telegram_bot_obj = TelegramBotClass.TelegramBot(token)
        result = telegram_bot_obj.send_message(chat_id, text)

        return HttpResponse(True)
    return HttpResponseRedirect(settings.HOSTNAME + '404/')