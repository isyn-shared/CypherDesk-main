from django.shortcuts import render
from django.http import HttpResponse
from News.models import NewsRecord

def index(request):
    records = NewsRecord.objects.all()
    result = {}

    result['news'] = []
    for record in records:
        result['news'].insert(0, {'title': record.title, 'short_news': record.short_news, 'news': record.news, 'date': record.date})

    return render(request, 'News/wrapper.html', result)