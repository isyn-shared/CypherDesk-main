from django.db import models

class NewsRecord (models.Model):
    title = models.CharField(max_length=64)
    short_news = models.CharField(max_length=128)
    news = models.TextField()
    date = models.DateTimeField()

    def __str__(self):
        return 'News Record: {}'.format(self.title)