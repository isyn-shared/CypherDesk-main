# Generated by Django 2.0.6 on 2018-07-12 10:16

from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ('News', '0003_newsrecord_image'),
    ]

    operations = [
        migrations.RenameField(
            model_name='newsrecord',
            old_name='image',
            new_name='image1',
        ),
        migrations.AddField(
            model_name='newsrecord',
            name='image2',
            field=models.FileField(blank=True, null=True, upload_to=''),
        ),
        migrations.AddField(
            model_name='newsrecord',
            name='image3',
            field=models.FileField(blank=True, null=True, upload_to=''),
        ),
    ]
