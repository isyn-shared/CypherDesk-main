from django.db import models

class Error500Record (models.Model):
    url = models.CharField(max_length=64)
    date = models.DateTimeField()

    def __str__(self):
        return '500 error on page: {}'.format(self.url)