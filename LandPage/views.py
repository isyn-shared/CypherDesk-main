from django.shortcuts import render
from django.http import HttpResponse
from django.conf import settings
from LandPage.models import Error500Record
from datetime import datetime, timezone, timedelta

import requests

def index(request):
    return render(request, 'LandPage/wrapper.html')

def error_500(request):
    data = {}
    url = settings.HOSTNAME + 'telegram/send/'

    if not Error500Record.objects.filter(url=url):
        page_error = settings.HOSTNAME + request.path[1:]
        telegram_message = "Ошибка 500 по адресу: " + page_error
        data = {'chat_name': 'feedback', 'token_name': 'feedback', 'text': telegram_message}
        result_telegram = requests.post(url, data=data).text
        data['telegram'] = result_telegram
        now = datetime.now(timezone.utc) + timedelta(minutes=180)
        error500_insert = Error500Record.objects.create(url=page_error, date=now)
        data['insertion'] = error500_insert

    return render(request, '500/error_500.html', data)

def error_404(request, exception):
    data = {}
    return render(request,'404/error_404.html', data)

def error_403(request):
    data = {"reason": "Извините, вы не имеете доступа к этим данным"}
    return render(request,'403/error_403.html', data)