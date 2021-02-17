# frozen_string_literal: true

class Message < ApplicationRecord
  # include Searchable
  include Elasticsearch::Model
  validates :number, presence: true
  belongs_to :chat

  after_create do
    UpdateElasticsearchIndexJob.perform_later(id, 'create')
  end

  after_update do
    UpdateElasticsearchIndexJob.perform_later(id, 'update')
  end

  mappings dynamic: 'false' do
    indexes :content, analyzer: 'english', index_options: 'offsets'
  end

  def as_json(options = {})
    super(options.merge({ only: [:number, :content] }))
  end

  def self.get_one(app_token, chat_number, message_number)
    joins(chat: [:app]).where(apps: { token: app_token }, chats: { number: chat_number }, messages: { number: message_number }).limit(1).first
  end

  def self.get_all(app_token, chat_number)
    joins(chat: [:app]).where(apps: { token: app_token }, chats: { number: chat_number })
  end

  __elasticsearch__.create_index! unless __elasticsearch__.index_exists?
end
