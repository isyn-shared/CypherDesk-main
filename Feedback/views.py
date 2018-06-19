from django.shortcuts import render
from django.http import HttpResponse, HttpResponseRedirect
from Feedback.models import FeedbackRecord
from difflib import get_close_matches
from django.core.mail import send_mail, BadHeaderError, EmailMultiAlternatives
from django.http import HttpResponse, HttpResponseRedirect
from django.conf import settings
from g_recaptcha.validate_recaptcha import validate_captcha

context = {
    'GOOGLE_RECAPTCHA_SITE_KEY': settings.GOOGLE_RECAPTCHA_SITE_KEY,
}

def index (request):
    return render (request, 'Feedback/wrapper.html')

@validate_captcha
def send (request):
    if request.POST:
        HttpResponseRedirect("/")
        feedback_data = request.POST
        user_name = feedback_data['user_name']
        user_email = feedback_data['user_email']
        message_title = feedback_data['message_title']
        message_text = "Hello WORLD!!! My name is " +  user_name   #feedback_data['message_text']

        """sending emails"""
        email_subject = "Что-то там " + message_title
        email_text_content = "Здороу, " + user_name
        from_email = "cypherdesk.isyn@gmail.com"
        email_html_content = "<h1> KekLOL </h1>"

        msg = EmailMultiAlternatives(email_subject, email_text_content, from_email, [user_email, from_email])
        msg.attach_alternative(email_html_content, "text/html")

        if email_subject and email_text_content and from_email and user_email:
            try:
                msg.send()
            except BadHeaderError:
                return HttpResponse('Invalid header found.')
            return HttpResponse('Thanks, your email was sended')
        else:
            # In reality we'd use a form class
            # to get proper validation errors.
            return HttpResponse('Make sure all fields are entered and valid.')

def found_titles (request):
    if request.GET:
        records = FeedbackRecord.objects.all()
        context = dict (records=records)
        current_title = request.GET['message_title']
        titles = []
        result = ''

        for record in records:
            if ComparsionTitles(current_title, record.title):
                titles += [record.title]
                result += record.title + '~'

        return HttpResponse (result)

def ComparsionTitles (s1, s2):
    matches = 0
    words_s1 = s1.split(" ")
    words_s2 = s2.split(" ")

    for word1 in words_s1:
        for word2 in words_s2:
            if len(get_close_matches(word1, [word2])) > 0:
                matches += 1
                break

    if matches >= min(len(words_s1), len(words_s2)):
        return True
    return False
