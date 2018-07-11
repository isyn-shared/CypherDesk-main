from django.db import models
import hashlib

class AdminPanelUser (models.Model):
    username = models.CharField(max_length=30)
    mail = models.CharField(max_length=64)
    first_name = models.CharField(max_length=64)
    last_name = models.CharField(max_length=64)
    password = models.CharField(max_length=64)
    status = models.CharField(max_length=2)

    def __str__(self):
        return 'AdmminPanelUser: {}'.format(self.username)

    def save(self, *args, **kwargs):
        if not self.pk:
            sha256 = hashlib.sha256()
            md5 = hashlib.md5()

            sha256.update(self.password.encode('utf-8'))
            md5.update(sha256.hexdigest().encode('utf-8'))
            self.password = md5.hexdigest()

        super(AdminPanelUser, self).save(args, kwargs)