class Chat < ApplicationRecord
  validates :number, presence: true
  belongs_to :app
  has_many :messages

  def as_json(options = {})
    super(options.merge({ only: [:number, :name] }))
  end
end
