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

  __elasticsearch__.create_index! unless __elasticsearch__.index_exists?
end
