# frozen_string_literal: true
require 'sneakers'
require 'sneakers/handlers/maxretry'

Sneakers.configure(amqp: ENV['RABBITMQ_URL'], daemonize: false, handler: Sneakers::Handlers::Maxretry, env: ENV['RAILS_ENV'])
Sneakers.logger.level = Logger::DEBUG
Sneakers.logger = Rails.logger
Sneakers::Worker.logger = Rails.logger
WORKER_OPTIONS = {
  retry_timeout: 5 * 1000, # 5 sec
  ack: true,
  threads: 10,
  prefetch: 10,
}
