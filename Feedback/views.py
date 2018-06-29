from django.shortcuts import render
from Feedback.models import FeedbackRecord
from difflib import get_close_matches
from django.core.mail import send_mail, BadHeaderError, EmailMultiAlternatives
from django.http import HttpResponse, HttpResponseRedirect
from django.conf import settings
from g_recaptcha.validate_recaptcha import validate_captcha
import requests, re

context = {
    'GOOGLE_RECAPTCHA_SITE_KEY': settings.GOOGLE_RECAPTCHA_SITE_KEY,
}

def index (request):
    return render (request, 'Feedback/wrapper.html')

#@validate_captcha
def send (request):
    if request.POST:
        successful_mail = True
        #esult_mail_message = ""

        feedback_data = request.POST
        user_name = feedback_data['user_name']
        user_email = feedback_data['user_email']
        message_title = feedback_data['message_title']
        message_text = feedback_data['message_text']

        """sending emails"""
        from_email = "cypherdesk.isyn@gmail.com"
        email_subject = open(settings.BASE_DIR + "/Feedback/templates/Feedback/mail/title.txt").read()
        email_subject = re.sub("{TITLE}", message_text, email_subject)
        email_html_content = open(settings.BASE_DIR + "/Feedback/templates/Feedback/mail/body.html").read()
        email_text_content = ""

        msg = EmailMultiAlternatives(email_subject, email_text_content, from_email, [user_email,])
        msg.attach_alternative(email_html_content, "text/html")

        """post request on telegram app"""
        telegram_message = "User name: " + user_name + "\nUser email: " + user_email + "\nMessage: " + message_text
        url = settings.HOSTNAME + 'telegram/send/'
        data = {'chat_name': 'feedback', 'token_name': 'feedback', 'text': telegram_message}
        result_telegram_mes = requests.post(url, data=data).text


        if email_subject and email_html_content and from_email and user_email:
            try:
                msg.send()
            except BadHeaderError:
                #result_mail_message = 'Invalid header found.'
                successful_mail = False
            #result_mail_message = 'Thanks, your email was sent'
        else:
            # In reality we'd use a form class
            # to get proper validation errors.
            #result_mail_message = 'Make sure all fields are entered and valid.'
            successful_mail = False

        #result = {'result_mail_message': result_mail_message, 'result_telegram': result_telegram_mes}
        return HttpResponse(successful_mail)

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
