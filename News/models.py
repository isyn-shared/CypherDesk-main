from django.db import models
import random

MAX_SIZE = 1048626
ALLOWED_EXTENSIONS = ['jpg', 'png', 'bmp', 'gif', 'mp3', 'mp4', 'mov', 'txt', 'docx', 'doc', 'jpeg', '7z', 'zip',]

class NewsRecord (models.Model):
    title = models.CharField(max_length=64)
    short_news = models.CharField(max_length=128)
    news = models.TextField()
    date = models.DateTimeField()
    file1 = models.FileField(null=True, blank=True)
    file2 = models.FileField(null=True, blank=True)
    file3 = models.FileField(null=True, blank=True)

    def save(self, *args, **kwargs):
        if not self.pk:
            if self.file1:
                parts = self.file1.name.split('.')
                file_extension = parts[len(parts) - 1]

                if NewsRecord.chkFile(file_extension, self.file1.size):
                    self.file1.name = "News_" + self.title + "_file1_" + str(random.randint(0, 9999)) + "." + file_extension
                else:
                    return
            if self.file2:
                parts = self.file2.name.split('.')
                file_extension = parts[len(parts) - 1]

                if NewsRecord.chkFile(file_extension, self.file1.size):
                    self.file2.name = "News_" + self.title + "_file2_" + str(random.randint(0, 9999)) + "." + file_extension
                else:
                    return
            if self.file3:
                parts = self.file3.name.split('.')
                file_extension = parts[len(parts) - 1]
                if NewsRecord.chkFile(file_extension, self.file2.size):
                    self.file3.name = "News_" + self.title + "_file3_" + str(random.randint(0, 9999)) + "." + file_extension
                else:
                    return
        super(NewsRecord, self).save(args, kwargs)

    @staticmethod
    def chkFile(file_extension, file_size):
        """
        if file_extension not in ALLOWED_EXTENSIONS:
            return False
        if file_size > MAX_SIZE:
            return False
        """
        return True

    def __str__(self):
        return 'News Record: {}'.format(self.title)