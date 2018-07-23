from django.contrib import admin
from Feedback.models import FeedbackRecord
from Feedback.models import FeedbackUserIP

admin.site.register(FeedbackRecord)
admin.site.register(FeedbackUserIP)