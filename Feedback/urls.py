from django.conf.urls import include
from django.urls import path
from . import views

urlpatterns = [
    path('', views.index, name='index'),
    path('send/', views.send, name='send'),
    path('found_titles/', views.found_titles, name='found_titles'),
]