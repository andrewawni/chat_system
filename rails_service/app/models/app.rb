class App < ApplicationRecord
  validates :token, presence: true
  has_many :chats

  def as_json(options = {})
    super(options.merge({ only: [:token, :name] }))
  end
  
end
