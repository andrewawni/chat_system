class App < ApplicationRecord
  validates :token, presence: true
  has_many :chats
end
