# frozen_string_literal: true

class PersistRecordsWorker
  include Sneakers::Worker

  from_queue 'chat_system'

  def work(raw_message)
    payload = JSON.parse(raw_message)
    status = case payload['type']
             when 'application'
               work_for_apps(payload)
             when 'chat'
               work_for_chats(payload)
             when 'message'
               work_for_messages(payload)
             else
               ack!
             end
    if status
      logger.info("[sneakers] #{payload['type']} consumed and created successfully")
      ack!
    else
      logger.info("[sneakers] #{payload['type']} consumed and rejected")
      reject!
    end
  end

  private

  def work_for_apps(payload)
    app = App.new(name: payload['application_name'], token: payload['application_token'])
    app.save
  end

  def work_for_chats(payload)
    app = App.find_by(token: payload['application_token'])
    chat = Chat.new(app: app, name: payload['chat_name'], number: payload['chat_number'])
    Redis.current.hincrby('rails_service:chats_counters', app.id, 1) if chat.valid?
    chat.save
  end

  def work_for_messages(payload)
    app = App.find_by(token: payload['application_token'])
    chat = app.chats.find_by(number: payload['chat_number'])
    message = Message.new(chat: chat, number: payload['message_number'], content: payload['message_content'])
    Redis.current.hincrby('rails_service:messages_counters', chat.id, 1) if message.valid?
    message.save
  end
end
