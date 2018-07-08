from django.contrib import admin
from django.conf.urls import include, handler404
from django.urls import path

urlpatterns = [
    path('', include('LandPage.urls')),
    path('admin/', admin.site.urls),
    path('feedback/', include('Feedback.urls')),
    path('news/', include('News.urls')),
    path('telegram/', include('TelegramBot.urls')),
    path('mail/', include('MailAgent.urls')),
]

#handler404 = 'LandPage.views.error_404'