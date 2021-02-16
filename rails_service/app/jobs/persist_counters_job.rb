class PersistCountersJob < ApplicationJob
  def perform()
    persist_chats_counters
    persist_messages_counters
    logger.info("[job] job executed")
  end
  
  private
  
  def persist_chats_counters()
    chats_counters_key = 'rails_service:chats_counters'
    chats_counters = {}
    Redis.current.multi do
      chats_counters = Redis.current.hgetall(chats_counters_key)
      Redis.current.del(chats_counters_key)
    end

    chats_counters.value.each do |app_id, count|
      app = App.find(app_id)
      app.with_lock do
        app.update(chats_count: app.chats_count + count.to_i)
      end
    end
  end

  def persist_messages_counters()
    messages_counters_key = 'rails_service:messages_counters'
    messages_counters = {}
    Redis.current.multi do 
      messages_counters = Redis.current.hgetall(messages_counters_key)
      Redis.current.del(messages_counters_key)
    end

    messages_counters.value.each do |chat_id, count|
      chat = Chat.find(chat_id)
      chat.with_lock do
        chat.update(messages_count: chat.messages_count + count.to_i)
      end
    end
  end
end
