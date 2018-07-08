from django.shortcuts import render
from django.http import HttpResponse

def index(request):
    return render(request, 'LandPage/wrapper.html')

def error_500(request):
    return HttpResponse("500, дебил")

def error_404(request):
    return HttpResponse("404, дебил")