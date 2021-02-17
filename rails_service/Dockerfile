FROM ruby:2.5.1

RUN apt-get update -qq && apt-get install -y build-essential libpq-dev cron
RUN mkdir /app

WORKDIR /app

ADD Gemfile /app/Gemfile
ADD Gemfile.lock /app/Gemfile.lock

RUN bundle install

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.7.3/wait /wait
RUN chmod +x /wait

ADD . /app
# CMD rake sneakers:run