from django.shortcuts import render
from django.http import HttpResponseRedirect
from News.models import NewsRecord
from django.conf import settings
import re, textile

def index(request):
    records = NewsRecord.objects.all()
    result = {}
    result['news'] = []
    for record in records:
        result['news'].insert(0, {'title': record.title, 'short_news': record.short_news, 'news': record.news,
                                  'date': record.date, 'id': record.id})

    return render(request, 'News/wrapper.html', result)

def full(request):
    result = {}
    im_tem1 = '[#file1#]'
    im_tem2 = '[#file2#]'
    im_tem3 = '[#file3#]'

    path = "/News/templates/News/"
    img_html_template = open(settings.BASE_DIR + path + "img_include.html").read()

    if (request.GET):
        id = request.GET['id']

        if NewsRecord.objects.filter(id=id):
            post = NewsRecord.objects.filter(id=id)[0]
            text = post.news

            if post.file1:
                A_SRC = settings.HOSTNAME + post.file1.url[1:]
                tmp_include = re.sub("{SRC}", post.file1.url, img_html_template)
                tmp_include = re.sub("{A_SRC}", A_SRC, tmp_include)
                text = text.replace(im_tem1, tmp_include)
            if post.file2:
                A_SRC = settings.HOSTNAME + post.file2.url[1:]
                tmp_include = re.sub("{SRC}", post.file2.url, img_html_template)
                tmp_include = re.sub("{A_SRC}", A_SRC, tmp_include)
                text = text.replace(im_tem2, tmp_include)
            if post.file3:
                A_SRC = settings.HOSTNAME + post.file3.url[1:]
                tmp_include = re.sub("{SRC}", post.file3.url, img_html_template)
                tmp_include = re.sub("{A_SRC}", A_SRC, tmp_include)
                text = text.replace(im_tem3, tmp_include)
            post.news = textile.textile(text)
            result['post'] = post
            return render(request, 'News/Full/wrapper.html', result)
    
    return HttpResponseRedirect(settings.HOSTNAME + '404/')