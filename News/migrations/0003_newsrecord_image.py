# Generated by Django 2.0.6 on 2018-07-12 09:21

from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ('News', '0002_auto_20180629_1446'),
    ]

    operations = [
        migrations.AddField(
            model_name='newsrecord',
            name='image',
            field=models.FileField(blank=True, null=True, upload_to=''),
        ),
    ]
