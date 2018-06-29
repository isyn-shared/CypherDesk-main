from django.shortcuts import render
from django.http import HttpResponse
from News.models import NewsRecord

def index(request):
    records = NewsRecord.objects.all()
    AllRecords = {}
    result = {}
    RecordsStr = ""

    result['news'] = []
    it = 0
    for record in records:
        result['news'] += [{'title': record.title, 'short_news': record.short_news, 'news': record.news, 'date': record.date},]
        it += 1

    return render(request, 'News/wrapper.html', result)