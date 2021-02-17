class Chat < ApplicationRecord
  validates :number, presence: true
  belongs_to :app
  has_many :messages

  def as_json(options = {})
    super(options.merge({ only: [:number, :name] }))
  end

  def self.get_one(app_token, chat_number)
    joins(:app).where(apps: { token: app_token }, chats: { number: chat_number }).limit(1).first
  end

  def self.get_all(app_token)
    joins(:app).where(apps: { token: app_token })
  end

end
