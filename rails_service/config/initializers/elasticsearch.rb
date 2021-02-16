# frozen_string_literal: true
require 'elasticsearch/model'

client = Elasticsearch::Client.new(url: ENV['ELASTICSEARCH_URL'])
Elasticsearch::Model.client = client
