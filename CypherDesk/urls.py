from django.contrib import admin
from django.urls import path

urlpatterns = [
    path('admin/', admin.site.urls),
    path('', include('LandPage.urls')),
    path('feedback/', include('Feedback.urls')),
]
