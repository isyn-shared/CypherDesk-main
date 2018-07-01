from django.shortcuts import render
from django.http import HttpResponse
from News.models import NewsRecord

def index(request):
    records = NewsRecord.objects.all()
    result = {}

    result['news'] = []
    for record in records:
        result['news'].insert(0, {'title': record.title, 'short_news': record.short_news, 'news': record.news, 'date': record.date, 'id': record.id})

    return render(request, 'News/wrapper.html', result)

def full(request):
    result = {}

    if (request.GET):
        id = request.GET['id']
        post = NewsRecord.objects.filter(id=id)[0]
        result['post'] = post
        return render(request, 'News/Full/wrapper.html', result)
    
    return HttpResponse('Ошибка 404 - Вы Кек')