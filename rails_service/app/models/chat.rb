class Chat < ApplicationRecord
  validates :number, presence: true
  belongs_to :app
  has_many :messages
end
