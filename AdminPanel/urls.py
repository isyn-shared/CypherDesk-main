from django.urls import path
from . import views

urlpatterns = [
    path('', views.index, name='index'),
    path('authorization/', views.authorization, name='authorization'),
    path('out/', views.out, name='out'),
]