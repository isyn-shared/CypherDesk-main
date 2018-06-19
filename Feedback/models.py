from django.db import models

class FeedbackRecord (models.Model):
    title = models.CharField(max_length=64) #текстовое поле (максимальное кол-во символов: 64)
    problem = models.TextField()
    answer = models.TextField()

    def __str__(self):
        return 'FeedbackRecord: {}'.format(self.title)