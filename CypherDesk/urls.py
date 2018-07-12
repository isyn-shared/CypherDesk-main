from django.contrib import admin
from django.conf.urls import include, handler404, handler500
from django.urls import path
from LandPage import views as Landpage_views

urlpatterns = [
    path('', include('LandPage.urls')),
    path('custom_admin/', include('AdminPanel.urls')),
    path('standart_admin/', admin.site.urls),
    path('feedback/', include('Feedback.urls')),
    path('news/', include('News.urls')),
    path('telegram/', include('TelegramBot.urls')),
    path('mail/', include('MailAgent.urls')),
]

handler404 = Landpage_views.error_404
handler500 = Landpage_views.error_500