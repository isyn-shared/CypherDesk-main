from django.shortcuts import render
from django.http import HttpResponse

def index(request):
    return render(request, 'LandPage/wrapper.html')

def error_500(request):
    data = {}
    return render(request, 'LandPage/500/error_500.html', data)

def error_404(request, exception):
    data = {}
    return render(request,'LandPage/404/error_404.html', data)