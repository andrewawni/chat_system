class Message < ApplicationRecord
  validates :number, presence: true
  belongs_to :chat
end
