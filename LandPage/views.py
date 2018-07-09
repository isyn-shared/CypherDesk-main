from django.shortcuts import render
from django.http import HttpResponse
from django.conf import settings
import requests

def index(request):
    return render(request, 'LandPage/wrapper.html')

def error_500(request):
    data = {}
    telegram_message = "Ошибка 500 по адресу: " + settings.HOSTNAME + request.path
    url = settings.HOSTNAME + 'telegram/send/'
    data = {'chat_name': 'feedback', 'token_name': 'feedback', 'text': telegram_message}
    result_telegram = requests.post(url, data=data).text

    data['telegram'] = result_telegram

    return render(request, 'LandPage/500/error_500.html', data)

def error_404(request, exception):
    data = {}
    return render(request,'LandPage/404/error_404.html', data)