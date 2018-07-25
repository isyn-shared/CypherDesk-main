from django.urls import path
from . import views, feedback_table, error_500_table, admin_table, news_table

urlpatterns = [
    path('', views.index, name='index'),
    path('authorization/', views.authorization, name='authorization'),
    path('out/', views.out, name='out'),
    path('feedback/', feedback_table.index, name='feedback_index'),
    path('news/', news_table.index, name='news_index'),
    path('admin/', admin_table.index, name='admin_index'),
    path('error_500/', error_500_table.index, name='error_500_index')
]