class UpdateElasticsearchIndexJob < ApplicationJob
  def perform(message_id, method)
    case method
    when 'create'
      message = Message.find(message_id)
      message.__elasticsearch__.index_document
    when 'update'
      message = Message.find(message_id)
      message.__elasticsearch__.update_document
    end
    logger.info("[job] job executed")
  end
end
  