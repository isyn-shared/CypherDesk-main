from django.shortcuts import render
from Feedback.models import FeedbackRecord
from django.http import HttpResponse, Http404, HttpResponseRedirect
from django.conf import settings
from g_recaptcha.validate_recaptcha import validate_captcha
import requests, re, distance, simplejson

context = {
    'GOOGLE_RECAPTCHA_SITE_KEY': settings.GOOGLE_RECAPTCHA_SITE_KEY,
}

def index (request):
    return render (request, 'Feedback/wrapper.html')

#@validate_captcha
def send (request):
    if request.POST:
        feedback_data = request.POST
        user_name = feedback_data['user_name']
        user_email = feedback_data['user_email']
        message_title = feedback_data['message_title']
        message_text = feedback_data['message_text']

        """sending emails"""
        static_mail_files_path = "/Feedback/templates/Feedback/mail/"
        from_email = settings.EMAIL_HOST_USER
        email_subject = open(settings.BASE_DIR + static_mail_files_path + "title.txt").read()
        email_subject = re.sub("{TITLE}", message_title, email_subject)
        email_html_content = open(settings.BASE_DIR + static_mail_files_path + "body.html").read()
        email_html_content = re.sub("{TITLE}", message_title, email_html_content)
        email_html_content = re.sub("{USERNAME}", user_name, email_html_content)
        email_html_content = re.sub("{CONTACTLINK}", open(settings.BASE_DIR + static_mail_files_path + 'contactlink.txt').read(),
                                    email_html_content)
        email_text_content = ""

        """post request on telegram app"""
        telegram_message = "User name: " + user_name + "\nUser email: " + user_email + "\nMessage: " + message_text
        url = settings.HOSTNAME + 'telegram/send/'
        data = {'chat_name': 'feedback', 'token_name': 'feedback', 'text': telegram_message}
        result_telegram = requests.post(url, data=data).text

        """post request on MailAgent app"""
        url = settings.HOSTNAME + 'mail/send/'
        data = {'from': from_email, 'subject': email_subject, 'html_content': email_html_content, 'text_content': email_text_content,
                'to': user_email}
        result_mail = requests.post(url, data=data).text

        return HttpResponse(result_mail and result_telegram)
    raise Http404()

def found_titles (request):
    if request.GET:
        records = FeedbackRecord.objects.all()
        context = dict(records=records)
        current_title = request.GET['message_title']
        result = {}

        for record in records:
            dis = ComparsionTitles(current_title, record.title)
            if dis >= 0.6:
                result[record.id] = {'title': record.title, 'dis': dis}

        result = simplejson.dumps(result)
        return HttpResponse(result)
    raise Http404()

def ComparsionTitles (s1, s2):
    arr1 = s1.split(' ')
    arr2 = s2.split(' ')

    sum_dis = 0
    for word1 in arr1:
        tmp_dis = 0
        for word2 in arr2:
            tmp_dis = max(tmp_dis, 1 - distance.jaccard(word1.lower(), word2.lower()))
        sum_dis += tmp_dis

    sum_dis = sum_dis / min(len(arr1), len(arr2))
    return(sum_dis)

def get_answer(request):
    result = {}
    if request.GET:
        id = request.GET['id']

        if FeedbackRecord.objects.filter(id=id):
            ans = FeedbackRecord.objects.filter(id=id)[0]
            result['ans'] = ans
            return render(request, 'Feedback/full_ans/wrapper.html', result)
    return HttpResponseRedirect(settings.HOSTNAME + '404/')