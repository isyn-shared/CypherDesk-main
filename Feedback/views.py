from django.shortcuts import render
from Feedback.models import FeedbackRecord
from Feedback.models import FeedbackUserIP
from django.http import HttpResponse, Http404, HttpResponseRedirect
from django.conf import settings
from g_recaptcha.validate_recaptcha import validate_captcha
import requests, re, distance, simplejson
from datetime import datetime, timezone

context = {
    'GOOGLE_RECAPTCHA_SITE_KEY': settings.GOOGLE_RECAPTCHA_SITE_KEY,
}

def index (request):
    return render (request, 'Feedback/wrapper.html')

#@validate_captcha
def send (request):
    if request.POST:
        IP = str(request.META['REMOTE_ADDR'])
        now = datetime.now(timezone.utc)

        if FeedbackUserIP.objects.filter(user_ip=IP):
            user_hist_obj = FeedbackUserIP.objects.filter(user_ip=IP)[0]
            period = now - user_hist_obj.date
            if (period.days > 0 or period.seconds > 60 * 60 * 2):
                user_hist_obj.update(date=now)
            else:
                return HttpResponse(2) #превысили лимит запросов
        else:
            FeedbackUserIP.objects.create(user_ip=IP, date=now)

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

        FeedbackRecord.objects.create(title=message_title, problem=message_text, answer="", user_name=user_name,
                                      user_email=user_email)

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

        if not result_mail:
            result = 1 # не удалось отправить почту
        else:
            result = 0
        return HttpResponse(result)
    raise Http404()

def sort_col(i):
    return i['dis']

def found_titles (request):
    TITLES_LIMIT = 5
    if request.GET:
        records = FeedbackRecord.objects.all()
        current_title = request.GET['message_title']
        result = []

        for record in records:
            if record.answer == "":
                continue
            cnt = ComparsionTitles(current_title, record.title)
            if cnt >= 1:
                result.append({'id': record.id, 'title': record.title, 'dis': cnt})

        result.sort(key=sort_col, reverse=True)

        if len(result) > TITLES_LIMIT:
            result = result[:TITLES_LIMIT]

        result = simplejson.dumps(result)
        return HttpResponse(result)
    raise Http404()

def ComparsionTitles (s1, s2):
    arr1 = s1.split(' ')
    arr2 = s2.split(' ')

    exc_dict = ['at', 'in', 'is']

    cnt = 0
    for word1 in arr1:
        for word2 in arr2:
            if word1 in exc_dict or word2 in exc_dict or len(word1) < 2  or len(word2) < 2:
                continue
            tmp_dis = 1 - distance.jaccard(word1.lower(), word2.lower())
            if tmp_dis >= 0.8:
                cnt += 1
    return(cnt)

def get_answer(request):
    result = {}
    if request.GET:
        id = request.GET['id']

        if FeedbackRecord.objects.filter(id=id):
            ans = FeedbackRecord.objects.filter(id=id)[0]
            result['ans'] = ans
            return render(request, 'Feedback/full_ans/wrapper.html', result)
    return HttpResponseRedirect(settings.HOSTNAME + '404/')

def faq(request):
    records = FeedbackRecord.objects.all()
    result = {}
    result['faq'] = []

    for record in records:
        if record.answer == "":
            continue
        result['faq'].append({'id': record.id, 'title': record.title})

    print (result['faq'])
    return render(request, 'Feedback/faq/wrapper.html', result)